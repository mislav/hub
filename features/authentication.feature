Feature: OAuth authentication
  Background:
    Given I am in "dotfiles" git repo

  Scenario: Ask for username & password, create authorization
    Given the GitHub API server:
      """
      require 'socket'
      require 'etc'
      machine_id = "#{Etc.getlogin}@#{Socket.gethostname}"

      post('/authorizations') {
        assert_basic_auth 'mislav', 'kitty'
        assert :scopes => ['repo'],
               :note => "hub for #{machine_id}",
               :note_url => 'http://hub.github.com/'
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

  Scenario: Rename & retry creating authorization if there's a token name collision
    Given the GitHub API server:
      """
      require 'socket'
      require 'etc'
      machine_id = "#{Etc.getlogin}@#{Socket.gethostname}"

      post('/authorizations') {
        assert_basic_auth 'mislav', 'kitty'
        if params[:note] == "hub for #{machine_id} 3"
          json :token => 'OTOKEN'
        else
          status 422
          json :message => 'Validation Failed',
               :errors => [{
                 :resource => 'OauthAccess',
                 :code => 'already_exists',
                 :field => 'description'
               }]
        end
      }
      get('/user') {
        json :login => 'MiSlAv'
      }
      post('/user/repos') {
        json :full_name => 'mislav/dotfiles'
      }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    Then the output should contain "github.com username:"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"

  Scenario: Avoid getting caught up in infinite recursion while retrying token names
    Given the GitHub API server:
      """
      tries = 0
      post('/authorizations') {
        tries += 1
        halt 400, json(:message => "too many tries") if tries >= 10
        status 422
        json :message => 'Validation Failed',
             :errors => [{
               :resource => 'OauthAccess',
               :code => 'already_exists',
               :field => 'description'
             }]
      }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    Then the output should contain:
      """
      Error creating repository: Unprocessable Entity (HTTP 422)
      Duplicate value for "description"
      """
    And the exit status should be 1
    And the file "../home/.config/hub" should not exist

  Scenario: Credentials from GITHUB_USER & GITHUB_PASSWORD
    Given the GitHub API server:
      """
      post('/authorizations') {
        assert_basic_auth 'mislav', 'kitty'
        json :token => 'OTOKEN'
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
      post('/authorizations') {
        assert_basic_auth 'mislav', 'kitty'
      }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "WRONG"
    Then the stderr should contain exactly:
      """
      Error creating repository: Unauthorized (HTTP 401)
      Bad credentials

      """
    And the exit status should be 1
    And the file "../home/.config/hub" should not exist

  Scenario: Personal access token used instead of password
    Given the GitHub API server:
      """
      post('/authorizations') {
        status 403
        json :message => "This API can only be accessed with username and password Basic Auth"
      }
      """
    When I run `hub create` interactively
    When I type "mislav"
    And I type "PERSONALACCESSTOKEN"
    Then the stderr should contain exactly:
      """
      Error creating repository: Forbidden (HTTP 403)
      This API can only be accessed with username and password Basic Auth

      """
    And the exit status should be 1
    And the file "../home/.config/hub" should not exist

  Scenario: Two-factor authentication, create authorization
    Given the GitHub API server:
      """
      post('/authorizations') {
        assert_basic_auth 'mislav', 'kitty'
        if request.env['HTTP_X_GITHUB_OTP'] == '112233'
          json :token => 'OTOKEN'
        else
          response.headers['X-GitHub-OTP'] = 'required; app'
          status 401
          json :message => "Must specify two-factor authentication OTP code."
        end
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
    And the output should not contain "warning: invalid two-factor code"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"

  Scenario: Retry entering two-factor authentication code
    Given the GitHub API server:
      """
      previous_otp_code = nil
      post('/authorizations') {
        assert_basic_auth 'mislav', 'kitty'
        if request.env['HTTP_X_GITHUB_OTP'] == '112233'
          halt 400 unless '666' == previous_otp_code
          json :token => 'OTOKEN'
        else
          previous_otp_code = request.env['HTTP_X_GITHUB_OTP']
          response.headers['X-GitHub-OTP'] = 'required; app'
          status 401
          json :message => "Must specify two-factor authentication OTP code."
        end
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
    And I type "666"
    And I type "112233"
    Then the output should contain "warning: invalid two-factor code"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"

  Scenario: Special characters in username & password
    Given the GitHub API server:
      """
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

  Scenario: Enterprise fork authentication with username & password, re-using existing authorization
    Given the GitHub API server:
      """
      require 'rack/auth/basic'
      post('/api/v3/authorizations', :host_name => 'git.my.org') {
        auth = Rack::Auth::Basic::Request.new(env)
        halt 401 unless auth.credentials == %w[mislav kitty]
        json :token => 'OTOKEN', :note_url => 'http://hub.github.com/'
      }
      get('/api/v3/user', :host_name => 'git.my.org') {
        json :login => 'mislav'
      }
      post('/api/v3/repos/evilchelu/dotfiles/forks', :host_name => 'git.my.org') { '' }
      """
    And "git.my.org" is a whitelisted Enterprise host
    And the "origin" remote has url "git@git.my.org:evilchelu/dotfiles.git"
    When I run `hub fork` interactively
    And I type "mislav"
    And I type "kitty"
    Then the output should contain "git.my.org password for mislav (never stored):"
    And the exit status should be 0
    And the file "../home/.config/hub" should contain "git.my.org"
    And the file "../home/.config/hub" should contain "user: mislav"
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"
    And the url for "mislav" should be "git@git.my.org:mislav/dotfiles.git"
