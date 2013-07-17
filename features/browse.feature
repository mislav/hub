Feature: hub browse
  Scenario: Project with owner
    When I successfully run `hub browse mislav/dotfiles`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles" should be run

  Scenario: Project without owner
    Given I am "mislav" on github.com
    When I successfully run `hub browse dotfiles`
    Then "open https://github.com/mislav/dotfiles" should be run

  Scenario: Explicit project overrides current
    Given I am in "git://github.com/josh/rails-behaviors.git" git repo
    And I am "mislav" on github.com
    When I successfully run `hub browse dotfiles`
    Then "open https://github.com/mislav/dotfiles" should be run

  Scenario: Project issues
    When I successfully run `hub browse mislav/dotfiles issues`
    Then "open https://github.com/mislav/dotfiles/issues" should be run

  Scenario: Project wiki
    When I successfully run `hub browse mislav/dotfiles wiki`
    Then "open https://github.com/mislav/dotfiles/wiki" should be run

  Scenario: Project commits on master
    When I successfully run `hub browse mislav/dotfiles commits`
    Then "open https://github.com/mislav/dotfiles/commits/master" should be run

  Scenario: Specific commit in project
    When I successfully run `hub browse mislav/dotfiles commit/4173c3b`
    Then "open https://github.com/mislav/dotfiles/commit/4173c3b" should be run

  Scenario: Output the URL instead of browse
    When I successfully run `hub browse -u mislav/dotfiles`
    Then the output should contain exactly "https://github.com/mislav/dotfiles\n"
    But "open https://github.com/mislav/dotfiles" should not be run

  Scenario: Current project
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    When I successfully run `hub browse`
    Then there should be no output
    And "open https://github.com/mislav/dotfiles" should be run

  Scenario: Commit in current project
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    When I successfully run `hub browse -- commit/abcd1234`
    Then "open https://github.com/mislav/dotfiles/commit/abcd1234" should be run

  Scenario: Current branch
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am on the "feature" branch with upstream "origin/experimental"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles/tree/experimental" should be run

  Scenario: Default branch
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And the default branch for "origin" is "develop"
    And I am on the "develop" branch with upstream "origin/develop"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles" should be run

  Scenario: Current branch, no tracking
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am on the "feature" branch
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles" should be run

  Scenario: Current branch with special chars
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am on the "fix-bug-#123" branch with upstream "origin/fix-bug-#123"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles/tree/fix-bug-%23123" should be run

  Scenario: Commits on current branch
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am on the "feature" branch with upstream "origin/experimental"
    When I successfully run `hub browse -- commits`
    Then "open https://github.com/mislav/dotfiles/commits/experimental" should be run

  Scenario: Complex branch
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am on the "foo/bar" branch with upstream "origin/baz/qux/moo"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles/tree/baz/qux/moo" should be run

  Scenario: Wiki repo
    Given I am in "git://github.com/defunkt/hub.wiki.git" git repo
    When I successfully run `hub browse`
    Then "open https://github.com/defunkt/hub/wiki" should be run

  Scenario: Wiki commits
    Given I am in "git://github.com/defunkt/hub.wiki.git" git repo
    When I successfully run `hub browse -- commits`
    Then "open https://github.com/defunkt/hub/wiki/_history" should be run

  Scenario: Wiki pages
    Given I am in "git://github.com/defunkt/hub.wiki.git" git repo
    When I successfully run `hub browse -- pages`
    Then "open https://github.com/defunkt/hub/wiki/_pages" should be run

  Scenario: Deprecated -p flag
    When I successfully run `hub browse -p defunkt/hub`
    Then the stderr should contain exactly:
      """
      Warning: the `-p` flag has no effect anymore\n
      """
    But "open https://github.com/defunkt/hub" should be run
