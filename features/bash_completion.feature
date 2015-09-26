@completion
Feature: bash tab-completion

  Background:
    Given my shell is bash
    And I'm using git-distributed base git completions

  Scenario: "pu" matches multiple commands including "pull-request"
    When I type "git pu" and press <Tab>
    Then the command should not expand
    When I press <Tab> again
    Then the completion menu should offer "pull pull-request push"

  Scenario: "ci-" expands to "ci-status"
    When I type "git ci-" and press <Tab>
    Then the command should expand to "git ci-status"

  Scenario: Offers pull-request flags
    When I type "git pull-request -" and press <Tab>
    When I press <Tab> again
    Then the completion menu should offer "-F -b -f -h -i -m -a -M -l" unsorted

  Scenario: Doesn't offer already used pull-request flags
    When I type "git pull-request -F myfile -h mybranch -" and press <Tab>
    When I press <Tab> again
    Then the completion menu should offer "-b -f -i -m -a -M -l" unsorted

  Scenario: Browse to issues
    When I type "git browse -- i" and press <Tab>
    Then the command should expand to "git browse -- issues"

  Scenario: Browse to punch-card graph
    When I type "git browse -- graphs/p" and press <Tab>
    Then the command should expand to "git browse -- graphs/punch-card"

  Scenario: Completion of fork argument
    When I type "git fork -" and press <Tab>
    Then the command should expand to "git fork --no-remote"

  Scenario: Completion of user/repo in "browse"
  Scenario: Completion of branch names in "compare"
  Scenario: Completion of "owner/repo:branch" in "pull-request -h/b"
