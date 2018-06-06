Feature: hub delete
  Background:
    Given I am "andreasbaumann" on github.com with OAuth token "OTOKEN"

  Scenario: No argument in current repo
    Given I am in "git://github.com/github/hub.git" git repo
    When I run `hub delete`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Usage: hub delete [-y] [<ORGANIZATION>/]<NAME>\n
      """

  Scenario: Successful confirmation
    Given the GitHub API server:
      """
      delete('/repos/andreasbaumann/my-repo') {
        status 204
      }
      """
    When I run `hub delete my-repo` interactively
    And I type "yes"
    Then the exit status should be 0
    And the output should contain:
      """
      Really delete repository 'andreasbaumann/my-repo' (yes/N)?
      """
    And the output should contain:
      """
      Deleted repository 'andreasbaumann/my-repo'.
      """

  Scenario: Org repo
    Given the GitHub API server:
      """
      delete('/repos/our-org/my-repo') {
        status 204
      }
      """
    When I run `hub delete our-org/my-repo` interactively
    And I type "yes"
    Then the exit status should be 0
    And the output should contain:
      """
      Really delete repository 'our-org/my-repo' (yes/N)?
      """
    And the output should contain:
      """
      Deleted repository 'our-org/my-repo'.
      """

  Scenario: Invalid confirmation
    When I run `hub delete my-repo` interactively
    And I type "y"
    Then the exit status should be 1
    And the output should contain:
      """
      Really delete repository 'andreasbaumann/my-repo' (yes/N)?
      """
    And the stderr should contain exactly:
      """
      Please type 'yes' for confirmation.\n
      """
