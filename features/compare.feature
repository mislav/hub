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
    And "open https://github.com/mislav/dotfiles/compare/feature;foo" should be run

  Scenario: No args, no upstream
    When I run `hub compare`
    Then the exit status should be 1
    And the stderr should contain:
      """
      hub compare [USER] [<START>...]<END>
      """

  Scenario: Can't compare default branch to self
    Given the default branch for "origin" is "develop"
    And I am on the "develop" branch with upstream "origin/develop"
    When I run `hub compare`
    Then the exit status should be 1
    And the stderr should contain:
      """
      hub compare [USER] [<START>...]<END>
      """

  Scenario: No args, has upstream branch
    Given I am on the "feature" branch with upstream "origin/experimental"
    And git "push.default" is set to "upstream"
    When I successfully run `hub compare`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles/compare/experimental" should be run

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
