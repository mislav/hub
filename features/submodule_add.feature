Feature: hub submodule add
  Background:
    Given I am in "dotfiles" git repo
    # make existing repo in subdirectory so git clone isn't triggered
    Given a git repo in "vendor/grit"

  Scenario: Add public submodule
    When I successfully run `hub submodule add mojombo/grit vendor/grit`
    Then the "vendor/grit" submodule url should be "git://github.com/mojombo/grit.git"
    And the output should contain exactly:
      """
      Adding existing repo at 'vendor/grit' to the index\n
      """

  Scenario: Add private submodule
    When I successfully run `hub submodule add -p mojombo/grit vendor/grit`
    Then the "vendor/grit" submodule url should be "git@github.com:mojombo/grit.git"

  Scenario: Add submodule with arguments
    When I successfully run `hub submodule add -b foo --name grit mojombo/grit vendor/grit`
    Then "git submodule add -b foo --name grit git://github.com/mojombo/grit.git vendor/grit" should be run

  Scenario: Add submodule with branch
    When I successfully run `hub submodule add --branch foo mojombo/grit vendor/grit`
    Then "git submodule add --branch foo git://github.com/mojombo/grit.git vendor/grit" should be run
