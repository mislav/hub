Feature: hub apply
  Background:
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"
    And I make a commit

  Scenario: Apply a local patch
    When I run `hub apply some.patch`
    Then the git command should be unchanged
    And the file "README.md" should not exist

  Scenario: Apply commits from pull request
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles/pulls/387') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3.patch;charset=utf-8'
        generate_patch "Create a README"
      }
      """
    When I successfully run `hub apply -3 https://github.com/mislav/dotfiles/pull/387`
    Then there should be no output
    Then a file named "README.md" should exist

  Scenario: Apply commits when TMPDIR is empty
    Given $TMPDIR is ""
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles/pulls/387') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3.patch;charset=utf-8'
        generate_patch "Create a README"
      }
      """
    When I successfully run `hub apply https://github.com/mislav/dotfiles/pull/387`
    Then a file named "README.md" should exist

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
    When I successfully run `hub apply https://git.my.org/mislav/dotfiles/pull/387`
    Then a file named "README.md" should exist

  Scenario: Apply patch from commit
    Given the GitHub API server:
      """
      get('/repos/davidbalbert/dotfiles/commits/fdb9921') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3.patch;charset=utf-8'
        generate_patch "Create a README"
      }
      """
    When I successfully run `hub apply https://github.com/davidbalbert/dotfiles/commit/fdb9921`
    Then a file named "README.md" should exist

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
    When I successfully run `hub apply https://gist.github.com/8da7fb575debd88c54cf`
    Then a file named "README.md" should exist
