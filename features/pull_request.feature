Feature: hub pull-request
  Background:
    Given I am in "dotfiles" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Detached HEAD
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    And I am in detached HEAD
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
        { :base  => 'master',
          :head  => 'mislav:master',
          :title => 'here we go'
        }.each do |param, value|
          if params[param] != value
            halt 422, json(
              :message => "expected %s to be %s; got %s" % [
                param.inspect,
                value.inspect,
                params[param].inspect
              ]
            )
          end
        end
        json :html_url => "https://github.com/Manganeez/repo/pull/12"
      }
      """
    When I successfully run `hub pull-request -m "here we go"`
    Then the output should contain exactly "https://github.com/Manganeez/repo/pull/12\n"

  Scenario: With Unicode characters
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 if request.content_charset != 'utf-8'
        halt 422 if params[:title] != 'ăéñøü'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request -m ăéñøü`
    Then the output should contain exactly "the://url\n"

  Scenario: Deprecated title argument
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
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

  Scenario: Non-existing base
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
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
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
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
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    And the text editor adds:
      """
      This title comes from vim!

      This body as well.
      """
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        { :title => 'This title comes from vim!',
          :body  => 'This body as well.'
        }.each do |param, value|
          if params[param] != value
            halt 422, json(
              :message => "expected %s to be %s; got %s" % [
                param.inspect,
                value.inspect,
                params[param].inspect
              ]
            )
          end
        end
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I successfully run `hub pull-request`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Failed pull request preserves previous message
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    And the text editor adds:
      """
      This title will fail
      """
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 422 if params[:title].include?("fail")
        halt 422 unless params[:body] == "This title will fail"
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I run `hub pull-request`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error creating pull request:  (HTTP 422)\n
      """
    Given the text editor adds:
      """
      But this title will prevail
      """
    When I successfully run `hub pull-request`
    Then the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Text editor fails
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    And the text editor exits with error status
    And an empty file named ".git/PULLREQ_EDITMSG"
    When I run `hub pull-request`
    Then the stderr should contain "error using text editor for pull request message"
    And the exit status should be 1
    And the file ".git/PULLREQ_EDITMSG" should not exist

  Scenario: Title and body from file
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        { :title => 'Title from file',
          :body  => "Body from file as well.\n\nMultiline, even!"
        }.each do |param, value|
          if params[param] != value
            halt 422, json(
              :message => "expected %s to be %s; got %s" % [
                param.inspect,
                value.inspect,
                params[param].inspect
              ]
            )
          end
        end
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
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        { :title => 'Unix piping is great',
          :body  => 'Just look at this'
        }.each do |param, value|
          if params[param] != value
            halt 422, json(
              :message => "expected %s to be %s; got %s" % [
                param.inspect,
                value.inspect,
                params[param].inspect
              ]
            )
          end
        end
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
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        { :title => 'I am just a pull',
          :body  => 'A little pull'
        }.each do |param, value|
          if params[param] != value
            halt 422, json(
              :message => "expected %s to be %s; got %s" % [
                param.inspect,
                value.inspect,
                params[param].inspect
              ]
            )
          end
        end
        json :html_url => "https://github.com/mislav/coral/pull/12"
      }
      """
    When I successfully run `hub pull-request -m "I am just a pull\n\nA little pull"`
    Then the output should contain exactly "https://github.com/mislav/coral/pull/12\n"
    And the file ".git/PULLREQ_EDITMSG" should not exist
