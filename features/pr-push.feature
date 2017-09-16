Feature: hub pr push <PULLREQ-NUMBER>
  Background:
    Given I am in "git://github.com/mojombo/jekyll.git" git repo
    And I am "mojombo" on github.com with OAuth token "OTOKEN"

  Scenario: Push to a pull request
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :base => {
          :ref => "fixes",
          :repo => {
            :name => "jekyll",
            :owner => { :login => "mojombo" },
          }
        }, :head => {
				  :ref => "patch-1",
          :repo => {
            :name => "jekyll",
            :html_url => "https://github.com/jridgewell/jekyll",
            :owner => { :login => "jridgewell" },
          }
        },
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
		Given HTTPS is preferred
    When I run `hub pr push 77`
		Then "git push https://github.com/jridgewell/jekyll.git master:patch-1" should be run

  Scenario: Push to a pull request (SSH fallback)
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :base => {
          :ref => "fixes",
          :repo => {
            :name => "jekyll",
            :owner => { :login => "mojombo" },
          }
        }, :head => {
				  :ref => "patch-1",
          :repo => {
            :name => "jekyll",
            :html_url => "https://github.com/jridgewell/jekyll",
            :owner => { :login => "jridgewell" },
          }
        },
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
    When I run `hub pr push 77`
    Then "git push git@github.com:jridgewell/jekyll.git master:patch-1" should be run

  Scenario: Push to a pull request with branch
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :base => {
          :ref => "fixes",
          :repo => {
            :name => "jekyll",
            :owner => { :login => "mojombo" },
          }
        }, :head => {
				  :ref => "patch-1",
          :repo => {
            :name => "jekyll",
            :html_url => "https://github.com/jridgewell/jekyll",
            :owner => { :login => "jridgewell" },
          }
        },
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
    When I run `hub pr push 77 local`
    Then "git push git@github.com:jridgewell/jekyll.git local:patch-1" should be run

  Scenario: Push to a pull request with force
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :base => {
          :ref => "fixes",
          :repo => {
            :name => "jekyll",
            :owner => { :login => "mojombo" },
          }
        }, :head => {
				  :ref => "patch-1",
          :repo => {
            :name => "jekyll",
            :html_url => "https://github.com/jridgewell/jekyll",
            :owner => { :login => "jridgewell" },
          }
        },
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
    When I run `hub pr push 77 -f`
    Then "git push git@github.com:jridgewell/jekyll.git master:patch-1 --force" should be run


  Scenario: Push to a pull request with force and gusto
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :base => {
          :ref => "fixes",
          :repo => {
            :name => "jekyll",
            :owner => { :login => "mojombo" },
          }
        }, :head => {
				  :ref => "patch-1",
          :repo => {
            :name => "jekyll",
            :html_url => "https://github.com/jridgewell/jekyll",
            :owner => { :login => "jridgewell" },
          }
        },
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
    When I run `hub pr push 77 --force`
    Then "git push git@github.com:jridgewell/jekyll.git master:patch-1 --force" should be run


  Scenario: Push to a pull request while setting upstream
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :base => {
          :ref => "fixes",
          :repo => {
            :name => "jekyll",
            :owner => { :login => "mojombo" },
          }
        }, :head => {
				  :ref => "patch-1",
          :repo => {
            :name => "jekyll",
            :html_url => "https://github.com/jridgewell/jekyll",
            :owner => { :login => "jridgewell" },
          }
        },
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
    When I run `hub pr push 77 -u`
    Then "git push git@github.com:jridgewell/jekyll.git master:patch-1 --set-upstream" should be run


  Scenario: Push to a pull request with force and gusto
    Given the GitHub API server:
      """
      get('/repos/mojombo/jekyll/pulls/77') {
        json :number => 77, :base => {
          :ref => "fixes",
          :repo => {
            :name => "jekyll",
            :owner => { :login => "mojombo" },
          }
        }, :head => {
				  :ref => "patch-1",
          :repo => {
            :name => "jekyll",
            :html_url => "https://github.com/jridgewell/jekyll",
            :owner => { :login => "jridgewell" },
          }
        },
        :html_url => 'https://github.com/mojombo/jekyll/pull/77'
      }
      """
    When I run `hub pr push 77 --set-upstream`
    Then "git push git@github.com:jridgewell/jekyll.git master:patch-1 --set-upstream" should be run

