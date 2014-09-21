Feature: hub am
  Background:
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Apply a local patch
    When I run `hub am some.patch`
    Then the git command should be unchanged

  Scenario: Apply commits from pull request
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles/pulls/387') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3.patch;charset=utf-8'
        generate_patch "Create a README"
      }
      """
    When I successfully run `hub am -q -3 https://github.com/mislav/dotfiles/pull/387`
    Then there should be no output
    Then the latest commit message should be "Create a README"

  Scenario: Apply commits when TMPDIR is empty
    Given $TMPDIR is ""
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles/pulls/387') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3.patch;charset=utf-8'
        generate_patch "Create a README"
      }
      """
    When I successfully run `hub am -q https://github.com/mislav/dotfiles/pull/387`
    Then the latest commit message should be "Create a README"

  Scenario: Enterprise repo
    Given I am in "git://git.my.org/mislav/dotfiles.git" git repo
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    Given the GitHub API server:
      """
      get('/api/v3/repos/mislav/dotfiles/pulls/387') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3.patch;charset=utf-8'
        generate_patch "Create a README"
      }
      """
    When I successfully run `hub am -q -3 https://git.my.org/mislav/dotfiles/pull/387`
    Then the latest commit message should be "Create a README"

  Scenario: Apply patch from commit
    Given the GitHub API server:
      """
      get('/repos/davidbalbert/dotfiles/commits/fdb9921') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3.patch;charset=utf-8'
        generate_patch "Create a README"
      }
      """
    When I successfully run `hub am -q https://github.com/davidbalbert/dotfiles/commit/fdb9921`
    Then the latest commit message should be "Create a README"

  Scenario: Apply patch from gist
    Given the GitHub API server:
      """
      get('/gists/8da7fb575debd88c54cf', :host_name => 'api.github.com') {
        json :files => {
          'file.diff' => {
            :raw_url => "https://gist.github.com/raw/8da7fb575debd88c54cf/SHA/file.diff"
          }
        }
      }
      get('/raw/8da7fb575debd88c54cf/SHA/file.diff', :host_name => 'gist.github.com') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'text/plain;charset=utf-8'
        generate_patch "Create a README"
      }
      """
    When I successfully run `hub am -q https://gist.github.com/8da7fb575debd88c54cf`
    Then the latest commit message should be "Create a README"
