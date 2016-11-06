Feature: git-hub compatibility
Scenario: If alias named branch exists, it should not be expanded.
  Given I am in "git://github.com/rtomayko/ronn.git" git repo
  And I am "mislav" on github.com with OAuth token "OTOKEN"
  Given the default branch for "origin" is "develop"
  And I am on the "develop" branch with upstream "origin/develop"
  When I successfully run `git config --global alias.branch "branch -a"`
  When I run `hub branch`
  Then the stdout should contain exactly:
    """
    * develop
      master\n
    """