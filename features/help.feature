Feature: hub help
  Scenario: Appends hub help to regular help text
    When I successfully run `hub help`
    Then the output should contain:
      """
      These GitHub commands are provided by hub:

         browse         Open a GitHub page in the default browser
      """
    And the output should contain "usage: git "

  Scenario: Appends hub commands to `--all` output
    When I successfully run `hub help -a`
    Then the output should contain "pull-request"

  Scenario: Shows help for a subcommand
    When I successfully run `hub help hub-help`
    Then the output should contain "Usage: hub help"

  Scenario: Doesn't sabotage --exec-path
    When I successfully run `hub --exec-path`
    Then the output should not contain "These GitHub commands"
