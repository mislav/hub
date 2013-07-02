@completion
Feature: bash tab-completion

  Scenario: "pu" matches multiple commands including "pull-request"
    Given my shell is bash
    And I'm using git-distributed base git completions
    When I type "git pu" and press <Tab>
    Then the command should not expand
    When I press <Tab> again
    Then the completion menu should offer "pull pull-request push"

  Scenario: "ci-" expands to "ci-status"
    Given my shell is bash
    And I'm using git-distributed base git completions
    When I type "git ci-" and press <Tab>
    Then the command should expand to "git ci-status"

  # In this combination, zsh uses completion support from a bash script.
  Scenario: "ci-" expands to "ci-status"
    Given my shell is zsh
    And I'm using git-distributed base git completions
    When I type "git ci-" and press <Tab>
    Then the command should expand to "git ci-status"
