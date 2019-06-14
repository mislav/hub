Feature: hub pr show
  Background:
    Given I am in "git://github.com/ashemesh/hub.git" git repo
    And I am "ashemesh" on github.com with OAuth token "OTOKEN"

  Scenario: Current branch
    Given I am on the "topic" branch
    Given the GitHub API server:
      """
      get('/repos/ashemesh/hub/pulls'){
        assert :state => "open",
               :head => "ashemesh:topic"
        json [
          { :html_url => "https://github.com/ashemesh/hub/pull/102" },
        ]
      }
      """
    When I successfully run `hub pr show`
    Then "open https://github.com/ashemesh/hub/pull/102" should be run

  Scenario: Current branch output URL
    Given I am on the "topic" branch
    Given the GitHub API server:
      """
      get('/repos/ashemesh/hub/pulls'){
        assert :state => "open",
               :head => "ashemesh:topic"
        json [
          { :html_url => "https://github.com/ashemesh/hub/pull/102" },
        ]
      }
      """
    When I successfully run `hub pr show -u`
    Then "open https://github.com/ashemesh/hub/pull/102" should not be run
    And the output should contain exactly:
      """
      https://github.com/ashemesh/hub/pull/102\n
      """

  Scenario: Current branch in fork
    Given the "upstream" remote has url "git@github.com:github/hub.git"
    And I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      get('/repos/github/hub/pulls'){
        assert :state => "open",
               :head => "ashemesh:topic"
        json [
          { :html_url => "https://github.com/github/hub/pull/102" },
        ]
      }
      """
    When I successfully run `hub pr show`
    Then "open https://github.com/github/hub/pull/102" should be run

  Scenario: Explicit head branch
    Given the GitHub API server:
      """
      get('/repos/ashemesh/hub/pulls'){
        assert :state => "open",
               :head => "ashemesh:topic"
        json [
          { :html_url => "https://github.com/ashemesh/hub/pull/102" },
        ]
      }
      """
    When I successfully run `hub pr show --head topic`
    Then "open https://github.com/ashemesh/hub/pull/102" should be run

  Scenario: Explicit head branch with owner
    Given the GitHub API server:
      """
      get('/repos/ashemesh/hub/pulls'){
        assert :state => "open",
               :head => "github:topic"
        json [
          { :html_url => "https://github.com/ashemesh/hub/pull/102" },
        ]
      }
      """
    When I successfully run `hub pr show --head github:topic`
    Then "open https://github.com/ashemesh/hub/pull/102" should be run

  Scenario: No pull request found
    Given the GitHub API server:
      """
      get('/repos/ashemesh/hub/pulls'){
        json []
      }
      """
    When I run `hub pr show --head topic`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      no open pull requests found for branch 'ashemesh:topic'\n
      """
