Feature: hub browse
  Background:
    Given I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: No repo
    When I run `hub browse`
    Then the exit status should be 1
    Then the output should contain exactly "Usage: hub browse [-uc] [[<USER>/]<REPOSITORY>|--] [<SUBPAGE>]\n"

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
    And git "push.default" is set to "upstream"
    And I am on the "feature" branch with upstream "origin/experimental"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles/tree/experimental" should be run

  Scenario: Current branch pushed to fork
    Given I am in "git://github.com/blueyed/dotfiles.git" git repo
    And the "mislav" remote has url "git@github.com:mislav/dotfiles.git"
    And I am on the "feature" branch with upstream "mislav/experimental"
    And git "push.default" is set to "upstream"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles/tree/experimental" should be run

  Scenario: Current branch pushed to fork with simple tracking
    Given I am in "git://github.com/blueyed/dotfiles.git" git repo
    And the "mislav" remote has url "git@github.com:mislav/dotfiles.git"
    And I am on the "feature" branch with upstream "mislav/feature"
    And git "push.default" is set to "simple"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles/tree/feature" should be run

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

  Scenario: Default branch in upstream repo as opposed to fork
    Given I am in "git://github.com/jashkenas/coffee-script.git" git repo
    And the "mislav" remote has url "git@github.com:mislav/coffee-script.git"
    And the default branch for "origin" is "master"
    And the "master" branch is pushed to "mislav/master"
    When I successfully run `hub browse`
    Then "open https://github.com/jashkenas/coffee-script" should be run

  Scenario: Current branch with special chars
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am on the "fix-bug-#123" branch with upstream "origin/fix-bug-#123"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles/tree/fix-bug-%23123" should be run

  Scenario: Commits on current branch
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And git "push.default" is set to "upstream"
    And I am on the "feature" branch with upstream "origin/experimental"
    When I successfully run `hub browse -- commits`
    Then "open https://github.com/mislav/dotfiles/commits/experimental" should be run

  Scenario: Issues subpage ignores tracking configuration
    Given I am in "git://github.com/jashkenas/coffee-script.git" git repo
    And the "mislav" remote has url "git@github.com:mislav/coffee-script.git"
    And git "push.default" is set to "upstream"
    And I am on the "feature" branch with upstream "mislav/experimental"
    When I successfully run `hub browse -- issues`
    Then "open https://github.com/jashkenas/coffee-script/issues" should be run

  Scenario: Issues subpage ignores current branch
    Given I am in "git://github.com/jashkenas/coffee-script.git" git repo
    And the "mislav" remote has url "git@github.com:mislav/coffee-script.git"
    And I am on the "feature" branch pushed to "mislav/feature"
    When I successfully run `hub browse -- issues`
    Then there should be no output
    Then "open https://github.com/jashkenas/coffee-script/issues" should be run

  Scenario: Forward Slash Delimited branch
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And git "push.default" is set to "upstream"
    And I am on the "foo/bar" branch with upstream "origin/baz/qux/moo"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles/tree/baz/qux/moo" should be run

  Scenario: No branch
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am in detached HEAD
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles" should be run

  Scenario: No branch to pulls
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am in detached HEAD
    When I successfully run `hub browse -- pulls`
    Then "open https://github.com/mislav/dotfiles/pulls" should be run

  Scenario: Dot Delimited branch
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And git "push.default" is set to "upstream"
    And I am on the "fix-glob-for.js" branch with upstream "origin/fix-glob-for.js"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles/tree/fix-glob-for.js" should be run

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

  Scenario: Repo with remote with local path
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And the "upstream" remote has url "../path/to/another/repo.git"
    When I successfully run `hub browse`
    Then "open https://github.com/mislav/dotfiles" should be run

  Scenario: Enterprise repo
    Given I am in "git://git.my.org/mislav/dotfiles.git" git repo
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    When I successfully run `hub browse`
    Then "open https://git.my.org/mislav/dotfiles" should be run

  Scenario: Multiple Enterprise repos
    Given I am in "git://git.my.org/mislav/dotfiles.git" git repo
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    And "git.another.org" is a whitelisted Enterprise host
    When I successfully run `hub browse`
    Then "open https://git.my.org/mislav/dotfiles" should be run

  Scenario: Enterprise repo over HTTP
    Given I am in "git://git.my.org/mislav/dotfiles.git" git repo
    And I am "mislav" on http://git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    When I successfully run `hub browse`
    Then "open http://git.my.org/mislav/dotfiles" should be run

  Scenario: SSH alias
    Given the SSH config:
      """
      Host gh
        User git
        HostName github.com
      """
    Given I am in "gh:singingwolfboy/sekrit.git" git repo
    When I successfully run `hub browse`
    Then "open https://github.com/singingwolfboy/sekrit" should be run

  Scenario: SSH GitHub alias
    Given the SSH config:
      """
      Host github.com
        HostName ssh.github.com
      """
    Given I am in "git@github.com:suan/git-sanity.git" git repo
    When I successfully run `hub browse`
    Then "open https://github.com/suan/git-sanity" should be run
