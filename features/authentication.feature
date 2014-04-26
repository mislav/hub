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
        assert :scopes => ['repo']
        json :token => 'OTOKEN'
      }
      get('/user') {
        halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN'
        json :login => 'MiSlAv'
      }
      post('/user/repos') {
        halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN'
        json :full_name => 'mislav/dotfiles'
      }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    Then the output should contain "github.com username:"
    And the output should contain "github.com password for mislav (never stored):"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "user: MiSlAv"
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
          {:token => 'SKIPPD', :note_url => 'http://example.com'},
          {:token => 'OTOKEN', :note_url => 'http://hub.github.com/'}
        ]
      }
      get('/user') {
        json :login => 'mislav'
      }
      post('/user/repos') {
        json :full_name => 'mislav/dotfiles'
      }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    Then the output should contain "github.com password for mislav (never stored):"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"

  Scenario: Re-use existing authorization with an old URL
    Given the GitHub API server:
      """
      require 'rack/auth/basic'
      get('/authorizations') {
        auth = Rack::Auth::Basic::Request.new(env)
        halt 401 unless auth.credentials == %w[mislav kitty]
        json [
          {:token => 'OTOKEN', :note => 'hub', :note_url => 'http://defunkt.io/hub/'}
        ]
      }
      post('/authorizations') {
        status 422
        json :message => "Validation Failed",
          :errors => [{:resource => "OauthAccess", :code => "already_exists", :field => "description"}]
      }
      get('/user') {
        json :login => 'mislav'
      }
      post('/user/repos') {
        json :full_name => 'mislav/dotfiles'
      }
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
          {:token => 'OTOKEN', :note_url => 'http://hub.github.com/'}
        ]
      }
      get('/user') {
        json :login => 'mislav'
      }
      post('/user/repos') {
        json :full_name => 'mislav/dotfiles'
      }
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

  Scenario: Two-factor authentication, create authorization
    Given the GitHub API server:
      """
      require 'rack/auth/basic'
      get('/authorizations') {
        auth = Rack::Auth::Basic::Request.new(env)
        halt 401 unless auth.credentials == %w[mislav kitty]
        if request.env['HTTP_X_GITHUB_OTP'] != "112233"
          response.headers['X-GitHub-OTP'] = "required;application"
          halt 401
        end
        json [ ]
      }
      post('/authorizations') {
        auth = Rack::Auth::Basic::Request.new(env)
        halt 401 unless auth.credentials == %w[mislav kitty]
        halt 412 unless params[:scopes]
        if request.env['HTTP_X_GITHUB_OTP'] != "112233"
          response.headers['X-GitHub-OTP'] = "required;application"
          halt 401
        end
        json :token => 'OTOKEN'
      }
      get('/user') {
        json :login => 'mislav'
      }
      post('/user/repos') {
        json :full_name => 'mislav/dotfiles'
      }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    And I type "112233"
    Then the output should contain "github.com password for mislav (never stored):"
    Then the output should contain "two-factor authentication code:"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"

  Scenario: Two-factor authentication, re-use existing authorization
    Given the GitHub API server:
      """
      token = 'OTOKEN'
      post('/authorizations') {
        assert_basic_auth 'mislav', 'kitty'
        token << 'SMS'
        status 412
      }
      get('/authorizations') {
        assert_basic_auth 'mislav', 'kitty'
        if request.env['HTTP_X_GITHUB_OTP'] != "112233"
          response.headers['X-GitHub-OTP'] = "required;application"
          halt 401
        end
        json [ {
          :token => token,
          :note_url => 'http://hub.github.com/'
          } ]
      }
      get('/user') {
        json :login => 'mislav'
      }
      post('/user/repos') {
        json :full_name => 'mislav/dotfiles'
      }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    And I type "112233"
    Then the output should contain "github.com password for mislav (never stored):"
    Then the output should contain "two-factor authentication code:"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "oauth_token: OTOKENSMS"

  Scenario: Special characters in username & password
    Given the GitHub API server:
      """
      get('/authorizations') { '[]' }
      post('/authorizations') {
        assert_basic_auth 'mislav@example.com', 'my pass@phrase ok?'
        json :token => 'OTOKEN'
      }
      get('/user') {
        json :login => 'mislav'
      }
      get('/repos/mislav/dotfiles') { status 200 }
      """
    When I run `hub create` interactively
    When I type "mislav@example.com"
    And I type "my pass@phrase ok?"
    Then the output should contain "github.com password for mislav@example.com (never stored):"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "user: mislav"
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"
