Feature: hub fetch
  Background:
    Given I am in "dotfiles" git repo
    And the "origin" remote has url "git://github.com/evilchelu/dotfiles.git"
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch existing remote
    When I successfully run `hub fetch origin`
    Then the git command should be unchanged
    And the output should not contain anything

  Scenario: Fetch existing remote from non-GitHub source
    Given the "origin" remote has url "ssh://dev@codeserver.dev.xxx.drush.in/~/repository.git"
    When I successfully run `hub fetch origin`
    Then the git command should be unchanged
    And the output should not contain anything

  Scenario: Fetch from non-GitHub source via refspec
    Given the "origin" remote has url "ssh://dev@codeserver.dev.xxx.drush.in/~/repository.git"
    When I successfully run `hub fetch ssh://myusername@a.specific.server:1234/existing-project/gerrit-project-name refs/changes/16/6116/1`
    Then the git command should be unchanged
    And the output should not contain anything

  Scenario: Fetch from local bundle
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :permissions => { :push => false }
      }
      """
    And a git bundle named "mislav"
    When I successfully run `hub fetch mislav`
    Then the git command should be unchanged
    And the output should not contain anything
    And there should be no "mislav" remote

  Scenario: Creates new remote
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub fetch mislav`
    Then "git fetch mislav" should be run
    And the url for "mislav" should be "https://github.com/mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Creates new remote with git
    Given git protocol is preferred
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub fetch mislav`
    Then "git fetch mislav" should be run
    And the url for "mislav" should be "git://github.com/mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Owner name with dash
    Given the GitHub API server:
      """
      get('/repos/ankit-maverick/dotfiles') {
        json :private => false,
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub fetch ankit-maverick`
    Then "git fetch ankit-maverick" should be run
    And the url for "ankit-maverick" should be "https://github.com/ankit-maverick/dotfiles.git"
    And the output should not contain anything

  Scenario: Private repo
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => true,
             :permissions => { :push => false }
      }
      """
    And git protocol is preferred
    When I successfully run `hub fetch mislav`
    Then "git fetch mislav" should be run
    And the url for "mislav" should be "git@github.com:mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Writeable repo
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :permissions => { :push => true }
      }
      """
    And git protocol is preferred
    When I successfully run `hub fetch mislav`
    Then "git fetch mislav" should be run
    And the url for "mislav" should be "git@github.com:mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Fetch with options
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub fetch --depth=1 mislav`
    Then "git fetch --depth=1 mislav" should be run

  Scenario: Fetch multiple
    Given the GitHub API server:
      """
      get('/repos/:owner/dotfiles') {
        json :private => false,
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub fetch --multiple mislav rtomayko`
    Then "git fetch --multiple mislav rtomayko" should be run
    And the url for "mislav" should be "https://github.com/mislav/dotfiles.git"
    And the url for "rtomayko" should be "https://github.com/rtomayko/dotfiles.git"

  Scenario: Fetch multiple with filtering
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :permissions => { :push => false }
      }
      """
    When I successfully run `git config remotes.mygrp "foo bar"`
    When I successfully run `hub fetch --multiple origin mislav mygrp https://example.com typo`
    Then "git fetch --multiple origin mislav mygrp https://example.com typo" should be run
    And the url for "mislav" should be "https://github.com/mislav/dotfiles.git"
    But there should be no "mygrp" remote
    And there should be no "typo" remote

  Scenario: Fetch multiple comma-separated
    Given the GitHub API server:
      """
      get('/repos/:owner/dotfiles') {
        json :private => false,
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub fetch mislav,rtomayko,dustinleblanc`
    Then "git fetch --multiple mislav rtomayko dustinleblanc" should be run
    And the url for "mislav" should be "https://github.com/mislav/dotfiles.git"
    And the url for "rtomayko" should be "https://github.com/rtomayko/dotfiles.git"
    And the url for "dustinleblanc" should be "https://github.com/dustinleblanc/dotfiles.git"

  Scenario: Doesn't create a new remote if repo doesn't exist on GitHub
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') { status 404 }
      """
    When I successfully run `hub fetch mislav`
    Then the git command should be unchanged
    And there should be no "mislav" remote
