Feature: hub compare
  Background:
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Compare branch
    When I successfully run `hub compare refactor`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/refactor" should be run

  Scenario: Compare complex branch
    When I successfully run `hub compare feature/foo`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/feature/foo" should be run

  Scenario: Compare branch with funky characters
    When I successfully run `hub compare 'my#branch!with.special+chars'`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/my%23branch!with.special%2Bchars" should be run

  Scenario: No args, no upstream
    When I run `hub compare`
    Then the exit status should be 1
    And the stderr should contain exactly "the current branch 'master' doesn't seem pushed to a remote\n"

  Scenario: Can't compare default branch to self
    Given the default branch for "origin" is "develop"
    And I am on the "develop" branch with upstream "origin/develop"
    When I run `hub compare`
    Then the exit status should be 1
    And the stderr should contain exactly "the branch to compare 'develop' is the default branch\n"

  Scenario: No args, has upstream branch
    Given I am on the "feature" branch with upstream "origin/experimental"
    And git "push.default" is set to "upstream"
    When I successfully run `hub compare`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/experimental" should be run

  Scenario: Current branch has funky characters
    Given I am on the "feature" branch with upstream "origin/my#branch!with.special+chars"
    And git "push.default" is set to "upstream"
    When I successfully run `hub compare`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/my%23branch!with.special%2Bchars" should be run

  Scenario: Current branch pushed to fork
    Given I am "monalisa" on github.com with OAuth token "MONATOKEN"
    And the "monalisa" remote has url "git@github.com:monalisa/dotfiles.git"
    And I am on the "topic" branch pushed to "monalisa/topic"
    When I successfully run `hub compare`
    Then "open https://github.com/mislav/dotfiles/compare/monalisa:topic" should be run

  Scenario: Current branch with full URL in upstream configuration
    Given I am on the "local-topic" branch
    When I successfully run `git config branch.local-topic.remote https://github.com/monalisa/dotfiles.git`
    When I successfully run `git config branch.local-topic.merge refs/remotes/remote-topic`
    When I successfully run `hub compare`
    Then "open https://github.com/mislav/dotfiles/compare/monalisa:remote-topic" should be run

  Scenario: Compare range
    When I successfully run `hub compare 1.0...fix`
    Then the output should not contain anything
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
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/master...experimental" should be run

  Scenario: Compare base in master branch
    Given I am on the "master" branch with upstream "origin/master"
    And git "push.default" is set to "upstream"
    When I successfully run `hub compare -b experimental`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/experimental...master" should be run

  Scenario: Compare base with same branch as the current branch
    Given I am on the "feature" branch with upstream "origin/experimental"
    And git "push.default" is set to "upstream"
    When I run `hub compare -b experimental`
    Then "open https://github.com/mislav/dotfiles/compare/experimental...experimental" should not be run
    And the exit status should be 1
    And the stderr should contain exactly "the branch to compare 'experimental' is the same as --base\n"

  Scenario: Compare base with parameters
    Given I am on the "master" branch with upstream "origin/master"
    When I run `hub compare -b master experimental..master`
    Then "open https://github.com/mislav/dotfiles/compare/experimental...master" should not be run
    And the exit status should be 1
    And the stderr should contain "Usage: hub compare"

  Scenario: Compare 2-dots range for tags
    When I successfully run `hub compare 1.0..fix`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/1.0...fix" should be run

  Scenario: Compare 2-dots range for SHAs
    When I successfully run `hub compare 1234abc..3456cde`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/1234abc...3456cde" should be run

  Scenario: Compare 2-dots range with "user:repo" notation
    When I successfully run `hub compare henrahmagix:master..2b10927`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/henrahmagix:master...2b10927" should be run

  Scenario: Compare 2-dots range with slashes in branch names
    When I successfully run `hub compare one/foo..two/bar/baz`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/one/foo...two/bar/baz" should be run

  Scenario: Complex range is unchanged
    When I successfully run `hub compare @{a..b}..@{c..d}`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/@{a..b}..@{c..d}" should be run

  Scenario: Compare wiki
    Given the "origin" remote has url "git://github.com/mislav/dotfiles.wiki.git"
    When I successfully run `hub compare 1.0..fix`
    Then the output should not contain anything
    And "open https://github.com/mislav/dotfiles/wiki/_compare/1.0...fix" should be run

  Scenario: Compare fork
    When I successfully run `hub compare anotheruser feature`
    Then the output should not contain anything
    And "open https://github.com/anotheruser/dotfiles/compare/feature" should be run

  Scenario: Enterprise repo over HTTP
    Given the "origin" remote has url "git://git.my.org/mislav/dotfiles.git"
    And I am "mislav" on http://git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    When I successfully run `hub compare refactor`
    Then the output should not contain anything
    And "open http://git.my.org/mislav/dotfiles/compare/refactor" should be run

  Scenario: Enterprise repo with explicit upstream project
    Given the "origin" remote has url "git://git.my.org/mislav/dotfiles.git"
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    When I successfully run `hub compare fehmicansaglam a..b`
    Then the output should not contain anything
    And "open https://git.my.org/fehmicansaglam/dotfiles/compare/a...b" should be run

  Scenario: Compare in non-GitHub repo
    Given the "origin" remote has url "git@bitbucket.org:mislav/dotfiles.git"
    And I am on the "feature" branch
    When I run `hub compare`
    Then the stdout should contain exactly ""
    And the stderr should contain exactly:
      """
      Aborted: could not find any git remote pointing to a GitHub repository\n
      """
    And the exit status should be 1

  Scenario: Comparing two branches while not on a local branch
    Given I am in detached HEAD
    And I run `hub compare refactor...master`
    Then the exit status should be 0
    And the output should not contain anything
    And "open https://github.com/mislav/dotfiles/compare/refactor...master" should be run
