Feature: hub pull-request
  Background:
    Given I am in "git://github.com/mislav/coral.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"
    And the git commit editor is "vim"

  Scenario: Basic pull request
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      KNOWN_PARAMS = %w[title body base head draft issue maintainer_can_modify]
      post('/repos/mislav/coral/pulls') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.shadow-cat-preview+json;charset=utf-8'
        halt 400 unless request.user_agent.include?('Hub')
        halt 400 if (params.keys - KNOWN_PARAMS).any?
        assert :title => 'hello',
               :body => nil,
               :base => 'master',
               :head => 'mislav:topic',
               :maintainer_can_modify => true,
               :draft => nil,
               :issue => nil
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hello`
    Then the output should contain exactly "the://url\n"

  Scenario: Detached HEAD
    Given I am in detached HEAD
    When I run `hub pull-request`
    Then the stderr should contain "Aborted: not currently on any branch.\n"
    And the exit status should be 1

  Scenario: Detached HEAD with explicit head
    Given I am in detached HEAD
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head => 'mislav:feature'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -h feature -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Non-GitHub repo
    Given the "origin" remote has url "mygh:Manganeez/repo.git"
    When I run `hub pull-request`
    Then the stderr should contain exactly:
      """
      Aborted: could not find any git remote pointing to a GitHub repository\n
      """
    And the exit status should be 1

  Scenario: Create pull request respecting "insteadOf" configuration
    Given the "origin" remote has url "mygh:Manganeez/repo.git"
    When I successfully run `git config url."git@github.com:".insteadOf mygh:`
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/Manganeez/repo/pulls') {
        assert :base  => 'master',
               :head  => 'Manganeez:topic',
               :title => 'here we go'
        status 201
        json :html_url => "https://github.com/Manganeez/repo/pull/12"
      }
      """
    When I successfully run `hub pull-request -m "here we go"`
    Then the output should contain exactly "https://github.com/Manganeez/repo/pull/12\n"

  Scenario: With Unicode characters
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 if request.content_charset != 'utf-8'
        assert :title => 'ăéñøü'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m ăéñøü`
    Then the output should contain exactly "the://url\n"

  Scenario: Invalid flag
    When I run `hub pull-request -yelp`
    Then the stderr should contain "unknown shorthand flag: 'y' in -yelp\n"
    And the exit status should be 1

  Scenario: With Unicode characters in the changelog
    Given the text editor adds:
      """
      I <3 encodings
      """
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 if request.content_charset != 'utf-8'
        assert :title => 'I <3 encodings',
               :body => 'ăéñøü'
        status 201
        json :html_url => "the://url"
      }
      """
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message "ăéñøü"
    And the "topic" branch is pushed to "origin/topic"
    When I successfully run `hub pull-request`
    Then the output should contain exactly "the://url\n"

  Scenario: Default message for single-commit pull request
    Given the text editor adds:
      """
      """
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 if request.content_charset != 'utf-8'
        assert :title => 'This is somewhat of a longish title that does not get wrapped & references #1234',
               :body => 'Hello'
        status 201
        json :html_url => "the://url"
      }
      """
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message:
      """
      This is somewhat of a longish title that does not get wrapped & references #1234

      Hello
      Signed-off-by: NAME <email@example.com>
      Co-authored-by: NAME <email@example.com>
      """
    And the "topic" branch is pushed to "origin/topic"
    When I successfully run `hub pull-request`
    Then the output should contain exactly "the://url\n"

  Scenario: Single-commit with pull request template
    Given the git commit editor is "true"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 if request.content_charset != 'utf-8'
        assert :title => 'Commit title',
               :body => <<BODY.chomp
      Commit body


       This is the pull request template

      Another line of template
      BODY
        status 201
        json :html_url => "the://url"
      }
      """
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    And I make a commit with message:
      """
      Commit title

      Commit body
      """
    And the "topic" branch is pushed to "origin/topic"
    Given a file named "pull_request_template.md" with:
      """
       This is the pull request template

      Another line of template
      """
    When I successfully run `hub pull-request`
    Then the output should contain exactly "the://url\n"

  Scenario: Single-commit with PULL_REQUEST_TEMPLATE directory
    Given the git commit editor is "true"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'Commit title',
               :body => 'Commit body'
        status 201
        json :html_url => "the://url"
      }
      """
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    And I make a commit with message:
      """
      Commit title

      Commit body
      """
    And the "topic" branch is pushed to "origin/topic"
    And a directory named "PULL_REQUEST_TEMPLATE"
    When I successfully run `hub pull-request`
    Then the output should contain exactly "the://url\n"

  Scenario: Single-commit pull request with "--no-edit"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'Commit title 1',
               :body => 'Commit body 1'
        status 201
        json :html_url => "the://url"
      }
      """
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message:
      """
      Commit title 1

      Commit body 1
      """
    And the "topic" branch is pushed to "origin/topic"
    When I successfully run `hub pull-request --no-edit`
    Then the output should contain exactly "the://url\n"

  Scenario: Multiple-commit pull request with "--no-edit"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'Commit title 1',
               :body => 'Commit body 1'
        status 201
        json :html_url => "the://url"
      }
      """
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message:
      """
      Commit title 1

      Commit body 1
      """
    Given I make a commit with message:
      """
      Commit title 2

      Commit body 2
      """
    And the "topic" branch is pushed to "origin/topic"
    When I successfully run `hub pull-request --no-edit`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request with "--push" and "--no-edit"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'Commit title 1',
               :body => 'Commit body 1'
        status 201
        json :html_url => "the://url"
      }
      """
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message:
      """
      Commit title 1

      Commit body 1
      """
    When I successfully run `hub pull-request --push --no-edit`
    Then the output should contain exactly "the://url\n"

  Scenario: No commits with "--no-edit"
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given the "topic" branch is pushed to "origin/topic"
    When I run `hub pull-request --no-edit`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Aborted: no commits detected between origin/master and topic\n
      """

  Scenario: Message template should include git log summary between base and head
    Given the text editor adds:
      """
      Hello
      """
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        status 500
      }
      """
    Given I am on the "master" branch
    And I make a commit with message "One on master"
    And I make a commit with message "Two on master"
    And the "master" branch is pushed to "origin/master"
    Given I successfully run `git reset --hard HEAD~2`
    And I successfully run `git checkout --quiet -B topic origin/master`
    Given I make a commit with message "One on topic"
    And I make a commit with message "Two on topic"
    Given the "topic" branch is pushed to "origin/topic"
    And I successfully run `git reset --hard HEAD~1`
    When I run `hub pull-request`
    Given the SHAs and timestamps are normalized in ".git/PULLREQ_EDITMSG"
    Then the file ".git/PULLREQ_EDITMSG" should contain exactly:
      """
      Hello


      # ------------------------ >8 ------------------------
      # Do not modify or remove the line above.
      # Everything below it will be ignored.

      Requesting a pull to mislav:master from mislav:topic

      Write a message for this pull request. The first block
      of text is the title and the rest is the description.

      Changes:

      SHA1SHA (Hub, 0 seconds ago)
         Two on topic

      SHA1SHA (Hub, 0 seconds ago)
         One on topic

      """

  Scenario: Non-existing base
    Given the GitHub API server:
      """
      post('/repos/origin/coral/pulls') { 404 }
      """
    When I run `hub pull-request -b origin:master -h topic -m here`
    Then the exit status should be 1
    Then the stderr should contain:
      """
      Error creating pull request: Not Found (HTTP 404)
      Are you sure that github.com/origin/coral exists?
      """

  Scenario: Text editor adds title and body
    Given I am on the "topic" branch pushed to "origin/topic"
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
        status 201
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I successfully run `hub pull-request`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Text editor adds title and body with multiple lines
    Given I am on the "topic" branch pushed to "origin/topic"
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
        status 201
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I successfully run `hub pull-request`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Text editor with custom commentchar
    Given I am on the "topic" branch pushed to "origin/topic"
    Given git "core.commentchar" is set to "/"
    And the text editor adds:
      """
      # Dat title

      / This line is not commented out.

      Dem body.
      """
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => '# Dat title',
               :body  => "/ This line is not commented out.\n\nDem body."
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request`
    Then the output should contain exactly "the://url\n"

  Scenario: Failed pull request preserves previous message
    Given I am on the "topic" branch pushed to "origin/topic"
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
        status 201
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

  Scenario: Text editor fails
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the text editor exits with error status
    And an empty file named ".git/PULLREQ_EDITMSG"
    When I run `hub pull-request`
    Then the stderr should contain "error using text editor for pull request message"
    And the exit status should be 1
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Title and body from file
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'Title from file',
               :body  => "Body from file as well.\n\nMultiline, even!"
        status 201
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

  Scenario: Edit title and body from file
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'Hello from editor',
               :body  => "Title from file\n\nBody from file as well."
        status 201
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    And a file named "pullreq-msg" with:
      """
      Title from file

      Body from file as well.
      """
    And the text editor adds:
      """
      Hello from editor
      """
    When I successfully run `hub pull-request -F pullreq-msg --edit`
    Then the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Title and body from stdin
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'Unix piping is great',
               :body  => 'Just look at this ăéñøü'
        status 201
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I run `hub pull-request -F -` interactively
    And I pass in:
      """
      Unix piping is great

      Just look at this ăéñøü
      """
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the exit status should be 0
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Title and body from command-line argument
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'I am just a pull',
               :body  => 'A little pull'
        status 201
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I successfully run `hub pull-request -m "I am just a pull\n\nA little pull"`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Title and body from multiple command-line arguments
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'I am just a pull',
               :body  => "A little pull\n\nAnd description"
        status 201
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I successfully run `hub pull-request -m "I am just a pull" -m "A little pull" -m "And description"`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"

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
        status 201
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
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -h mojombo:feature -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Explicit base
    Given I am on the "feature" branch pushed to "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :base => 'develop',
               :head => 'mislav:feature'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -b develop -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Implicit base by detecting main branch
    Given the default branch for "origin" is "develop"
    And the "master" branch is pushed to "origin/master"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :base => 'develop',
               :head => 'mislav:master'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Explicit base with owner
    Given I am on the "master" branch pushed to "origin/master"
    Given the GitHub API server:
      """
      post('/repos/mojombo/coral/pulls') {
        assert :base => 'develop'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -b mojombo:develop -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Explicit base with owner and repo name
    Given I am on the "master" branch pushed to "origin/master"
    Given the GitHub API server:
      """
      post('/repos/mojombo/coralify/pulls') {
        assert :base => 'develop'
        status 201
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
        status 201
        json :html_url => "the://url"
      }
      """
    When I make 2 commits
    And I successfully run `hub pull-request -f -m message`
    Then the output should contain exactly "the://url\n"

  Scenario: Error from an unpushed branch
    Given I am on the "feature" branch
    When I run `hub pull-request -m hello`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Aborted: the current branch seems not yet pushed to a remote
      (use `-p` to push the branch or `-f` to skip this check)\n
      """

  Scenario: Error from an unpushed branch with upstream same as base branch
    Given I am on the "feature" branch with upstream "origin/master"
    When I run `hub pull-request -m hello`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Aborted: the current branch seems not yet pushed to a remote
      (use `-p` to push the branch or `-f` to skip this check)\n
      """

  Scenario: Pull request fails on the server
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      tries = 0
      post('/repos/mislav/coral/pulls') {
        tries += 1
        if tries == 1
          status 422
          json :message => 'Validation Failed',
               :errors => [{
                 :resource => 'PullRequest',
                 :code => 'invalid',
                 :field => 'head'
               }]
        else
          status 400
        end
      }
      """
    When I run `hub pull-request -m message`
    Then the stderr should contain exactly:
      """
      Error creating pull request: Unprocessable Entity (HTTP 422)
      Invalid value for "head"\n
      """
    And the exit status should be 1

  Scenario: Convert issue to pull request
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :issue => 92
        status 201
        json :html_url => "https://github.com/mislav/coral/pull/92"
      }
      """
    When I successfully run `hub pull-request -i 92`
    Then the output should contain exactly:
      """
      https://github.com/mislav/coral/pull/92\n
      """

  Scenario: Convert issue URL to pull request
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :issue => 92
        status 201
        json :html_url => "https://github.com/mislav/coral/pull/92"
      }
      """
    When I successfully run `hub pull-request https://github.com/mislav/coral/issues/92`
    Then the output should contain exactly:
      """
      https://github.com/mislav/coral/pull/92\n
      """

  Scenario: Enterprise host
    Given the "origin" remote has url "git@git.my.org:mislav/coral.git"
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    And I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/api/v3/repos/mislav/coral/pulls', :host_name => 'git.my.org') {
        assert :base => 'master',
               :head => 'mislav:topic'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m enterprisey`
    Then the output should contain exactly "the://url\n"

  Scenario: Enterprise remote witch matching branch but no tracking
    Given the "origin" remote has url "git@git.my.org:mislav/coral.git"
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    And I am on the "feature" branch pushed to "origin/feature"
    Given the GitHub API server:
      """
      post('/api/v3/repos/mislav/coral/pulls', :host_name => 'git.my.org') {
        assert :base => 'master',
               :head => 'mislav:feature'
        status 201
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
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Create pull request from branch on the personal fork case sensitive
    Given the "origin" remote has url "git://github.com/github/coral.git"
    And the "doge" remote has url "git://github.com/Mislav/coral.git"
    And I am on the "feature" branch pushed to "doge/feature"
    Given the GitHub API server:
      """
      post('/repos/github/coral/pulls') {
        assert :base  => 'master',
               :head  => 'Mislav:feature',
               :title => 'hereyougo'
        status 201
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
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Create pull request from branch on the personal fork, capitalized
    Given the "origin" remote has url "git://github.com/LightAlf/FirstRepo.git"
    And the "Kristinita" remote has url "git@github.com:Kristinita/FirstRepo.git"
    And I am on the "add-py3kwarn" branch pushed to "Kristinita/add-py3kwarn"
    And I am "Kristinita" on github.com with OAuth token "OTOKEN"
    Given the GitHub API server:
      """
      post('/repos/LightAlf/FirstRepo/pulls') {
        assert :base  => 'master',
               :head  => 'Kristinita:add-py3kwarn',
               :title => 'hereyougo'
        status 201
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
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Create pull request to "upstream" remote with differently-named default branch
    Given I am on the "master" branch pushed to "origin/master"
    And the "upstream" remote has url "git://github.com/github/coral.git"
    And the default branch for "upstream" is "develop"
    Given the GitHub API server:
      """
      post('/repos/github/coral/pulls') {
        assert :base  => 'develop',
               :head  => 'mislav:master',
               :title => 'hereyougo'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Create pull request to "github" remote when "upstream" is non-GitHub
    Given I am on the "master" branch pushed to "origin/master"
    And the "github" remote has url "git://github.com/github/coral.git"
    And the "upstream" remote has url "git://example.com/coral.git"
    Given the GitHub API server:
      """
      post('/repos/github/coral/pulls') {
        assert :base  => 'master',
               :head  => 'mislav:master',
               :title => 'hereyougo'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Create pull request to "github" remote when "origin" is non-GitHub
    Given the "github" remote has url "git@github.com:sam-hart-swanson/debug.git"
    Given the "origin" remote has url "ssh://git@private.server.com/path/to/repo.git"
    And I am on the "feat/123-some-branch" branch pushed to "github/feat/123-some-branch"
    And an empty file named ".git/refs/remotes/origin/feat/123-some-branch"
    Given the GitHub API server:
      """
      post('/repos/sam-hart-swanson/debug/pulls') {
        assert :base  => 'master',
               :head  => 'sam-hart-swanson:feat/123-some-branch',
               :title => 'hereyougo'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Open pull request in web browser
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -o -m hereyougo`
    Then "open the://url" should be run

  Scenario: Current branch is tracking local branch
    Given git "push.default" is set to "upstream"
    And I am on the "feature" branch pushed to "origin/feature"
    When I successfully run `git config branch.feature.remote .`
    When I successfully run `git config branch.feature.merge refs/heads/master`
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :base  => 'master',
               :head  => 'mislav:feature'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Current branch is pushed to remote without upstream configuration
    Given the "upstream" remote has url "git://github.com/lestephane/coral.git"
    And I am on the "feature" branch pushed to "origin/feature"
    And git "push.default" is set to "upstream"
    Given the GitHub API server:
      """
      post('/repos/lestephane/coral/pulls') {
        assert :base  => 'master',
               :head  => 'mislav:feature'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Branch with quotation mark in name
    Given I am on the "feat'ure" branch with upstream "origin/feat'ure"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head  => "mislav:feat'ure"
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request with assignees
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head  => "mislav:feature"
        status 201
        json :html_url => "the://url", :number => 1234
      }
      patch('/repos/mislav/coral/issues/1234') {
        assert :assignees => ["mislav", "josh", "pcorpet"], :labels => :no
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo -a mislav,josh -apcorpet`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request with reviewers
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head  => "mislav:feature"
        status 201
        json :html_url => "the://url", :number => 1234
      }
      post('/repos/mislav/coral/pulls/1234/requested_reviewers') {
        assert :reviewers => ["mislav", "josh", "pcorpet"]
        assert :team_reviewers => ["robots", "js"]
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo -r mislav,josh -rgithub/robots -rpcorpet -r github/js`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request with reviewers from CODEOWNERS
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head  => "mislav:feature"
        status 201
        json :html_url => "the://url", :number => 1234,
          :requested_reviewers => [{ :login => "josh" }],
          :requested_teams => [{ :slug => "robots" }]
      }
      post('/repos/mislav/coral/pulls/1234/requested_reviewers') {
        assert :reviewers => ["mislav", "pcorpet"]
        assert :team_reviewers => ["js"]
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo -r mislav,josh -rgithub/robots -rpcorpet -r github/js`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request avoids re-requesting reviewers
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head  => "mislav:feature"
        status 201
        json :html_url => "the://url", :number => 1234,
          :requested_reviewers => [{ :login => "josh" }, { :login => "mislav" }],
          :requested_teams => [{ :slug => "robots" }]
      }
      """
    When I successfully run `hub pull-request -m hereyougo -r mislav,josh -rgithub/robots`
    Then the output should contain exactly "the://url\n"

  Scenario: Requesting reviewers failed
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        status 201
        json :html_url => "the://url", :number => 1234
      }
      post('/repos/mislav/coral/pulls/1234/requested_reviewers') {
        status 422
        json :message => "Validation Failed",
          :errors => ["Could not add requested reviewers to pull request."],
          :documentation_url => "https://developer.github.com/v3/pulls/review_requests/#create-a-review-request"
      }
      """
    When I run `hub pull-request -m hereyougo -r pedrohc`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error requesting reviewer: Unprocessable Entity (HTTP 422)
      Could not add requested reviewers to pull request.\n
      """

  Scenario: Pull request with milestone
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      get('/repos/mislav/coral/milestones') {
        status 200
        json [
          { :number => 237, :title => "prerelease" },
          { :number => 1337, :title => "v1" },
          { :number => 41319, :title => "Hello World!" }
        ]
      }
      post('/repos/mislav/coral/pulls') {
        assert :head  => "mislav:feature"
        status 201
        json :html_url => "the://url", :number => 1234
      }
      patch('/repos/mislav/coral/issues/1234') {
        assert :milestone => 41319
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo -M "Hello World!"`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request with case-insensitive milestone
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      get('/repos/mislav/coral/milestones') {
        status 200
        json [
          { :number => 237, :title => "prerelease" },
          { :number => 1337, :title => "v1" },
          { :number => 41319, :title => "Hello World!" }
        ]
      }
      post('/repos/mislav/coral/pulls') {
        assert :head  => "mislav:feature"
        status 201
        json :html_url => "the://url", :number => 1234
      }
      patch('/repos/mislav/coral/issues/1234') {
        assert :milestone => 41319
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo -M "hello world!"`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request uses integer milestone number for BC
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      get('/repos/mislav/coral/milestones') {
        status 200
        json [{ :number => 237, :title => "prerelease" }]
      }
      post('/repos/mislav/coral/pulls') {
        assert :head  => "mislav:feature"
        status 201
        json :html_url => "the://url", :number => 1234
      }
      patch('/repos/mislav/coral/issues/1234') {
        assert :milestone => 55
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo -M 55`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request fails with unknown milestone before it's created
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      get('/repos/mislav/coral/milestones') {
        status 200
        json []
      }
      """
    When I run `hub pull-request -m hereyougo -M "unknown"`
    Then the exit status should be 1
    And the stderr should contain exactly "error: no milestone found with name 'unknown'\n"

  Scenario: Pull request with labels
    Given I am on the "feature" branch with upstream "origin/feature"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :head  => "mislav:feature"
        status 201
        json :html_url => "the://url", :number => 1234
      }
      patch('/repos/mislav/coral/issues/1234') {
        assert :labels => ["feature", "release", "docs"], :assignees => :no
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo -l feature,release -ldocs`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request to a fetch-only upstream
    Given the "upstream" remote has url "git://github.com/github/coral.git"
    And the "upstream" remote has push url "no_push"
    And I am on the "master" branch pushed to "origin/master"
    Given the GitHub API server:
      """
      post('/repos/github/coral/pulls') {
        assert :base  => 'master',
               :head  => 'mislav:master',
               :title => 'hereyougo'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request with 307 redirect
    Given the "origin" remote has url "https://github.com/mislav/coral.git"
    And I am on the "feature" branch pushed to "origin/feature"
    Given the GitHub API server:
      """
      get('/repos/mislav/coral') {
        redirect 'https://api.github.com/repositories/12345', 301
      }
      get('/repositories/12345') {
        json :name => 'coralify', :owner => { :login => 'coral-org' }
      }
      post('/repos/mislav/coral/pulls') {
        redirect 'https://api.github.com/repositories/12345', 307
      }
      post('/repositories/12345', :host_name => 'api.github.com') {
        assert :base  => 'master',
               :head  => 'coral-org:feature',
               :title => 'hereyougo'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hereyougo`
    Then the output should contain exactly "the://url\n"

  Scenario: Pull request with 301 redirect
    Given the "origin" remote has url "https://github.com/mislav/coral.git"
    And I am on the "feature" branch pushed to "origin/feature"
    Given the GitHub API server:
      """
      get('/repos/mislav/coral') {
        redirect 'https://api.github.com/repositories/12345', 301
      }
      get('/repositories/12345') {
        json :name => 'coralify', :owner => { :login => 'coral-org' }
      }
      post('/repos/mislav/coral/pulls') {
        redirect 'https://api.github.com/repositories/12345/pulls', 301
      }
      post('/repositories/12345/pulls', :host_name => 'api.github.com') {
        assert :base  => 'master',
               :head  => 'coral-org:feature',
               :title => 'hereyougo'
        status 201
        json :html_url => "the://url"
      }
      """
    When I run `hub pull-request -m hereyougo`
    Then the exit status should be 1
    And stderr should contain exactly:
      """
      Error creating pull request: Post https://api.github.com/repositories/12345/pulls: refusing to follow HTTP 301 redirect for a POST request
      Have your site admin use HTTP 308 for this kind of redirect\n
      """

  Scenario: Default message with --push
    Given the git commit editor is "true"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :title => 'The commit I never pushed',
               :body => nil
        status 201
        json :html_url => "the://url"
      }
      """
    Given I am on the "master" branch pushed to "origin/master"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message "The commit I never pushed"
    When I successfully run `hub pull-request -p`
    Then the output should contain exactly "the://url\n"
    And "git push --set-upstream origin HEAD:topic" should be run

  Scenario: Text editor fails with --push
    Given the text editor exits with error status
    And I am on the "master" branch pushed to "origin/master"
    And an empty file named ".git/PULLREQ_EDITMSG"
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit
    When I run `hub pull-request -p`
    Then the stderr should contain "error using text editor for pull request message"
    And the exit status should be 1
    And the file ".git/PULLREQ_EDITMSG" should not exist
    And "git push --set-upstream origin HEAD:topic" should not be run

  Scenario: Triangular workflow with --push
    Given the "upstream" remote has url "git://github.com/github/coral.git"
    And I am on the "master" branch pushed to "upstream/master"
    # TODO: head should be "mislav:topic"
    Given the GitHub API server:
      """
      post('/repos/github/coral/pulls') {
        assert :base  => 'master',
               :head  => 'github:topic',
               :title => 'hereyougo'
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `git checkout --quiet -b topic`
    Given I make a commit with message "Fork commit"
    When I successfully run `hub pull-request -p -m hereyougo`
    Then the output should contain exactly "the://url\n"
    # TODO: the push should be to the "origin" remote instead
    And "git push --set-upstream upstream HEAD:topic" should be run

  Scenario: Automatically retry when --push resulted in 422
    Given the default aruba exit timeout is 7 seconds
    And the text editor adds:
      """
      hello!
      """
    Given the GitHub API server:
      """
      first_try_at = nil
      tries = 0

      post('/repos/mislav/coral/pulls') {
        tries += 1
        assert :title => 'hello!', :head => 'mislav:topic'

        if !first_try_at || (Time.now - first_try_at) < 5
          first_try_at ||= Time.now
          status 422
          json :message => 'Validation Failed',
               :errors => [{
                 :resource => 'PullRequest',
                 :code => 'invalid',
                 :field => 'head'
               }]
        else
          status 201
          json :html_url => "the://url?tries=#{tries}"
        end
      }
      """
    Given I am on the "topic" branch
    When I successfully run `hub pull-request -p`
    Then the output should contain exactly "the://url?tries=3\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Eventually give up on retries for --push
    Given the default aruba exit timeout is 7 seconds
    And the text editor adds:
      """
      hello!
      """
    And $HUB_RETRY_TIMEOUT is "5"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        status 422
        json :message => 'Validation Failed',
             :errors => [{
               :resource => 'PullRequest',
               :code => 'invalid',
               :field => 'head'
             }]
      }
      """
    Given I am on the "topic" branch
    When I run `hub pull-request -p`
    Then the stderr should contain:
      """
      Error creating pull request: Unprocessable Entity (HTTP 422)
      Invalid value for "head"\n
      """
    And the output should match /Given up after retrying for 5\.\d seconds\./
    And a file named ".git/PULLREQ_EDITMSG" should exist

  Scenario: Draft pull request
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.shadow-cat-preview+json;charset=utf-8'
        assert :draft => true
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -d -m wip`
    Then the output should contain exactly "the://url\n"

  Scenario: Disallow edits from maintainers
    Given I am on the "topic" branch pushed to "origin/topic"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        assert :maintainer_can_modify => false
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m hello --no-maintainer-edits`
    Then the output should contain exactly "the://url\n"
