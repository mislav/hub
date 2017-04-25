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
        halt 415 unless request.accept?('application/vnd.github.v3+json')
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
        }, :maintainer_can_modify => false
      }
      """
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git fetch origin refs/pull/77/head:mislav-fixes" should be run
    And "git checkout -f mislav-fixes -q" should be run
    And "mislav-fixes" should merge "refs/pull/77/head" from remote "origin"

  Scenario: No matching remotes for pull request base
    Given the GitHub API server:
      """
      get('/repos/mislav/jekyll/pulls/77') {
        json :number => 77, :base => {
          :repo => {
            :name => 'jekyll',
            :html_url => 'https://github.com/mislav/jekyll',
            :owner => { :login => "mislav" },
          }
        }
      }
      """
    When I run `hub checkout -f https://github.com/mislav/jekyll/pull/77 -q`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      could not find git remote for mislav/jekyll\n
      """

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
        }, :maintainer_can_modify => false
      }
      """
    When I run `hub checkout https://github.com/mojombo/jekyll/pull/77 fixes-from-mislav`
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
        }
      }
      """
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git fetch origin +refs/heads/fixes:refs/remotes/origin/fixes" should be run
    And "git checkout -f -b fixes --track origin/fixes -q" should be run

  Scenario: Same-repo with custom branch name
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
        }
      }
      """
    When I run `hub checkout https://github.com/mojombo/jekyll/pull/77 mycustombranch`
    Then "git fetch origin +refs/heads/fixes:refs/remotes/origin/fixes" should be run
    And "git checkout -b mycustombranch --track origin/fixes" should be run

  Scenario: Unavailable fork
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :head => {
          :ref => "fixes",
          :repo => nil
        }, :base => {
          :repo => {
            :name => "jekyll",
            :html_url => "https://github.com/mojombo/jekyll",
            :owner => { :login => "mojombo" },
          }
        }
      }
      """
    When I run `hub checkout https://github.com/mojombo/jekyll/pull/77`
    Then "git fetch origin refs/pull/77/head:pr-77" should be run
    And "git checkout pr-77" should be run
    And "pr-77" should merge "refs/pull/77/head" from remote "origin"

  Scenario: Reuse existing remote for head branch
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
        }
      }
      """
    And the "mislav" remote has url "git://github.com/mislav/jekyll.git"
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git fetch mislav +refs/heads/fixes:refs/remotes/mislav/fixes" should be run
    And "git checkout -f -b fixes --track mislav/fixes -q" should be run

  Scenario: Reuse existing remote and branch
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
        }
      }
      """
    And the "mislav" remote has url "git://github.com/mislav/jekyll.git"
    And I am on the "fixes" branch
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git fetch mislav +refs/heads/fixes:refs/remotes/mislav/fixes" should be run
    And "git checkout -f fixes -q" should be run
    And "git merge --ff-only refs/remotes/mislav/fixes" should be run

  Scenario: Modifiable fork
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :head => {
          :ref => "fixes",
          :repo => {
            :owner => { :login => "mislav" },
            :name => "jekyll",
            :html_url => "https://github.com/mislav/jekyll.git",
            :private => false
          },
        }, :base => {
          :repo => {
            :name => 'jekyll',
            :html_url => 'https://github.com/mojombo/jekyll',
            :owner => { :login => "mojombo" },
          }
        }, :maintainer_can_modify => true
      }
      """
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git fetch origin refs/pull/77/head:mislav-fixes" should be run
    And "git checkout -f mislav-fixes -q" should be run
    And "mislav-fixes" should merge "refs/heads/fixes" from remote "git@github.com:mislav/jekyll.git"

  Scenario: Modifiable fork with HTTPS
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :head => {
          :ref => "fixes",
          :repo => {
            :owner => { :login => "mislav" },
            :name => "jekyll",
            :html_url => "https://github.com/mislav/jekyll.git",
            :private => false
          },
        }, :base => {
          :repo => {
            :name => 'jekyll',
            :html_url => 'https://github.com/mojombo/jekyll',
            :owner => { :login => "mojombo" },
          }
        }, :maintainer_can_modify => true
      }
      """
    And HTTPS is preferred
    When I run `hub checkout -f https://github.com/mojombo/jekyll/pull/77 -q`
    Then "git fetch origin refs/pull/77/head:mislav-fixes" should be run
    And "git checkout -f mislav-fixes -q" should be run
    And "mislav-fixes" should merge "refs/heads/fixes" from remote "https://github.com/mislav/jekyll.git"
