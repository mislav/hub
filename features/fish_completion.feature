@completion
Feature: fish tab-completion

  Background:
    Given my shell is fish

  Scenario: "pu" matches multiple commands including "pull-request"
    When I type "git pu" and press <Tab>
    Then the command should not expand
    When I press <Tab> again
    Then the completion menu should offer "pull push pull-request" unsorted

  Scenario: "ci-" expands to "ci-status"
    When I type "git ci-" and press <Tab>
    Then the command should expand to "git ci-status"

  Scenario: Offers pull-request flags
    When I type "git pull-request -" and press <Tab>
    When I press <Tab> again
    Then the completion menu should offer "-F -b -f -h -m -a -M -l -o --browse -p --help" unsorted

  Scenario: Browse to issues
    When I type "git browse -- i" and press <Tab>
    Then the command should expand to "git browse -- issues"
