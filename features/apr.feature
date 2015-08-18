Feature: hub am
  Background:
    Given I am in "git://github.com/mislav/dotfiles.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Apply commits from pull request
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles/pulls/387') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3.patch;charset=utf-8'
        generate_patch "Create a README"
      }
      """
    When I successfully run `hub apr -q 387`
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
    When I successfully run `hub apr -q 387`
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
    When I successfully run `hub apr -q 387`
    Then the latest commit message should be "Create a README"
