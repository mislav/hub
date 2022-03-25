Feature: hub remote add
  Background:
    Given I am "EvilChelu" on GitHub.com
    And I am in "dotfiles" git repo

  Scenario: Add origin remote for my own repo
    Given there are no remotes
    When I successfully run `hub remote add origin`
    Then the url for "origin" should be "https://github.com/EvilChelu/dotfiles.git"
    And the output should not contain anything

  Scenario: Add origin remote for my own repo using -C
    Given there are no remotes
    And I cd to ".."
    When I successfully run `hub -C dotfiles remote add origin`
    And I cd to "dotfiles"
    Then the url for "origin" should be "https://github.com/EvilChelu/dotfiles.git"
    And the output should not contain anything

  Scenario: Unchanged public remote add
    When I successfully run `hub remote add origin http://github.com/defunkt/resque.git`
    Then the url for "origin" should be "http://github.com/defunkt/resque.git"
    And the output should not contain anything

  Scenario: Unchanged private remote add
    When I successfully run `hub remote add origin git@github.com:defunkt/resque.git`
    Then the url for "origin" should be "git@github.com:defunkt/resque.git"
    And the output should not contain anything

  Scenario: Unchanged local path remote add
    When I successfully run `hub remote add myremote ./path`
    Then the git command should be unchanged
    And the output should not contain anything

  Scenario: Unchanged local absolute path remote add
    When I successfully run `hub remote add myremote /path`
    Then the git command should be unchanged
    And the output should not contain anything

  Scenario: Unchanged remote add with host alias
    When I successfully run `hub remote add myremote server:/git/repo.git`
    Then the git command should be unchanged
    And the output should not contain anything

  Scenario: Add new remote for Enterprise repo
    Given "git.my.org" is a whitelisted Enterprise host
    And git protocol is preferred
    And I am "ProLoser" on git.my.org with OAuth token "FITOKEN"
    And the "origin" remote has url "git@git.my.org:mislav/topsekrit.git"
    When I successfully run `hub remote add another`
    Then the url for "another" should be "git@git.my.org:another/topsekrit.git"
    And the output should not contain anything

  Scenario: Add public remote
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub remote add mislav`
    Then the url for "mislav" should be "https://github.com/mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Add detected private remote
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => true,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => false }
      }
      """
    And git protocol is preferred
    When I successfully run `hub remote add mislav`
    Then the url for "mislav" should be "git@github.com:mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Add remote with push access
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => true }
      }
      """
    When I successfully run `hub remote add mislav`
    Then the url for "mislav" should be "https://github.com/mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Add remote for missing repo
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        status 404
      }
      """
    When I run `hub remote add mislav`
    Then the exit status should be 1
    And the output should contain exactly:
      """
      Error: repository mislav/dotfiles doesn't exist\n
      """

  Scenario: Add explicitly private remote
    Given git protocol is preferred
    When I successfully run `hub remote add -p mislav`
    Then the url for "mislav" should be "git@github.com:mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Remote for my own repo is automatically private
    Given git protocol is preferred
    When I successfully run `hub remote add evilchelu`
    Then the url for "evilchelu" should be "git@github.com:EvilChelu/dotfiles.git"
    And the output should not contain anything

  Scenario: Add remote with arguments
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub remote add -f mislav`
    Then "git remote add -f mislav https://github.com/mislav/dotfiles.git" should be run
    And the output should not contain anything

  Scenario: Add remote with branch argument
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub remote add -f -t feature mislav`
    Then "git remote add -f -t feature mislav https://github.com/mislav/dotfiles.git" should be run
    And the output should not contain anything

  Scenario: Add named public remote
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub remote add mm mislav`
    Then the url for "mm" should be "https://github.com/mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: set-url
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => false }
      }
      """
    Given the "origin" remote has url "https://github.com/evilchelu/dotfiles.git"
    When I successfully run `hub remote set-url origin mislav`
    Then the url for "origin" should be "https://github.com/mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Add public remote including repo name
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfilez.js') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub remote add mislav/dotfilez.js`
    Then the url for "mislav" should be "https://github.com/mislav/dotfilez.js.git"
    And the output should not contain anything

  Scenario: Add named public remote including repo name
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfilez.js') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub remote add mm mislav/dotfilez.js`
    Then the url for "mm" should be "https://github.com/mislav/dotfilez.js.git"
    And the output should not contain anything

  Scenario: Add named private remote
    Given git protocol is preferred
    When I successfully run `hub remote add -p mm mislav`
    Then the url for "mm" should be "git@github.com:mislav/dotfiles.git"
    And the output should not contain anything

  Scenario: Add private remote including repo name
    When I successfully run `hub remote add -p mislav/dotfilez.js`
    Then the url for "mislav" should be "https://github.com/mislav/dotfilez.js.git"
    And the output should not contain anything

  Scenario: Add named private remote including repo name
    When I successfully run `hub remote add -p mm mislav/dotfilez.js`
    Then the url for "mm" should be "https://github.com/mislav/dotfilez.js.git"
    And the output should not contain anything

  Scenario: Add named private remote for my own repo including repo name
    When I successfully run `hub remote add ec evilchelu/dotfilez.js`
    Then the url for "ec" should be "https://github.com/EvilChelu/dotfilez.js.git"
    And the output should not contain anything

  Scenario: Avoid crash in argument parsing
    When I successfully run `hub --noop remote add a b evilchelu`
    Then the output should contain exactly "git remote add a b evilchelu\n"
