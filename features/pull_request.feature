Feature: hub pull-request
  Background:
    Given I am in "git://github.com/mislav/coral.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"
    And the git commit editor is "vim"

  Scenario: Detached HEAD
    Given I am in detached HEAD
    When I run `hub pull-request`
    Then the stderr should contain "Aborted: not currently on any branch.\n"
    And the exit status should be 1

  Scenario: Non-GitHub repo
    Given the "origin" remote has url "mygh:Manganeez/repo.git"
    When I run `hub pull-request`
    Then the stderr should contain "Aborted: the origin remote doesn't point to a GitHub repository.\n"
    And the exit status should be 1

  Scenario: Create pull request respecting "insteadOf" configuration
    Given the "origin" remote has url "mygh:Manganeez/repo.git"
    When I successfully run `git config url."git@github.com:".insteadOf mygh:`
    Given the GitHub API server:
      """
      post('/repos/Manganeez/repo/pulls') {
        assert :base  => 'master',
               :head  => 'Manganeez:master',
               :title => 'here we go'
        json :html_url => "https://github.com/Manganeez/repo/pull/12"
      }
      """
    When I successfully run `hub pull-request -m "here we go"`
    Then the output should contain exactly "https://github.com/Manganeez/repo/pull/12\n"

  Scenario: With Unicode characters
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 if request.content_charset != 'utf-8'
        assert :title => 'ăéñøü'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m ăéñøü`
    Then the output should contain exactly "the://url\n"

  Scenario: Deprecated title argument
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 422 if params[:title] != 'mytitle'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request mytitle`
    Then the stderr should contain exactly:
      """
      hub: Specifying pull request title without a flag is deprecated.
      Please use one of `-m' or `-F' options.\n
      """
    And the stdout should contain exactly "the://url\n"

  Scenario: Deprecated title argument can't start with a dash
    When I run `hub pull-request -help`
    Then the stderr should contain "invalid argument: -help\n"
    And the exit status should be 1

  Scenario: Non-existing base
    Given the GitHub API server:
      """
      post('/repos/origin/coral/pulls') { 404 }
      """
    When I run `hub pull-request -b origin:master -m here`
    Then the exit status should be 1
    Then the stderr should contain:
      """
      Error creating pull request: Not Found (HTTP 404)
      Are you sure that github.com/origin/coral exists?
      """

  Scenario: Supplies User-Agent string to API calls
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 unless request.user_agent.include?('Hub')
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m useragent`
    Then the output should contain exactly "the://url\n"

  Scenario: Text editor adds title and body
    Given the text editor adds:
      """
      This title comes from vim!

      This body as well.
      """
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'This title comes from vim!',
               :body  => 'This body as well.'
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I successfully run `hub pull-request`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Text editor adds title and body with multiple lines
    Given the text editor adds:
      """


      This title is on the third line


      This body


      has multiple
      lines.

      """
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'This title is on the third line',
               :body  => "This body\n\n\nhas multiple\nlines."
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I successfully run `hub pull-request`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Failed pull request preserves previous message
    Given the text editor adds:
      """
      This title will fail
      """
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 422 if params[:title].include?("fail")
        assert :body => "This title will fail",
               :title => "But this title will prevail"
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I run `hub pull-request`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error creating pull request: Unprocessable Entity (HTTP 422)\n
      """
    Given the text editor adds:
      """
      But this title will prevail
      """
    When I successfully run `hub pull-request`
    Then the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Ignore outdated PULLREQ_EDITMSG
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => "Added interactively", :body => nil
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    And a file named ".git/PULLREQ_EDITMSG" with:
      """
      Outdated message from old version of hub
      """
    Given the file named ".git/PULLREQ_EDITMSG" is older than hub source
    And the text editor adds:
      """
      Added interactively
      """
    When I successfully run `hub pull-request`
    Then the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Text editor fails
    Given the text editor exits with error status
    And an empty file named ".git/PULLREQ_EDITMSG"
    When I run `hub pull-request`
    Then the stderr should contain "error using text editor for pull request message"
    And the exit status should be 1
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Title and body from file
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'Title from file',
               :body  => "Body from file as well.\n\nMultiline, even!"
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    And a file named "pullreq-msg" with:
      """
      Title from file

      Body from file as well.

      Multiline, even!
      """
    When I successfully run `hub pull-request -F pullreq-msg`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Title and body from stdin
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'Unix piping is great',
               :body  => 'Just look at this'
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I run `hub pull-request -F -` interactively
    And I pass in:
      """
      Unix piping is great

      Just look at this
      """
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the exit status should be 0
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Title and body from command-line argument
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'I am just a pull',
               :body  => 'A little pull'
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I successfully run `hub pull-request -m "I am just a pull\n\nA little pull"`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Error when implicit head is the same as base
    Given I am on the "master" branch with upstream "origin/master"
    When I run `hub pull-request`
    Then the stderr should contain exactly:
      """
      Aborted: head branch is the same as base ("master")
      (use `-h <branch>` to specify an explicit pull request head)\n
      """

  Scenario: Explicit head
    Given I am on the "master" branch
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head => 'mislav:feature'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -h feature -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Explicit head with owner
    Given I am on the "master" branch
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head => 'mojombo:feature'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -h mojombo:feature -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Explicit base
    Given I am on the "feature" branch
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :base => 'develop'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -b develop -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Implicit base by detecting main branch
    Given the default branch for "origin" is "develop"
    And I am on the "master" branch
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :base => 'develop',
               :head => 'mislav:master'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Explicit base with owner
    Given I am on the "master" branch
    Given the GitHub API server:
      """
      post('/repos/mojombo/coral/pulls') {
        assert :base => 'develop'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -b mojombo:develop -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Explicit base with owner and repo name
    Given I am on the "master" branch
    Given the GitHub API server:
      """
      post('/repos/mojombo/coralify/pulls') {
        assert :base => 'develop'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -b mojombo/coralify:develop -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Error when there are unpushed commits
    Given I am on the "feature" branch with upstream "origin/feature"
    When I make 2 commits
    And I run `hub pull-request`
    Then the stderr should contain exactly:
      """
      Aborted: 2 commits are not yet pushed to origin/feature
      (use `-f` to force submit a pull request anyway)\n
      """

  Scenario: Ignore unpushed commits with `-f`
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head => 'mislav:feature'
        json :html_url => "the://url"
      }
      """
    When I make 2 commits
    And I successfully run `hub pull-request -f -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request fails on the server
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        status 422
        json(:message => "I haz fail!")
      }
      """
    When I run `hub pull-request -m message`
    Then the stderr should contain exactly:
      """
      Error creating pull request: Unprocessable Entity (HTTP 422)
      I haz fail!\n
      """

  Scenario: Convert issue to pull request
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :issue => '92'
        json :html_url => "https://github.com/mislav/coral/pull/92"
      }
      """
    When I successfully run `hub pull-request -i 92`
    Then the output should contain exactly:
      """
      https://github.com/mislav/coral/pull/92
      Warning: Issue to pull request conversion is deprecated and might not work in the future.\n
      """

  Scenario: Convert issue URL to pull request
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :issue => '92'
        json :html_url => "https://github.com/mislav/coral/pull/92"
      }
      """
    When I successfully run `hub pull-request https://github.com/mislav/coral/issues/92`
    Then the output should contain exactly:
      """
      https://github.com/mislav/coral/pull/92
      Warning: Issue to pull request conversion is deprecated and might not work in the future.\n
      """

  Scenario: Enterprise host
    Given the "origin" remote has url "git@git.my.org:mislav/coral.git"
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    Given the GitHub API server:
      """
      post('/api/v3/repos/mislav/coral/pulls') {
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m enterprisey`
    Then the output should contain exactly "the://url\n"

  Scenario: Create pull request from branch on the same remote
    Given the "origin" remote has url "git://github.com/github/coral.git"
    And the "mislav" remote has url "git://github.com/mislav/coral.git"
    And I am on the "feature" branch pushed to "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/github/coral/pulls') {
        assert :base  => 'master',
               :head  => 'github:feature',
               :title => 'hereyougo'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Create pull request from branch on the personal fork
    Given the "origin" remote has url "git://github.com/github/coral.git"
    And the "doge" remote has url "git://github.com/mislav/coral.git"
    And I am on the "feature" branch pushed to "doge/feature"
    Given the GitHub API server:
      """
      post('/repos/github/coral/pulls') {
        assert :base  => 'master',
               :head  => 'mislav:feature',
               :title => 'hereyougo'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Create pull request to "upstream" remote
    Given the "upstream" remote has url "git://github.com/github/coral.git"
    And I am on the "master" branch pushed to "origin/master"
    Given the GitHub API server:
      """
      post('/repos/github/coral/pulls') {
        assert :base  => 'master',
               :head  => 'mislav:master',
               :title => 'hereyougo'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Open pull request in web browser
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -o -m hereyougo`
    Then "open the://url" should be run
