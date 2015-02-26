Feature: hub merge
  Background:
    Given I am in "hub" git repo
    And the "origin" remote has url "git://github.com/defunkt/hub.git"
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Merge pull request
    Given the GitHub API server:
      """
      get('/repos/defunkt/hub/pulls/164') { json \
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
    And there is a commit named "jfirebaugh/hub_merge"
    When I successfully run `hub merge https://github.com/defunkt/hub/pull/164`
    Then "git fetch git://github.com/jfirebaugh/hub.git +refs/heads/hub_merge:refs/remotes/jfirebaugh/hub_merge" should be run
    And "git merge jfirebaugh/hub_merge --no-ff -m Merge pull request #164 from jfirebaugh/hub_merge" should be run
    When I successfully run `git show -s --format=%B`
    Then the output should contain:
      """
      Merge pull request #164 from jfirebaugh/hub_merge

      Add `hub merge` command
      """

  Scenario: Merge pull request with --ff-only option
    Given the GitHub API server:
      """
      require 'json'
      get('/repos/defunkt/hub/pulls/164') { json \
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
    And there is a commit named "jfirebaugh/hub_merge"
    When I successfully run `hub merge --ff-only https://github.com/defunkt/hub/pull/164`
    Then "git fetch git://github.com/jfirebaugh/hub.git +refs/heads/hub_merge:refs/remotes/jfirebaugh/hub_merge" should be run
    And "git merge --ff-only jfirebaugh/hub_merge -m Merge pull request #164 from jfirebaugh/hub_merge" should be run
    When I successfully run `git show -s --format=%B`
    Then the output should contain:
      """
      Fast-forward (no commit created; -m option ignored)
      """

  Scenario: Merge pull request with --squash option
    Given the GitHub API server:
      """
      require 'json'
      get('/repos/defunkt/hub/pulls/164') { json \
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
    And there is a commit named "jfirebaugh/hub_merge"
    When I successfully run `hub merge --squash https://github.com/defunkt/hub/pull/164`
    Then "git fetch git://github.com/jfirebaugh/hub.git +refs/heads/hub_merge:refs/remotes/jfirebaugh/hub_merge" should be run
    And "git merge --squash jfirebaugh/hub_merge -m Merge pull request #164 from jfirebaugh/hub_merge" should be run
    When I successfully run `git show -s --format=%B`
    Then the output should contain:
      """
      Fast-forward (no commit created; -m option ignored)
      """

  Scenario: Merge private pull request
    Given the GitHub API server:
      """
      require 'json'
      get('/repos/defunkt/hub/pulls/164') { json \
        :head => {
          :ref => "hub_merge",
          :repo => {
            :owner => { :login => "jfirebaugh" },
            :name => "hub",
            :private => true
          }
        },
        :title => "Add `hub merge` command"
      }
      """
    And there is a commit named "jfirebaugh/hub_merge"
    When I successfully run `hub merge https://github.com/defunkt/hub/pull/164`
    Then "git fetch git@github.com:jfirebaugh/hub.git +refs/heads/hub_merge:refs/remotes/jfirebaugh/hub_merge" should be run

  Scenario: Missing repo
    Given the GitHub API server:
      """
      require 'json'
      get('/repos/defunkt/hub/pulls/164') { json \
        :head => {
          :ref => "hub_merge",
          :repo => nil
        }
      }
      """
    When I run `hub merge https://github.com/defunkt/hub/pull/164`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error: that fork is not available anymore\n
      """

  Scenario: Renamed repo
    Given the GitHub API server:
      """
      require 'json'
      get('/repos/defunkt/hub/pulls/164') { json \
        :head => {
          :ref => "hub_merge",
          :repo => {
            :owner => { :login => "jfirebaugh" },
            :name => "hub-1",
            :private => false
          }
        }
      }
      """
    And there is a commit named "jfirebaugh/hub_merge"
    When I successfully run `hub merge https://github.com/defunkt/hub/pull/164`
    Then "git fetch git://github.com/jfirebaugh/hub-1.git +refs/heads/hub_merge:refs/remotes/jfirebaugh/hub_merge" should be run

  Scenario: Unchanged merge
    When I run `hub merge master`
    Then "git merge master" should be run
