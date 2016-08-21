Feature: hub compare
  Background:
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Compare branch
    When I successfully run `hub compare refactor`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/refactor" should be run

  Scenario: Compare complex branch
    When I successfully run `hub compare feature/foo`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/feature/foo" should be run

  Scenario: Compare branch with funky characters
    When I successfully run `hub compare 'my#branch!with.special+chars'`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/my%23branch!with.special%2Bchars" should be run

  Scenario: No args, no upstream
    When I run `hub compare`
    Then the exit status should be 1
    And the stderr should contain:
      """
      Usage: hub compare [-u] [-b <BASE>] [<USER>] [[<START>...]<END>]
      """

  Scenario: Can't compare default branch to self
    Given the default branch for "origin" is "develop"
    And I am on the "develop" branch with upstream "origin/develop"
    When I run `hub compare`
    Then the exit status should be 1
    And the stderr should contain:
      """
      Usage: hub compare [-u] [-b <BASE>] [<USER>] [[<START>...]<END>]
      """

  Scenario: No args, has upstream branch
    Given I am on the "feature" branch with upstream "origin/experimental"
    And git "push.default" is set to "upstream"
    When I successfully run `hub compare`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/experimental" should be run

  Scenario: Current branch has funky characters
    Given I am on the "feature" branch with upstream "origin/my#branch!with.special+chars"
    And git "push.default" is set to "upstream"
    When I successfully run `hub compare`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/my%23branch!with.special%2Bchars" should be run

  Scenario: Compare range
    When I successfully run `hub compare 1.0...fix`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/1.0...fix" should be run

  Scenario: Output URL without opening the browser
    When I successfully run `hub compare -u 1.0...fix`
    Then "open https://github.com/mislav/dotfiles/compare/1.0...fix" should not be run
    And the stdout should contain exactly:
      """
      https://github.com/mislav/dotfiles/compare/1.0...fix\n
      """

  Scenario: Compare base in branch that is not master
    Given I am on the "feature" branch with upstream "origin/experimental"
    And git "push.default" is set to "upstream"
    When I successfully run `hub compare -b master`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/master...experimental" should be run

  Scenario: Compare base in master branch
    Given I am on the "master" branch with upstream "origin/master"
    And git "push.default" is set to "upstream"
    When I successfully run `hub compare -b experimental`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/experimental...master" should be run

  Scenario: Compare base with same branch as the current branch
    Given I am on the "feature" branch with upstream "origin/experimental"
    And git "push.default" is set to "upstream"
    When I run `hub compare -b experimental`
    Then "open https://github.com/mislav/dotfiles/compare/experimental...experimental" should not be run
    And the exit status should be 1
    And the stderr should contain:
      """
      Usage: hub compare [-u] [-b <BASE>] [<USER>] [[<START>...]<END>]
      """

  Scenario: Compare base with parameters
    Given I am on the "master" branch with upstream "origin/master"
    When I run `hub compare -b master experimental..master`
    Then "open https://github.com/mislav/dotfiles/compare/experimental...master" should not be run
    And the exit status should be 1
    And the stderr should contain:
      """
      Usage: hub compare [-u] [-b <BASE>] [<USER>] [[<START>...]<END>]
      """

  Scenario: Compare 2-dots range for tags
    When I successfully run `hub compare 1.0..fix`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/1.0...fix" should be run

  Scenario: Compare 2-dots range for SHAs
    When I successfully run `hub compare 1234abc..3456cde`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/1234abc...3456cde" should be run

  Scenario: Compare 2-dots range with "user:repo" notation
    When I successfully run `hub compare henrahmagix:master..2b10927`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/henrahmagix:master...2b10927" should be run

  Scenario: Complex range is unchanged
    When I successfully run `hub compare @{a..b}..@{c..d}`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/@{a..b}..@{c..d}" should be run

  Scenario: Compare wiki
    Given the "origin" remote has url "git://github.com/mislav/dotfiles.wiki.git"
    When I successfully run `hub compare 1.0..fix`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/wiki/_compare/1.0...fix" should be run

  Scenario: Compare fork
    When I successfully run `hub compare anotheruser feature`
    Then there should be no output
    And "open https://github.com/anotheruser/dotfiles/compare/feature" should be run

  Scenario: Enterprise repo over HTTP
    Given the "origin" remote has url "git://git.my.org/mislav/dotfiles.git"
    And I am "mislav" on http://git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    When I successfully run `hub compare refactor`
    Then there should be no output
    And "open http://git.my.org/mislav/dotfiles/compare/refactor" should be run

  Scenario: Compare in non-GitHub repo
    Given the "origin" remote has url "git@bitbucket.org:mislav/dotfiles.git"
    And I am on the "feature" branch
    When I run `hub compare`
    Then the stdout should contain exactly ""
    And the stderr should contain exactly:
      """
      Aborted: the origin remote doesn't point to a GitHub repository.\n
      """
    And the exit status should be 1
