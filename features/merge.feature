Feature: hub merge
  Background:
    Given I am in "hub" git repo
    And the "origin" remote has url "git://github.com/defunkt/hub.git"
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Normal merge
    When I run `hub merge master`
    Then the git command should be unchanged

  Scenario: Merge pull request
    Given the GitHub API server:
      """
      get('/repos/defunkt/hub/pulls/164') { json \
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
    And there is a git FETCH_HEAD
    When I successfully run `hub merge https://github.com/defunkt/hub/pull/164`
    Then "git fetch origin refs/pull/164/head" should be run
    And "git merge FETCH_HEAD --no-ff -m Merge pull request #164 from jfirebaugh/hub_merge" should be run
    When I successfully run `git show -s --format=%B`
    Then the output should contain:
      """
      Merge pull request #164 from jfirebaugh/hub_merge

      Add `hub merge` command
      """

  Scenario: Merge pull request with options
    Given the GitHub API server:
      """
      require 'json'
      get('/repos/defunkt/hub/pulls/164') { json \
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
    And there is a git FETCH_HEAD
    When I successfully run `hub merge --squash https://github.com/defunkt/hub/pull/164 --no-edit`
    Then "git fetch origin refs/pull/164/head" should be run
    And "git merge --squash --no-edit FETCH_HEAD -m Merge pull request #164 from jfirebaugh/hub_merge" should be run

  Scenario: Merge pull request no repo
    Given the GitHub API server:
      """
      require 'json'
      get('/repos/defunkt/hub/pulls/164') { json \
        :base => {
          :repo => {
            :owner => { :login => "defunkt" },
            :name => "hub",
            :private => false
          }
        },
        :head => {
          :ref => "hub_merge",
          :repo => nil
        },
        :title => "Add `hub merge` command"
      }
      """
    When I run `hub merge https://github.com/defunkt/hub/pull/164`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error: that fork is not available anymore\n
      """
