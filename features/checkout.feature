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
        halt 406 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3+json;charset=utf-8'
        json :head => {
          :ref => "fixes",
          :repo => {
            :owner => { :login => "mislav" },
            :name => "jekyll",
            :private => false
          }
        }
      }
      """
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git fetch git://github.com/mojombo/jekyll.git pull/77/head:mislav-fixes" should be run
    And "git checkout -f mislav-fixes -q" should be run

  Scenario: Pull request from a renamed fork
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :head => {
          :ref => "fixes",
          :repo => {
            :owner => { :login => "mislav" },
            :name => "jekyll-blog",
            :private => false
          }
        }
      }
      """
    When I run `hub checkout https://github.com/mojombo/jekyll/pull/77`
    Then "git fetch git://github.com/mojombo/jekyll.git pull/77/head:mislav-fixes" should be run
    And "git checkout mislav-fixes" should be run

  Scenario: Custom name for new branch
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :head => {
          :ref => "fixes",
          :repo => {
            :owner => { :login => "mislav" },
            :name => "jekyll",
            :private => false
          }
        }
      }
      """
    When I run `hub checkout https://github.com/mojombo/jekyll/pull/77 fixes-from-mislav`
    Then "git fetch git://github.com/mojombo/jekyll.git pull/77/head:fixes-from-mislav" should be run
    And "git checkout fixes-from-mislav" should be run

  Scenario: Private pull request
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :head => {
          :ref => "fixes",
          :repo => {
            :owner => { :login => "mislav" },
            :name => "jekyll",
            :private => true
          }
        }
      }
      """
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git fetch git://github.com/mojombo/jekyll.git pull/77/head:mislav-fixes" should be run
    And "git checkout -f mislav-fixes -q" should be run

  Scenario: Remote for user already exists
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :head => {
          :ref => "fixes",
          :repo => {
            :owner => { :login => "mislav" },
            :name => "jekyll",
            :private => false
          }
        }
      }
      """
    And the "mislav" remote has url "git://github.com/mislav/jekyll.git"
    When I run `hub checkout https://github.com/mojombo/jekyll/pull/77`
    Then "git fetch mislav pull/77/head:mislav-fixes" should be run
    And "git checkout mislav-fixes" should be run
