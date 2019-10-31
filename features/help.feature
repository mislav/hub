Feature: hub help
  Scenario: Appends hub help to regular help text
    When I successfully run `hub help`
    Then the output should contain:
      """
      These GitHub commands are provided by hub:

         api            Low-level GitHub API request interface
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

  Scenario: Shows help for a hub extension
    When I successfully run `hub help hub-help`
    Then "man hub-help" should be run

  Scenario: Shows help for a hub command
    When I successfully run `hub help fork`
    Then "man hub-fork" should be run

  Scenario: Show help in HTML format
    When I successfully run `hub help -w fork`
    Then "man hub-fork" should not be run
    And "git web--browse PATH/hub-fork.1.html" should be run

  Scenario: Show help in HTML format by default
    Given I successfully run `git config --global help.format html`
    When I successfully run `hub help fork`
    Then "git web--browse PATH/hub-fork.1.html" should be run

  Scenario: Override HTML format back to man
    Given I successfully run `git config --global help.format html`
    When I successfully run `hub help -m fork`
    Then "man hub-fork" should be run

  Scenario: The --help flag opens man page
    When I successfully run `hub fork --help`
    Then "man hub-fork" should be run

  Scenario: The --help flag expands alias first
    Given I successfully run `git config --global alias.ci ci-status`
    When I successfully run `hub ci --help`
    Then "man hub-ci-status" should be run
