Feature: hub checkout <PULLREQ-URL>
  Background:
    Given I am in "git://github.com/mojombo/jekyll.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Unchanged command
    When I run `hub checkout master`
    Then "git checkout master" should be run

  Scenario: Checkout a pull request
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        halt 406 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3+json'
        json :head => {
          :label => 'mislav:fixes',
          :repo => { :private => false }
        }
      }
      """
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git remote add -f -t fixes mislav git://github.com/mislav/jekyll.git" should be run
    And "git checkout -f --track -B mislav-fixes mislav/fixes -q" should be run

  Scenario: Custom name for new branch
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :head => {
          :label => 'mislav:fixes',
          :repo => { :private => false }
        }
      }
      """
    When I run `hub checkout https://github.com/mojombo/jekyll/pull/77 fixes-from-mislav`
    Then "git remote add -f -t fixes mislav git://github.com/mislav/jekyll.git" should be run
    And "git checkout --track -B fixes-from-mislav mislav/fixes" should be run

  Scenario: Private pull request
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :head => {
          :label => 'mislav:fixes',
          :repo => { :private => true }
        }
      }
      """
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git remote add -f -t fixes mislav git@github.com:mislav/jekyll.git" should be run
    And "git checkout -f --track -B mislav-fixes mislav/fixes -q" should be run

  Scenario: Custom name for new branch
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :head => {
          :label => 'mislav:fixes',
          :repo => { :private => false }
        }
      }
      """
    When I run `hub checkout https://github.com/mojombo/jekyll/pull/77 fixes-from-mislav`
    Then "git remote add -f -t fixes mislav git://github.com/mislav/jekyll.git" should be run
    And "git checkout --track -B fixes-from-mislav mislav/fixes" should be run

  Scenario: Remote for user already exists
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :head => {
          :label => 'mislav:fixes',
          :repo => { :private => false }
        }
      }
      """
    And the "mislav" remote has url "git://github.com/mislav/jekyll.git"
    When I run `hub checkout https://github.com/mojombo/jekyll/pull/77`
    Then "git remote set-branches --add mislav fixes" should be run
    And "git fetch mislav +refs/heads/fixes:refs/remotes/mislav/fixes" should be run
    And "git checkout --track -B mislav-fixes mislav/fixes" should be run
