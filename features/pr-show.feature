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

  Scenario: Format Current branch output URL
    Given I am on the "topic" branch
    Given the GitHub API server:
      """
      get('/repos/ashemesh/hub/pulls'){
        assert :state => "open",
               :head => "ashemesh:topic"
        json [{
          :number => 102,
          :state => "open",
          :base => {
            :ref => "master",
            :label => "github:master",
            :repo => { :owner => { :login => "github" } }
          },
          :head => { :ref => "patch-1", :label => "octocat:patch-1" },
          :user => { :login => "octocat" },
          :requested_reviewers => [
            { :login => "rey" },
          ],
          :requested_teams => [
            { :slug => "troopers" },
            { :slug => "cantina-band" },
          ],
          :html_url => "https://github.com/ashemesh/hub/pull/102",
        }]
      }
      """
    When I successfully run `hub pr show -f "%sC%>(8)%i %rs%n"`
    Then "open https://github.com/ashemesh/hub/pull/102" should not be run
    And the output should contain exactly:
      """
          #102 rey, github/troopers, github/cantina-band\n\n
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

  Scenario: Differently named branch in fork
    Given the "upstream" remote has url "git@github.com:github/hub.git"
    And I am on the "local-topic" branch with upstream "origin/remote-topic"
    Given the GitHub API server:
      """
      get('/repos/github/hub/pulls'){
        assert :head => "ashemesh:remote-topic"
        json [
          { :html_url => "https://github.com/github/hub/pull/102" },
        ]
      }
      """
    When I successfully run `hub pr show`
    Then "open https://github.com/github/hub/pull/102" should be run

  Scenario: Upstream configuration with HTTPS URL
    Given I am on the "local-topic" branch
    When I successfully run `git config branch.local-topic.remote https://github.com/octocat/hub.git`
    When I successfully run `git config branch.local-topic.merge refs/remotes/remote-topic`
    Given the GitHub API server:
      """
      get('/repos/ashemesh/hub/pulls'){
        assert :head => "octocat:remote-topic"
        json [
          { :html_url => "https://github.com/github/hub/pull/102" },
        ]
      }
      """
    When I successfully run `hub pr show`
    Then "open https://github.com/github/hub/pull/102" should be run

  Scenario: Upstream configuration with SSH URL
    Given I am on the "local-topic" branch
    When I successfully run `git config branch.local-topic.remote git@github.com:octocat/hub.git`
    When I successfully run `git config branch.local-topic.merge refs/remotes/remote-topic`
    Given the GitHub API server:
      """
      get('/repos/ashemesh/hub/pulls'){
        assert :head => "octocat:remote-topic"
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

  Scenario: Show pull request by number
    When I successfully run `hub pr show 102`
    Then "open https://github.com/ashemesh/hub/pull/102" should be run

  Scenario: Format pull request by number
    Given the GitHub API server:
      """
      get('/repos/ashemesh/hub/pulls/102') {
        json :number => 102,
          :title => "First",
          :state => "open",
          :base => {
            :ref => "master",
            :label => "github:master",
            :repo => { :owner => { :login => "github" } }
          },
          :head => { :ref => "patch-1", :label => "octocat:patch-1" },
          :user => { :login => "octocat" },
          :requested_reviewers => [
            { :login => "rey" },
          ],
          :requested_teams => [
            { :slug => "troopers" },
            { :slug => "cantina-band" },
          ]
      }
      """
    When I successfully run `hub pr show 102 -f "%sC%>(8)%i %rs%n"`
    Then "open https://github.com/ashemesh/hub/pull/102" should not be run
    And the output should contain exactly:
      """
          #102 rey, github/troopers, github/cantina-band\n\n
      """

  Scenario: Show pull request by invalid number
    When I run `hub pr show XYZ`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      invalid pull request number: 'XYZ'\n
      """
