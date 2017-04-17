Feature: hub pr
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
    When I successfully run `hub pr --am -q -3 387`
    Then there should be no output
    Then the latest commit message should be "Create a README"

  Scenario: Browse
    When I successfully run `hub pr --browse 1`
    Then "open https://github.com/mislav/dotfiles/pull/1" should be run
