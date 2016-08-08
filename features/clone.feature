Feature: hub clone
  Background:
    Given I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Clone a public repo
    Given the GitHub API server:
      """
      get('/repos/rtomayko/ronn') {
        json :private => false,
             :name => 'ronn', :owner => { :login => 'rtomayko' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub clone rtomayko/ronn`
    Then it should clone "git://github.com/rtomayko/ronn.git"
    And there should be no output

  Scenario: Clone a public repo with period in name
    Given the GitHub API server:
      """
      get('/repos/hookio/hook.js') {
        json :private => false,
             :name => 'hook.js', :owner => { :login => 'hookio' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub clone hookio/hook.js`
    Then it should clone "git://github.com/hookio/hook.js.git"
    And there should be no output

  Scenario: Clone a public repo that starts with a period
    Given the GitHub API server:
      """
      get('/repos/zhuangya/.vim') {
        json :private => false,
             :name => '.vim', :owner => { :login => 'zhuangya' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub clone zhuangya/.vim`
    Then it should clone "git://github.com/zhuangya/.vim.git"
    And there should be no output

  Scenario: Clone a repo even if same-named directory exists
    Given the GitHub API server:
      """
      get('/repos/rtomayko/ronn') {
        json :private => false,
             :name => 'ronn', :owner => { :login => 'rtomayko' },
             :permissions => { :push => false }
      }
      """
    And a directory named "rtomayko/ronn"
    When I successfully run `hub clone rtomayko/ronn`
    Then it should clone "git://github.com/rtomayko/ronn.git"
    And there should be no output

  Scenario: Clone a public repo with HTTPS
    Given HTTPS is preferred
    Given the GitHub API server:
      """
      get('/repos/rtomayko/ronn') {
        json :private => false,
             :name => 'ronn', :owner => { :login => 'rtomayko' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub clone rtomayko/ronn`
    Then it should clone "https://github.com/rtomayko/ronn.git"
    And there should be no output

  Scenario: Clone command aliased
    Given the GitHub API server:
      """
      get('/repos/rtomayko/ronn') {
        json :private => false,
             :name => 'ronn', :owner => { :login => 'rtomayko' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `git config --global alias.c "clone --bare"`
    And I successfully run `hub c rtomayko/ronn`
    Then "git clone --bare git://github.com/rtomayko/ronn.git" should be run
    And there should be no output

  Scenario: Unchanged public clone
    When I successfully run `hub clone git://github.com/rtomayko/ronn.git`
    Then the git command should be unchanged

  Scenario: Unchanged public clone with path
    When I successfully run `hub clone git://github.com/rtomayko/ronn.git ronnie`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Unchanged private clone
    When I successfully run `hub clone git@github.com:rtomayko/ronn.git`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Unchanged clone with complex arguments
    When I successfully run `hub clone --template=one/two git://github.com/defunkt/resque.git --origin master resquetastic`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Unchanged local clone
    When I successfully run `hub clone ./dotfiles`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Unchanged local clone with destination
    Given a directory named ".git"
    When I successfully run `hub clone -l . ../copy`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Unchanged local clone from bare repo
    Given a bare git repo in "rtomayko/ronn"
    When I successfully run `hub clone rtomayko/ronn`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Unchanged clone with host alias
    When I successfully run `hub clone shortcut:git/repo.git`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Preview cloning a private repo
    Given the GitHub API server:
      """
      get('/repos/rtomayko/ronn') {
        json :private => false,
             :name => 'ronn', :owner => { :login => 'rtomayko' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub --noop clone -p rtomayko/ronn`
    Then the output should contain exactly "git clone git@github.com:rtomayko/ronn.git\n"
    But it should not clone anything

  Scenario: Clone a private repo
    Given the GitHub API server:
      """
      get('/repos/rtomayko/ronn') {
        json :private => false,
             :name => 'ronn', :owner => { :login => 'rtomayko' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub clone -p rtomayko/ronn`
    Then it should clone "git@github.com:rtomayko/ronn.git"
    And there should be no output

  Scenario: Clone my repo
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => true }
      }
      """
    When I successfully run `hub clone dotfiles`
    Then it should clone "git@github.com:mislav/dotfiles.git"
    And there should be no output

  Scenario: Clone my repo that doesn't exist
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') { status 404 }
      """
    When I run `hub clone dotfiles`
    Then the exit status should be 1
    And the stdout should contain exactly ""
    And the stderr should contain exactly "Error: repository mislav/dotfiles doesn't exist\n"
    And it should not clone anything

  Scenario: Clone my repo with arguments
    Given the GitHub API server:
      """
      get('/repos/mislav/dotfiles') {
        json :private => false,
             :name => 'dotfiles', :owner => { :login => 'mislav' },
             :permissions => { :push => true }
      }
      """
    When I successfully run `hub clone --bare -o master dotfiles`
    Then "git clone --bare -o master git@github.com:mislav/dotfiles.git" should be run
    And there should be no output

  Scenario: Clone repo to which I have push access to
    Given the GitHub API server:
      """
      get('/repos/sstephenson/rbenv') {
        json :private => false,
             :name => 'rbenv', :owner => { :login => 'sstephenson' },
             :permissions => { :push => true }
      }
      """
    When I successfully run `hub clone sstephenson/rbenv`
    Then "git clone git@github.com:sstephenson/rbenv.git" should be run
    And there should be no output

  Scenario: Preview cloning a repo I have push access to
    Given the GitHub API server:
      """
      get('/repos/sstephenson/rbenv') {
        json :private => false,
             :name => 'rbenv', :owner => { :login => 'sstephenson' },
             :permissions => { :push => true }
      }
      """
    When I successfully run `hub --noop clone sstephenson/rbenv`
    Then the output should contain exactly "git clone git@github.com:sstephenson/rbenv.git\n"
    But it should not clone anything

  Scenario: Clone my Enterprise repo
    Given I am "mifi" on git.my.org with OAuth token "FITOKEN"
    And $GITHUB_HOST is "git.my.org"
    Given the GitHub API server:
      """
      get('/api/v3/repos/myorg/myrepo') {
        json :private => true,
             :name => 'myrepo', :owner => { :login => 'myorg' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub clone myorg/myrepo`
    Then it should clone "git@git.my.org:myorg/myrepo.git"
    And there should be no output

  Scenario: Clone from existing directory is a local clone
    Given a directory named "dotfiles/.git"
    When I successfully run `hub clone dotfiles`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Clone from git bundle is a local clone
    Given a git bundle named "my-bundle"
    When I successfully run `hub clone my-bundle`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Clone a wiki
    Given the GitHub API server:
      """
      get('/repos/rtomayko/ronn') {
        json :private => false,
             :name => 'ronin', :owner => { :login => 'RTomayko' },
             :permissions => { :push => false },
             :has_wiki => true
      }
      """
    When I successfully run `hub clone rtomayko/ronn.wiki`
    Then it should clone "git://github.com/RTomayko/ronin.wiki.git"
    And there should be no output

  Scenario: Clone a nonexisting wiki
    Given the GitHub API server:
      """
      get('/repos/rtomayko/ronn') {
        json :private => false,
             :name => 'ronin', :owner => { :login => 'RTomayko' },
             :permissions => { :push => false },
             :has_wiki => false
      }
      """
    When I run `hub clone rtomayko/ronn.wiki`
    Then the exit status should be 1
    And the stdout should contain exactly ""
    And the stderr should contain exactly "Error: RTomayko/ronin doesn't have a wiki\n"
    And it should not clone anything

  Scenario: Clone a redirected repo
    Given the GitHub API server:
      """
      get('/repos/rtomayko/ronn') {
        redirect 'https://api.github.com/repositories/12345', 301
      }
      get('/repositories/12345', :host_name => 'api.github.com') {
        halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN'
        json :private => false,
             :name => 'ronin', :owner => { :login => 'RTomayko' },
             :permissions => { :push => false }
      }
      """
    When I successfully run `hub clone rtomayko/ronn`
    Then it should clone "git://github.com/RTomayko/ronin.git"
    And there should be no output
