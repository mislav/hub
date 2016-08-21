Feature: hub remote add
  Background:
    Given I am "EvilChelu" on GitHub.com
    And I am in "dotfiles" git repo

  Scenario: Add origin remote for my own repo
    Given there are no remotes
    When I successfully run `hub remote add origin`
    Then the url for "origin" should be "git@github.com:EvilChelu/dotfiles.git"
    And there should be no output

  Scenario: Add origin remote for my own repo using -C
    Given there are no remotes
    And I cd to ".."
    When I successfully run `hub -C dotfiles remote add origin`
    And I cd to "dotfiles"
    Then the url for "origin" should be "git@github.com:EvilChelu/dotfiles.git"
    And there should be no output

  Scenario: Unchanged public remote add
    When I successfully run `hub remote add origin http://github.com/defunkt/resque.git`
    Then the url for "origin" should be "http://github.com/defunkt/resque.git"
    And there should be no output

  Scenario: Unchanged private remote add
    When I successfully run `hub remote add origin git@github.com:defunkt/resque.git`
    Then the url for "origin" should be "git@github.com:defunkt/resque.git"
    And there should be no output

  Scenario: Unchanged local path remote add
    When I successfully run `hub remote add myremote ./path`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Unchanged local absolute path remote add
    When I successfully run `hub remote add myremote /path`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Unchanged remote add with host alias
    When I successfully run `hub remote add myremote server:/git/repo.git`
    Then the git command should be unchanged
    And there should be no output

  Scenario: Add new remote for Enterprise repo
    Given "git.my.org" is a whitelisted Enterprise host
    And I am "ProLoser" on git.my.org with OAuth token "FITOKEN"
    And the "origin" remote has url "git@git.my.org:mislav/topsekrit.git"
    When I successfully run `hub remote add another`
    Then the url for "another" should be "git@git.my.org:another/topsekrit.git"
    And there should be no output

  Scenario: Add public remote
    When I successfully run `hub remote add mislav`
    Then the url for "mislav" should be "git://github.com/mislav/dotfiles.git"
    And there should be no output

  Scenario: Add private remote
    When I successfully run `hub remote add -p mislav`
    Then the url for "mislav" should be "git@github.com:mislav/dotfiles.git"
    And there should be no output

  Scenario: Remote for my own repo is automatically private
    When I successfully run `hub remote add evilchelu`
    Then the url for "evilchelu" should be "git@github.com:EvilChelu/dotfiles.git"
    And there should be no output

  Scenario: Add remote with arguments
    When I successfully run `hub remote add -f mislav`
    Then "git remote add -f mislav git://github.com/mislav/dotfiles.git" should be run
    And there should be no output

  Scenario: Add HTTPS protocol remote
    Given HTTPS is preferred
    When I successfully run `hub remote add mislav`
    Then the url for "mislav" should be "https://github.com/mislav/dotfiles.git"
    And there should be no output

  Scenario: Add named public remote
    When I successfully run `hub remote add mm mislav`
    Then the url for "mm" should be "git://github.com/mislav/dotfiles.git"
    And there should be no output

  Scenario: set-url
    Given the "origin" remote has url "git://github.com/evilchelu/dotfiles.git"
    When I successfully run `hub remote set-url origin mislav`
    Then the url for "origin" should be "git://github.com/mislav/dotfiles.git"
    And there should be no output

  Scenario: Add public remote including repo name
    When I successfully run `hub remote add mislav/dotfilez.js`
    Then the url for "mislav" should be "git://github.com/mislav/dotfilez.js.git"
    And there should be no output

  Scenario: Add named public remote including repo name
    When I successfully run `hub remote add mm mislav/dotfilez.js`
    Then the url for "mm" should be "git://github.com/mislav/dotfilez.js.git"
    And there should be no output

  Scenario: Add named private remote
    When I successfully run `hub remote add -p mm mislav`
    Then the url for "mm" should be "git@github.com:mislav/dotfiles.git"
    And there should be no output

  Scenario: Add private remote including repo name
    When I successfully run `hub remote add -p mislav/dotfilez.js`
    Then the url for "mislav" should be "git@github.com:mislav/dotfilez.js.git"
    And there should be no output

  Scenario: Add named private remote including repo name
    When I successfully run `hub remote add -p mm mislav/dotfilez.js`
    Then the url for "mm" should be "git@github.com:mislav/dotfilez.js.git"
    And there should be no output

  Scenario: Add named private remote for my own repo including repo name
    When I successfully run `hub remote add ec evilchelu/dotfilez.js`
    Then the url for "ec" should be "git@github.com:EvilChelu/dotfilez.js.git"
    And there should be no output
