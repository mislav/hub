Feature: hub help
  Scenario: Appends hub help to regular help text
    When I successfully run `hub help`
    Then the output should contain:
      """
      These GitHub commands are provided by hub:

         browse         Open a GitHub page in the default browser
      """
    And the output should contain "usage: git "

  Scenario: Shows help text with no arguments
    When I run `hub`
    Then the stdout should contain "usage: git "
    And the stderr should contain exactly ""
    And the exit status should be 1

  Scenario: Appends hub commands to `--all` output
    When I successfully run `hub help -a`
    Then the output should contain "pull-request"

  Scenario: Shows help for a subcommand
    When I successfully run `hub help hub-help`
    Then the output should contain "`hub help` hub-<COMMAND>"
