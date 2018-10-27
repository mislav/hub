Feature: hub issue
  Background:
    Given I am in "git://github.com/github/hub.git" git repo
    And I am "cornwe19" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch issues
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues/102') { json \
        :base => {
          :repo => {
            :owner => { :login => "defunkt" },
            :name => "hub",
            :private => false
          }
        },
        :head => {
          :ref => "hub_merge",
          :repo => {
            :owner => { :login => "jfirebaugh" },
            :name => "hub",
            :private => false
          }
        },
        :title => "Add `hub merge` command"
      }

      get('/repos/github/hub/issues/102/comments') { json \
              :base => {
                :repo => {
                  :owner => { :login => "defunkt" },
                  :name => "hub",
                  :private => false
                }
              },
              :head => {
                :ref => "hub_merge",
                :repo => {
                  :owner => { :login => "jfirebaugh" },
                  :name => "hub",
                  :private => false
                }
              },
              :title => "Add `hub merge` command"
      }
      """
    When I successfully run `hub issue view 102`
    Then the output should contain exactly:
      """
          #102  First issue
      """

  Scenario: Fetch single issue
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues/102') { json \
        :number => 102,
        :body => "I want this feature",
        :title => "Feature request for hub issue view",
        :created_at => "2017-04-14T16:00:49Z",
        :user => { :login => "royels" }
     }
      get('/repos/github/hub/issues/102/comments') {
      json [{
              :id => 1,
              :body => "I am from the future",
              :created_at => "2011-04-14T16:00:49Z",
              :user => { :login => "octocat" }}
      ]
      }
      """
    When I successfully run `hub issue view 102`
    Then the output should contain exactly:
      """

      # Feature request for hub issue view

      * created by @royels on 2017-04-14 16:00:49 +0000 UTC
      * assignees: \nI want this feature
      ## Comments:
      ### comment by @octocat on 2011-04-14 16:00:49 +0000 UTC

      I am from the future

      """

  Scenario: Did not supply an issue number
    When I run `hub issue view`
    Then the exit status should be 1
    Then the output should contain exactly "Usage: hub issue view <NUMBER>\n"


  Scenario: Show error message if http code is not 200 for issues endpoint
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues/102') {
    status 500
    json \
        :number => 102,
        :body => "I want this feature",
        :title => "Feature request for hub issue view",
        :created_at => "2017-04-14T16:00:49Z",
        :user => { :login => "royels" }
     }
      get('/repos/github/hub/issues/102/comments') {
      json [{
              :id => 1,
              :body => "I am from the future",
              :created_at => "2011-04-14T16:00:49Z",
              :user => { :login => "octocat" }}
      ]
      }
      """
    When I run `hub issue view 102`
    Then the output should contain exactly:
      """
      Unable to find issue with number: 102
      
      """


  Scenario: Show error message if http code is not 200 for comments endpoint
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues/102') { json \
        :number => 102,
        :body => "I want this feature",
        :title => "Feature request for hub issue view",
        :created_at => "2017-04-14T16:00:49Z",
        :user => { :login => "royels" }
     }
      get('/repos/github/hub/issues/102/comments') {
      status 404
      json [{
              :id => 1,
              :body => "I am from the future",
              :created_at => "2011-04-14T16:00:49Z",
              :user => { :login => "octocat" }}
      ]
      }
      """
    When I run `hub issue view 102`
    Then the output should contain exactly:
      """
      Unable to get comments for issue 102

      """