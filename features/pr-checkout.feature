Feature: hub pr checkout <PULLREQ-NUMBER>
  Background:
    Given I am in "git://github.com/mojombo/jekyll.git" git repo
    And I am "mojombo" on github.com with OAuth token "OTOKEN"

  Scenario: Checkout a pull request
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :head => {
          :ref => "fixes",
          :repo => {
            :owner => { :login => "mislav" },
            :name => "jekyll",
            :private => false
          }
        }, :base => {
          :repo => {
            :name => 'jekyll',
            :html_url => 'https://github.com/mojombo/jekyll',
            :owner => { :login => "mojombo" },
          }
        },
        :maintainer_can_modify => false,
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
    When I run `hub pr checkout 77`
    Then "git fetch origin refs/pull/77/head:mislav-fixes" should be run
    And "git checkout mislav-fixes" should be run
    And "mislav-fixes" should merge "refs/pull/77/head" from remote "origin"

  Scenario: Custom name for new branch
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :head => {
          :ref => "fixes",
          :repo => {
            :name => "jekyll",
            :owner => { :login => "mislav" },
          }
        }, :base => {
          :repo => {
            :name => 'jekyll',
            :html_url => 'https://github.com/mojombo/jekyll',
            :owner => { :login => "mojombo" },
          }
        },
        :maintainer_can_modify => false,
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
    When I run `hub pr checkout 77 fixes-from-mislav`
    Then "git fetch origin refs/pull/77/head:fixes-from-mislav" should be run
    And "git checkout fixes-from-mislav" should be run
    And "fixes-from-mislav" should merge "refs/pull/77/head" from remote "origin"

  Scenario: Same-repo
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :head => {
          :ref => "fixes",
          :repo => {
            :name => "jekyll",
            :owner => { :login => "mojombo" },
          }
        }, :base => {
          :repo => {
            :name => "jekyll",
            :html_url => "https://github.com/mojombo/jekyll",
            :owner => { :login => "mojombo" },
          }
        },
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
    When I run `hub pr checkout 77`
    Then "git fetch origin +refs/heads/fixes:refs/remotes/origin/fixes" should be run
    And "git checkout -b fixes --track origin/fixes" should be run
