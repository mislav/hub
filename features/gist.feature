Feature: hub gist
  Background:
    Given I am "octokitten" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch a gist with a single file
    Given the GitHub API server:
      """
      get('/gists/myhash') {
        json({
          :files => {
            'hub_gist1.txt' => {
              'content' => "my content is here",
            }
          },
          :description => "my gist",
        })
      }
      """
    When I successfully run `hub gist show myhash`
    Then the output should contain exactly:
      """
      my content is here\n
      """
  
  Scenario: Fetch a gist with many files
    Given the GitHub API server:
      """
      get('/gists/myhash') {
        json({
          :files => {
            'hub_gist1.txt' => {
              'content' => "my content is here"
            },
            'hub_gist2.txt' => {
              'content' => "more content is here"
            }
          },
          :description => "my gist",
          :id => "myhash",
        })
      }
      """
    When I run `hub gist show myhash`
    Then the exit status should be 1
    Then the stderr should contain:
      """
      This gist contains multiple files, you must specify one:
        hub_gist1.txt
        hub_gist2.txt
      """

  Scenario: Fetch a single file from gist
    Given the GitHub API server:
      """
      get('/gists/myhash') {
        json({
          :files => {
            'hub_gist1.txt' => {
              'content' => "my content is here"
            },
            'hub_gist2.txt' => {
              'content' => "more content is here"
            }
          },
          :description => "my gist",
          :id => "myhash",
        })
      }
      """
    When I successfully run `hub gist show myhash hub_gist1.txt`
    Then the output should contain exactly:
      """
      my content is here\n
      """

  Scenario: Create a gist from file
    Given the GitHub API server:
      """
      post('/gists') {
        status 201
        json :html_url => 'http://gists.github.com/somehash'
      }
      """
    Given a file named "testfile.txt" with:
      """
      this is a test file
      """
    When I successfully run `hub gist create testfile.txt`
    Then the output should contain exactly:
      """
      http://gists.github.com/somehash\n
      """

  Scenario: Open the new gist in a browser
    Given the GitHub API server:
      """
      post('/gists') {
        status 201
        json :html_url => 'http://gists.github.com/somehash'
      }
      """
    Given a file named "testfile.txt" with:
      """
      this is a test file
      """
    When I successfully run `hub gist create -o testfile.txt`
    Then the output should contain exactly ""
    And "open http://gists.github.com/somehash" should be run

  Scenario: Create a gist with multiple files
    Given the GitHub API server:
      """
      post('/gists') {
        halt 400 unless params[:files]["testfile.txt"]["content"]
        halt 400 unless params[:files]["testfile2.txt"]["content"]
        status 201
        json({
          :html_url => 'http://gists.github.com/somehash',
        })
      }
      """
    Given a file named "testfile.txt" with:
      """
      this is a test file
      """
    Given a file named "testfile2.txt" with:
      """
      this is another test file
      """
    When I successfully run `hub gist create testfile.txt testfile2.txt`
    Then the output should contain exactly:
      """
      http://gists.github.com/somehash\n
      """

  Scenario: Create a gist from stdin
    Given the GitHub API server:
      """
      post('/gists') {
        halt 400 unless params[:files]["gistfile1.txt"]["content"] == "hello\n"
        status 201
        json :html_url => 'http://gists.github.com/somehash'
      }
      """
    When I run `hub gist create` interactively
    And I pass in:
      """
      hello
      """
    Then the output should contain exactly:
      """
      http://gists.github.com/somehash\n
      """

  Scenario: Insufficient OAuth scopes
    Given the GitHub API server:
      """
      post('/gists') {
        status 404
        response.headers['x-accepted-oauth-scopes'] = 'gist'
        response.headers['x-oauth-scopes'] = 'repos'
        json({})
      }
      """
    Given a file named "testfile.txt" with:
      """
      this is a test file
      """
    When I run `hub gist create testfile.txt`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error creating gist: Not Found (HTTP 404)
      Your access token may have insufficient scopes. Visit http://github.com/settings/tokens
      to edit the 'hub' token and enable one of the following scopes: gist\n
      """

  Scenario: Infer correct OAuth scopes for gist
    Given the GitHub API server:
      """
      post('/gists') {
        status 404
        response.headers['x-oauth-scopes'] = 'repos'
        json({})
      }
      """
    Given a file named "testfile.txt" with:
      """
      this is a test file
      """
    When I run `hub gist create testfile.txt`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error creating gist: Not Found (HTTP 404)
      Your access token may have insufficient scopes. Visit http://github.com/settings/tokens
      to edit the 'hub' token and enable one of the following scopes: gist\n
      """

  Scenario: Create error
    Given the GitHub API server:
      """
      post('/gists') {
        status 404
        response.headers['x-accepted-oauth-scopes'] = 'gist'
        response.headers['x-oauth-scopes'] = 'repos, gist'
        json({})
      }
      """
    Given a file named "testfile.txt" with:
      """
      this is a test file
      """
    When I run `hub gist create testfile.txt`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Error creating gist: Not Found (HTTP 404)\n
      """

