Feature: hub create
  Background:
    Given I am in "dotfiles" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Create repo
    Given the GitHub API server:
      """
      post('/user/repos') {
        halt 400 if params[:private]
        status 200
      }
      """
    When I successfully run `hub create`
    Then the url for "origin" should be "git@github.com:mislav/dotfiles.git"
    And the output should contain exactly "created repository: mislav/dotfiles\n"

  Scenario: Create private repo
    Given the GitHub API server:
      """
      post('/user/repos') {
        halt 400 unless params[:private]
        status 200
      }
      """
    When I successfully run `hub create -p`
    Then the url for "origin" should be "git@github.com:mislav/dotfiles.git"

  Scenario: HTTPS is preferred
    Given the GitHub API server:
      """
      post('/user/repos') { status 200 }
      """
    And HTTPS is preferred
    When I successfully run `hub create`
    Then the url for "origin" should be "https://github.com/mislav/dotfiles.git"

  Scenario: Create in organization
    Given the GitHub API server:
      """
      post('/orgs/acme/repos') { status 200 }
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
        halt 400 unless params[:name] == 'myconfig'
        status 200
      }
      """
    When I successfully run `hub create myconfig`
    Then the url for "origin" should be "git@github.com:mislav/myconfig.git"

  Scenario: With description and homepage
    Given the GitHub API server:
      """
      post('/user/repos') {
        halt 400 unless params[:description] == 'mydesc' and
          params[:homepage] == 'http://example.com'
        status 200
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
      post('/user/repos') { status 200 }
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
