Feature: git-hub compatibility
  Scenario: If alias named branch exists, it should not be expanded.
    Given I am in "git://github.com/rtomayko/ronn.git" git repo
    And the default branch for "origin" is "master"
    When I successfully run `git config --global alias.branch "branch -a"`
    When I run `hub branch`
    Then the stdout should contain exactly "* master\n"

  Scenario: List commands
    When I successfully run `hub --list-cmds=others`
    Then the stdout should contain exactly:
      """
      add
      branch
      commit
      alias
      api
      browse
      ci-status
      compare
      create
      delete
      fork
      gist
      issue
      pr
      pull-request
      release
      sync\n
      """

  Scenario: Doesn't sabotage --exec-path
    When I successfully run `hub --exec-path`
    Then the output should not contain "These GitHub commands"

  Scenario: Shows help with --git-dir
    When I run `hub --git-dir=.git`
    Then the exit status should be 1
    And the output should contain "usage: git "
