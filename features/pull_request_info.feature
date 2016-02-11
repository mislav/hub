Feature: hub pull-request info
  Background:
    Given I am in "git://github.com/github/hub.git" git repo
    And I am "cornwe19" on github.com with OAuth token "OTOKEN"

  Scenario: Basic pull request
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message "ăéñøü"
    And the "topic" branch is pushed to "origin/topic"
    And the GitHub API server:
      """
      get('/repos/github/hub/pulls') {
        assert :base => "master",
               :head => "github:topic"
        json([
          { :title => "Description of my own pull request",
            :comments_url => "https://api.github.com/repos/github/hub/issues/4242/comments",
            :html_url => "https://github.com/github/hub/pull/4242"
          }
        ])
      }
      get('/repos/github/hub/issues/4242/comments') {
        json([
          { :user => {
              :login => "pcorpet"
            },
            :body => "A nice comment.",
            :created_at => "2016-02-11T11:02:50Z",
            :updated_at => "2016-02-11T15:27:33Z"
          },
          { :user => {
              :login => "lascap"
            },
            :body => "A second comment.",
            :created_at => "2016-02-11T11:02:54Z",
            :updated_at => "2016-02-11T11:02:54Z",
          }
        ])
      }
      """
    When I run `hub pull-request info`
    Then the output should contain exactly:
      """
      Title: "Description of my own pull request"
      * 2016-02-11 15:27:33 - pcorpet: A nice comment.
      * 2016-02-11 11:02:54 - lascap: A second comment.
      URL: https://github.com/github/hub/pull/4242

      """

  Scenario: No pull request for the branch yet
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message "ăéñøü"
    And the "topic" branch is pushed to "origin/topic"
    And the GitHub API server:
      """
      get('/repos/github/hub/pulls') {
        assert :base => "master",
               :head => "github:topic"
        json([])
      }
      """
    When I run `hub pull-request info`
    Then the output should contain exactly "no such pull request\n"

  Scenario: Very long lines
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message "ăéñøü"
    And the "topic" branch is pushed to "origin/topic"
    And the GitHub API server:
      """
      get('/repos/github/hub/pulls') {
        assert :base => "master",
               :head => "github:topic"
        json([
          { :title => "Very long description that could go on and on for hours, but I am just too lazy to make it even bigger. You got the spirit I think.",
            :comments_url => "https://api.github.com/repos/github/hub/issues/4242/comments",
            :html_url => "https://github.com/github/hub/pull/4242"
          }
        ])
      }
      get('/repos/github/hub/issues/4242/comments') {
        json([
          { :user => {
              :login => "pcorpet"
            },
            :body => "A very nice comment that is so long that it wouldn't fit on one line in a terminal, even with a very very large terminal I really doubt that it would fit entirely. Except if you had reduced the font on purpose.",
            :created_at => "2016-02-11T11:02:50Z",
            :updated_at => "2016-02-11T15:27:33Z"
          },
          { :user => {
              :login => "lascap"
            },
            :body => "A second comment. This one is also very long, but don't let that stop you. I could have written some very understanding towards the end of the line. Or not.",
            :created_at => "2016-02-11T11:02:54Z",
            :updated_at => "2016-02-11T11:02:54Z",
          }
        ])
      }
      """
    When I run `hub pull-request info`
    Then the output should contain exactly:
      """
      Title: "Very long description that could go on and on for hours, but I am just too lazy …"
      * 2016-02-11 15:27:33 - pcorpet: A very nice comment that is so long that it wouldn't fit on one line in a termin…
      * 2016-02-11 11:02:54 - lascap: A second comment. This one is also very long, but don't let that stop you. I cou…
      URL: https://github.com/github/hub/pull/4242

      """

  Scenario: Text on several lines
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message "ăéñøü"
    And the "topic" branch is pushed to "origin/topic"
    And the GitHub API server:
      """
      get('/repos/github/hub/pulls') {
        assert :base => "master",
               :head => "github:topic"
        json([
          { :title => "A title is always on one line only.",
            :comments_url => "https://api.github.com/repos/github/hub/issues/4242/comments",
            :html_url => "https://github.com/github/hub/pull/4242"
          }
        ])
      }
      get('/repos/github/hub/issues/4242/comments') {
        json([
          { :user => {
              :login => "pcorpet"
            },
            :body => "A comment\nthat is on 2 lines.",
            :updated_at => "2016-02-11T15:27:33Z"
          },
          { :user => {
              :login => "lascap"
            },
            :body => "Another comment.\nAlso on 2 lines.",
            :created_at => "2016-02-11T11:02:54Z",
            :updated_at => "2016-02-11T11:02:54Z",
          }
        ])
      }
      """
    When I run `hub pull-request info`
    Then the output should contain exactly:
      """
      Title: "A title is always on one line only."
      * 2016-02-11 15:27:33 - pcorpet: A comment…
      * 2016-02-11 11:02:54 - lascap: Another comment.…
      URL: https://github.com/github/hub/pull/4242

      """
