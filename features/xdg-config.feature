Feature: XDG BaseDir Spec

  According to the XDG Base Directory Specification 
  (http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html),
  if the XDG_CONFIG_HOME environment variable is set, this should be used
  to store application specific configuration. If unset, then a default of
  ~/.config should be used.

  Background:
    Given I am in "dotfiles" git repo

  Scenario: XDG_CONFIG_HOME empty, store authorization in default directory
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
    Given $XDG_CONFIG_HOME is ""
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    And the file "../home/.config/hub" should contain "user: MiSlAv"
    And the file "../home/.config/hub" should contain "oauth_token: OTOKEN"
    And the file "../home/.config/hub" should have mode "0600"

Scenario: XDG_CONFIG_HOME is not empty, store authorization in specified directory
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
    Given $XDG_CONFIG_HOME is "../home/.xdgconfig"
    Given a directory named "../home/.xdconfig"
    When I run `hub create` interactively
    When I type "mislav"
    And I type "kitty"
    And the file "../home/.xdgconfig/hub" should contain "user: MiSlAv"
    And the file "../home/.xdgconfig/hub" should contain "oauth_token: OTOKEN"
    And the file "../home/.xdgconfig/hub" should have mode "0600"
