Feature: hub submodule add
  Background:
    Given I am "mislav" on github.com with OAuth token "OTOKEN"
    Given I am in "dotfiles" git repo
    # make existing repo in subdirectory so git clone isn't triggered
    Given a git repo in "vendor/grit"

  Scenario: Add public submodule
    Given the GitHub API server:
      """
      get('/repos/mojombo/grit') {
        json :private => false,
             :name => 'grit', :owner => { :login => 'mojombo' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub submodule add mojombo/grit vendor/grit`
    Then the "vendor/grit" submodule url should be "git://github.com/mojombo/grit.git"
    And the output should contain exactly:
      """
      Adding existing repo at 'vendor/grit' to the index\n
      """

  Scenario: Add private submodule
    Given the GitHub API server:
      """
      get('/repos/mojombo/grit') {
        json :private => false,
             :name => 'grit', :owner => { :login => 'mojombo' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub submodule add -p mojombo/grit vendor/grit`
    Then the "vendor/grit" submodule url should be "git@github.com:mojombo/grit.git"

  Scenario: A submodule for my own repo is public nevertheless
    Given the GitHub API server:
      """
      get('/repos/mislav/grit') {
        json :private => false,
             :name => 'grit', :owner => { :login => 'mislav' },
             :permissions => { :push => true }
      }
      """
    When I successfully run `hub submodule add grit vendor/grit`
    Then the "vendor/grit" submodule url should be "git://github.com/mislav/grit.git"

  Scenario: Add submodule with arguments
    Given the GitHub API server:
      """
      get('/repos/mojombo/grit') {
        json :private => false,
             :name => 'grit', :owner => { :login => 'mojombo' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub submodule add -b foo --name grit mojombo/grit vendor/grit`
    Then "git submodule add -b foo --name grit git://github.com/mojombo/grit.git vendor/grit" should be run

  Scenario: Add submodule with branch
    Given the GitHub API server:
      """
      get('/repos/mojombo/grit') {
        json :private => false,
             :name => 'grit', :owner => { :login => 'mojombo' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub submodule add --branch foo mojombo/grit vendor/grit`
    Then "git submodule add --branch foo git://github.com/mojombo/grit.git vendor/grit" should be run
