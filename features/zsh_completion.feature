@completion
Feature: zsh tab-completion

  Background:
    Given my shell is zsh
    And I'm using zsh-distributed base git completions

  Scenario: "pu" expands to "pull-request" after "pull"
    When I type "git pu" and press <Tab>
    Then the completion menu should offer "pull-request" with description "open a pull request on GitHub"
    When I press <Tab> again
    Then the command should expand to "git pull"
    When I press <Tab> again
    Then the command should expand to "git pull-request"

  Scenario: "ci-" expands to "ci-status"
    When I type "git ci-" and press <Tab>
    Then the command should expand to "git ci-status"

  Scenario: Completion of pull-request arguments
    When I type "git pull-request -" and press <Tab>
    Then the completion menu should offer:
      | -b | base                                 |
      | -h | head                                 |
      | -m | message                              |
      | -F | file                                 |
      | -i | issue                                |
      | -f | force (skip check for local commits) |
      | -a | user                                 |
      | -M | milestone                            |
      | -l | labels                               |

  Scenario: Completion of fork arguments
    When I type "git fork -" and press <Tab>
    Then the command should expand to "git fork --no-remote"

  Scenario: Completion of 2nd browse argument
    When I type "git browse -- i" and press <Tab>
    Then the command should expand to "git browse -- issues"

  # In this combination, zsh uses completion support from a bash script.
  Scenario: "ci-" expands to "ci-status"
    Given I'm using git-distributed base git completions
    When I type "git ci-" and press <Tab>
    Then the command should expand to "git ci-status"
