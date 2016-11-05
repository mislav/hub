Feature: git-hub compatibility
  Scenario: If alias named branch exists, it should not be expanded.
    Given I am in "git://github.com/rtomayko/ronn.git" git repo
    And the default branch for "origin" is "master"
    When I successfully run `git config --global alias.branch "branch -a"`
    When I run `hub branch`
    Then the stdout should contain exactly "* master\n"
