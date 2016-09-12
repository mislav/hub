Feature: hub cherry-pick
  Background:
    Given I am in "git://github.com/rtomayko/ronn.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Unchanged
    When I run `hub cherry-pick a319d88`
    Then the git command should be unchanged

  Scenario: From GitHub commit URL
    When I run `hub cherry-pick https://github.com/rtomayko/ronn/commit/a319d88#comments`
    Then "git fetch origin" should be run
    And "git cherry-pick a319d88" should be run

  Scenario: From GitHub pull request URL
    Given the GitHub API server:
      """
      get('/repos/blueyed/ronn/pulls/560') {
        json :head => {
          :ref => "fixes",
          :repo => {
            :owner => { :login => "mislav" },
            :name => "ronin",
            :private => true
          }
        }
      }
      """
    When I run `hub cherry-pick https://github.com/blueyed/ronn/pull/560/commits/a319d88`
    Then "git remote add _hub-cherry-pick git@github.com:mislav/ronin.git" should be run
    And "git fetch -q --no-tags _hub-cherry-pick" should be run
    And "git remote rm _hub-cherry-pick" should be run
    And "git cherry-pick a319d88" should be run

  Scenario: From fork that has existing remote
    Given the "mislav" remote has url "git@github.com:mislav/ronn.git"
    When I run `hub cherry-pick https://github.com/mislav/ronn/commit/a319d88`
    Then "git fetch mislav" should be run
    And "git cherry-pick a319d88" should be run

  Scenario: Using GitHub owner@SHA notation
    Given the "mislav" remote has url "git@github.com:mislav/ronn.git"
    When I run `hub cherry-pick mislav@a319d88`
    Then "git fetch mislav" should be run
    And "git cherry-pick a319d88" should be run

  Scenario: Using GitHub owner@SHA notation that is too short
    When I run `hub cherry-pick mislav@a319`
    Then the git command should be unchanged

  Scenario: Unsupported GitHub owner/repo@SHA notation
    When I run `hub cherry-pick mislav/ronn@a319d88`
    Then the git command should be unchanged

  Scenario: Skips processing if `-m/--mainline` is specified
    When I run `hub cherry-pick -m 42 mislav@a319d88`
    Then the git command should be unchanged
    When I run `hub cherry-pick --mainline 42 mislav@a319d88`
    Then the git command should be unchanged

  Scenario: Using GitHub owner@SHA notation with remote add
    When I run `hub cherry-pick mislav@a319d88`
    Then "git remote add _hub-cherry-pick git://github.com/mislav/ronn.git" should be run
    And "git fetch -q --no-tags _hub-cherry-pick" should be run
    And "git remote rm _hub-cherry-pick" should be run
    And "git cherry-pick a319d88" should be run

  Scenario: From fork that doesn't have a remote
    When I run `hub cherry-pick https://github.com/jingweno/ronn/commit/a319d88`
    Then "git remote add _hub-cherry-pick git://github.com/jingweno/ronn.git" should be run
    And "git fetch -q --no-tags _hub-cherry-pick" should be run
    And "git remote rm _hub-cherry-pick" should be run
    And "git cherry-pick a319d88" should be run
