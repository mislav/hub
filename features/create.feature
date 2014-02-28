Feature: hub create
  Background:
    Given I am in "dotfiles" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Create repo
    Given the GitHub API server:
      """
      post('/user/repos') {
        assert :private => false
        json :full_name => 'mislav/dotfiles'
      }
      """
    When I successfully run `hub create`
    Then the url for "origin" should be "git@github.com:mislav/dotfiles.git"
    And the output should contain exactly "created repository: mislav/dotfiles\n"

  Scenario: Create private repo
    Given the GitHub API server:
      """
      post('/user/repos') {
        assert :private => true
        json :full_name => 'mislav/dotfiles'
      }
      """
    When I successfully run `hub create -p`
    Then the url for "origin" should be "git@github.com:mislav/dotfiles.git"

  Scenario: HTTPS is preferred
    Given the GitHub API server:
      """
      post('/user/repos') {
        json :full_name => 'mislav/dotfiles'
      }
      """
    And HTTPS is preferred
    When I successfully run `hub create`
    Then the url for "origin" should be "https://github.com/mislav/dotfiles.git"

  Scenario: Create in organization
    Given the GitHub API server:
      """
      post('/orgs/acme/repos') {
        json :full_name => 'acme/dotfiles'
      }
      """
    When I successfully run `hub create acme/dotfiles`
    Then the url for "origin" should be "git@github.com:acme/dotfiles.git"
    And the output should contain exactly "created repository: acme/dotfiles\n"

  Scenario: Creating repo failed
    Given the GitHub API server:
      """
      post('/user/repos') { status 500 }
      """
    When I run `hub create`
    Then the stderr should contain "Error creating repository: Internal Server Error (HTTP 500)"
    And the exit status should be 1
    And there should be no "origin" remote

  Scenario: With custom name
    Given the GitHub API server:
      """
      post('/user/repos') {
        assert :name => 'myconfig'
        json :full_name => 'mislav/myconfig'
      }
      """
    When I successfully run `hub create myconfig`
    Then the url for "origin" should be "git@github.com:mislav/myconfig.git"

  Scenario: With description and homepage
    Given the GitHub API server:
      """
      post('/user/repos') {
        assert :description => 'mydesc',
               :homepage => 'http://example.com'
        json :full_name => 'mislav/dotfiles'
      }
      """
    When I successfully run `hub create -d mydesc -h http://example.com`
    Then the url for "origin" should be "git@github.com:mislav/dotfiles.git"

  Scenario: Not in git repo
    Given the current dir is not a repo
    When I run `hub create`
    Then the stderr should contain "'create' must be run from inside a git repository"
    And the exit status should be 1

  Scenario: Origin remote already exists
    Given the GitHub API server:
      """
      post('/user/repos') {
        json :full_name => 'mislav/dotfiles'
      }
      """
    And the "origin" remote has url "git://github.com/mislav/dotfiles.git"
    When I successfully run `hub create`
    Then the url for "origin" should be "git://github.com/mislav/dotfiles.git"

  Scenario: GitHub repo already exists
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') { status 200 }
      """
    When I successfully run `hub create`
    Then the output should contain "mislav/dotfiles already exists on github.com\n"
    And the url for "origin" should be "git@github.com:mislav/dotfiles.git"

  Scenario: API response changes the clone URL
    Given the GitHub API server:
      """
      post('/user/repos') {
        json :full_name => 'Mooslav/myconfig'
      }
      """
    When I successfully run `hub create`
    Then the url for "origin" should be "git@github.com:Mooslav/myconfig.git"
    And the output should contain exactly "created repository: Mooslav/myconfig\n"

  Scenario: Current directory contains spaces
    Given I am in "my dot files" git repo
    Given the GitHub API server:
      """
      post('/user/repos') {
        assert :name => 'my-dot-files'
        json :full_name => 'mislav/my-dot-files'
      }
      """
    When I successfully run `hub create`
    Then the url for "origin" should be "git@github.com:mislav/my-dot-files.git"

  Scenario: Verbose API output
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') { status 404 }
      post('/user/repos') {
        response['location'] = 'http://disney.com'
        json :full_name => 'mislav/dotfiles'
      }
      """
    And $HUB_VERBOSE is "on"
    When I successfully run `hub create`
    Then the stderr should contain:
      """
      > GET https://api.github.com/repos/mislav/dotfiles
      > Authorization: token [REDACTED]
      < HTTP 404
      """
    And the stderr should contain:
      """
      > POST https://api.github.com/user/repos
      > Authorization: token [REDACTED]
      """
    And the stderr should contain:
      """
      < HTTP 200
      < Location: http://disney.com
      {"full_name":"mislav/dotfiles"}\n
      """
