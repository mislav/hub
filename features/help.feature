Feature: hub help
  Scenario: Appends hub help to regular help text
    When I successfully run `hub help`
    Then the output should contain:
      """
      These GitHub commands are provided by hub:

         pull-request   Open a pull request on GitHub
      """
    And the output should contain "usage: git "

  Scenario: Appends hub commands to `--all` output
    When I successfully run `hub help -a`
    Then the output should contain "pull-request"

  Scenario: Shows help for a subcommand
    When I successfully run `hub help hub-help`
    Then the output should contain "Usage: hub help"
