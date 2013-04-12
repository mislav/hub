Feature: hub collaborator
  Background:
    Given I am in "dotfiles" git repo
    And the "origin" remote has url "git://github.com/mislav/dotfiles.git"
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: No args
    When I run `hub collaborator`
    Then the exit status should be 1
    And the stderr should contain:
      """
      hub collaborator [ACTION] <[USER1] [USER2]...>
      """

  Scenario: hub collaborator list
    When I successfully run `hub collaborator list`
    Then the stdout should contain exactly:
    """
    Some folks
    """

  Scenario: hub collaborator add with no user
    When I run `hub collaborator add`
    Then the exit status should be 1
    And the stderr should contain:
      """
      hub collaborator add [USER1] <[USER2] ...>
      """

  Scenario: hub collaborator add with one user

  Scenario: hub collaborator add with many users

  Scenario: hub collaborator remove with no user
    When I run `hub collaborator remove`
    Then the exit status should be 1
    And the stderr should contain:
      """
      hub collaborator remove [USER1] <[USER2] ...>
      """

  Scenario: hub collaborator remove with one user

  Scenario: hub collaborator remove with many users
