Feature: OAuth authentication
  Background:
    Given I am in "dotfiles" git repo

  Scenario: Ask for username & password, create authorization
    Given the GitHub API server:
      """
      require 'rack/auth/basic'
      get('/authorizations') { '[]' }
      post('/authorizations') {
        auth = Rack::Auth::Basic::Request.new(env)
        halt 401 unless auth.credentials == %w[mislav kitty]
        halt 400 unless params[:scopes] == ['repo']
        json :token => 'OTOKEN'
      }
      post('/user/repos') { status 200 }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    Then the output should contain "github.com username:"
    And the output should contain "github.com password for mislav (never stored):"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"
    And the file "../home/.config/hub" should have mode "0600"

  Scenario: Ask for username & password, re-use existing authorization
    Given the GitHub API server:
      """
      require 'rack/auth/basic'
      get('/authorizations') {
        auth = Rack::Auth::Basic::Request.new(env)
        halt 401 unless auth.credentials == %w[mislav kitty]
        json [
          {:token => 'SKIPPD', :app => {:url => 'http://example.com'}},
          {:token => 'OTOKEN', :app => {:url => 'http://defunkt.io/hub/'}}
        ]
      }
      post('/user/repos') { status 200 }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    Then the output should contain "github.com password for mislav (never stored):"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"

  Scenario: Credentials from GITHUB_USER & GITHUB_PASSWORD
    Given the GitHub API server:
      """
      require 'rack/auth/basic'
      get('/authorizations') {
        auth = Rack::Auth::Basic::Request.new(env)
        halt 401 unless auth.credentials == %w[mislav kitty]
        json [
          {:token => 'OTOKEN', :app => {:url => 'http://defunkt.io/hub/'}}
        ]
      }
      post('/user/repos') { status 200 }
      """
    Given $GITHUB_USER is "mislav"
    And $GITHUB_PASSWORD is "kitty"
    When I successfully run `hub create`
    Then the output should not contain "github.com password for mislav"
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"

  Scenario: Wrong password
    Given the GitHub API server:
      """
      require 'rack/auth/basic'
      get('/authorizations') {
        auth = Rack::Auth::Basic::Request.new(env)
        halt 401 unless auth.credentials == %w[mislav kitty]
      }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "WRONG"
    Then the stderr should contain "Error creating repository: Unauthorized (HTTP 401)"
    And the exit status should be 1
    And the file "../home/.config/hub" should not exist
