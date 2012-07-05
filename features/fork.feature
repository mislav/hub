Feature: hub fork
  Background:
    Given I am in "dotfiles" git repo
    And the "origin" remote has url "git://github.com/evilchelu/dotfiles.git"
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Fork the repository
    Given the GitHub API server:
      """
      before { halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN' }
      get('/repos/evilchelu/dotfiles', :host_name => 'api.github.com') { '' }
      post('/repos/evilchelu/dotfiles/forks', :host_name => 'api.github.com') { '' }
      """
    When I successfully run `hub fork`
    Then the output should contain exactly "new remote: mislav\n"
    And "git remote add -f mislav git@github.com:mislav/dotfiles.git" should be run
    And the url for "mislav" should be "git@github.com:mislav/dotfiles.git"

  Scenario: --no-remote
    Given the GitHub API server:
      """
      get('/repos/evilchelu/dotfiles') { '' }
      post('/repos/evilchelu/dotfiles/forks') { '' }
      """
    When I successfully run `hub fork --no-remote`
    Then there should be no output
    And there should be no "mislav" remote

  Scenario: Fork failed
    Given the GitHub API server:
      """
      get('/repos/evilchelu/dotfiles') { '' }
      post('/repos/evilchelu/dotfiles/forks') { halt 500 }
      """
    When I run `hub fork`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error creating fork: Internal Server Error (HTTP 500)\n
      """
    And there should be no "mislav" remote

  Scenario: Fork already exists
    Given the GitHub API server:
      """
      get('/repos/evilchelu/dotfiles') { '' }
      get('/repos/mislav/dotfiles') { '' }
      """
    When I run `hub fork`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error creating fork: mislav/dotfiles already exists on github.com\n
      """
    And there should be no "mislav" remote

  Scenario: Invalid OAuth token
    Given the GitHub API server:
      """
      before { halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN' }
      get('/repos/evilchelu/dotfiles') { '' }
      """
    And I am "mislav" on github.com with OAuth token "WRONGTOKEN"
    When I run `hub fork`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error creating fork: Unauthorized (HTTP 401)\n
      """

  Scenario: HTTPS is preferred
    Given the GitHub API server:
      """
      get('/repos/evilchelu/dotfiles') { '' }
      post('/repos/evilchelu/dotfiles/forks') { '' }
      """
    And HTTPS is preferred
    When I successfully run `hub fork`
    Then the output should contain exactly "new remote: mislav\n"
    And the url for "mislav" should be "https://github.com/mislav/dotfiles.git"

  Scenario: Not in repo
    Given the current dir is not a repo
    When I run `hub fork`
    Then the exit status should be 1
    And the stderr should contain "fatal: Not a git repository"

  Scenario: Unknown host
    Given the "origin" remote has url "git@git.my.org:evilchelu/dotfiles.git"
    When I run `hub fork`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error: repository under 'origin' remote is not a GitHub project\n
      """

  Scenario: Enterprise fork
    Given the GitHub API server:
      """
      before { halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token FITOKEN' }
      get('/api/v3/repos/evilchelu/dotfiles', :host_name => 'git.my.org') { '' }
      post('/api/v3/repos/evilchelu/dotfiles/forks', :host_name => 'git.my.org') { '' }
      """
    And the "origin" remote has url "git@git.my.org:evilchelu/dotfiles.git"
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    When I successfully run `hub fork`
    Then the url for "mislav" should be "git@git.my.org:mislav/dotfiles.git"
