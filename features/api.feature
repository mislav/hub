Feature: hub api
  Background:
    Given I am "octokitten" on github.com with OAuth token "OTOKEN"

  Scenario: GET resource
    Given the GitHub API server:
      """
      get('/hello/world') {
        halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN'
        halt 401 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3+json;charset=utf-8'
        json :name => "Ed"
      }
      """
    When I successfully run `hub api hello/world`
    Then the output should contain exactly:
      """
      {"name":"Ed"}\n
      """

  Scenario: GET Enterprise resource
    Given I am "octokitten" on git.my.org with OAuth token "FITOKEN"
    Given the GitHub API server:
      """
      get('/api/v3/hello/world', :host_name => 'git.my.org') {
        halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token FITOKEN'
        json :name => "Ed"
      }
      """
    And $GITHUB_HOST is "git.my.org"
    When I successfully run `hub api hello/world`
    Then the output should contain exactly:
      """
      {"name":"Ed"}\n
      """

  Scenario: Non-success response
    Given the GitHub API server:
      """
      get('/hello/world') {
        status 400
        json :name => "Ed"
      }
      """
    When I run `hub api hello/world`
    Then the exit status should be 1
    And the stdout should contain exactly ""
    And the stderr should contain exactly:
      """
      Error: HTTP 400 Bad Request
      {"name":"Ed"}\n
      """

  Scenario: Non-success response flat output
    Given the GitHub API server:
      """
      get('/hello/world') {
        status 400
        json :name => "Ed"
      }
      """
    When I run `hub api -t hello/world`
    Then the exit status should be 1
    And the stdout should contain exactly ""
    And the stderr should contain exactly:
      """
      Error: HTTP 400 Bad Request
      .name	Ed\n
      """

  Scenario: Non-success response doesn't choke on non-JSON
    Given the GitHub API server:
      """
      get('/hello/world') {
        status 400
        content_type :text
        'Something went wrong'
      }
      """
    When I run `hub api -t hello/world`
    Then the exit status should be 1
    And the stdout should contain exactly ""
    And the stderr should contain exactly:
      """
      Error: HTTP 400 Bad Request
      Something went wrong\n
      """

  Scenario: GET query string
    Given the GitHub API server:
      """
      get('/hello/world') {
        json Hash[*params.sort.flatten]
      }
      """
    When I successfully run `hub api -XGET -Fname=Ed -Fnum=12 -Fbool=false -Fvoid=null hello/world`
    Then the output should contain exactly:
      """
      {"bool":"false","name":"Ed","num":"12","void":""}\n
      """

  Scenario: GET full URL
    Given the GitHub API server:
      """
      get('/hello/world', :host_name => 'api.github.com') {
        halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN'
        json :name => "Faye"
      }
      """
    When I successfully run `hub api https://api.github.com/hello/world`
    Then the output should contain exactly:
      """
      {"name":"Faye"}\n
      """

  Scenario: Avoid leaking token to a 3rd party
    Given the GitHub API server:
      """
      get('/hello/world', :host_name => 'example.com') {
        halt 401 unless request.env['HTTP_AUTHORIZATION'].nil?
        json :name => "Jet"
      }
      """
    When I successfully run `hub api http://example.com/hello/world`
    Then the output should contain exactly:
      """
      {"name":"Jet"}\n
      """

  Scenario: Custom headers
    Given the GitHub API server:
      """
      get('/hello/world') {
        json :accept => request.env['HTTP_ACCEPT'],
             :foo => request.env['HTTP_X_FOO']
      }
      """
      When I successfully run `hub api hello/world -H 'x-foo:bar' -H 'Accept: text/json'`
    Then the output should contain exactly:
      """
      {"accept":"text/json","foo":"bar"}\n
      """

  Scenario: POST fields
    Given the GitHub API server:
      """
      post('/hello/world') {
        json Hash[*params.sort.flatten]
      }
      """
    When I successfully run `hub api -f name=@hubot -Fnum=12 -Fbool=false -Fvoid=null hello/world`
    Then the output should contain exactly:
      """
      {"bool":false,"name":"@hubot","num":12,"void":null}\n
      """

  Scenario: POST raw fields
    Given the GitHub API server:
      """
      post('/hello/world') {
        json Hash[*params.sort.flatten]
      }
      """
    When I successfully run `hub api -fnum=12 -fbool=false hello/world`
    Then the output should contain exactly:
      """
      {"bool":"false","num":"12"}\n
      """

  Scenario: POST from stdin
    Given the GitHub API server:
      """
      post('/graphql') {
        json :query => params[:query]
      }
      """
    When I run `hub api -t -F query=@- graphql` interactively
    And I pass in:
      """
      query {
        repository
      }
      """
    Then the output should contain exactly:
      """
      .query	query {\n  repository\n}\n\n
      """

  Scenario: Pass extra GraphQL variables
    Given the GitHub API server:
      """
      post('/graphql') {
        json(params[:variables])
      }
      """
    When I successfully run `hub api -F query='query {}' -Fname=Jet -Fsize=2 graphql`
    Then the output should contain exactly:
      """
      {"name":"Jet","size":2}\n
      """

  Scenario: Repo context
    Given I am in "git://github.com/octocat/Hello-World.git" git repo
    Given the GitHub API server:
      """
      get('/repos/octocat/Hello-World/commits') {
        json :commits => 12
      }
      """
    When I successfully run `hub api repos/{owner}/{repo}/commits`
    Then the output should contain exactly:
      """
      {"commits":12}\n
      """

  Scenario: Repo context in graphql
    Given I am in "git://github.com/octocat/Hello-World.git" git repo
    Given the GitHub API server:
      """
      post('/graphql') {
        json :query => params[:query]
      }
      """
    When I run `hub api -t -F query=@- graphql` interactively
    And I pass in:
      """
      repository(owner: "{owner}", name: "{repo}")
      """
    Then the output should contain exactly:
      """
      .query	repository(owner: "octocat", name: "Hello-World")\n\n
      """

  Scenario: Cache response
    Given the GitHub API server:
      """
      count = 0
      get('/count') {
        count += 1
        json :count => count
      }
      """
    When I successfully run `hub api -t 'count?a=1&b=2' --cache 5`
    And I successfully run `hub api -t 'count?b=2&a=1' --cache 5`
    Then the output should contain exactly:
      """
      .count	1
      .count	1\n
      """

  Scenario: Cache graphql response
    Given the GitHub API server:
      """
      count = 0
      post('/graphql') {
        halt 400 unless params[:query] =~ /^Q\d$/
        count += 1
        json :count => count
      }
      """
    When I successfully run `hub api -t graphql -F query=Q1 --cache 5`
    And I successfully run `hub api -t graphql -F query=Q1 --cache 5`
    And I successfully run `hub api -t graphql -F query=Q2 --cache 5`
    Then the output should contain exactly:
      """
      .count	1
      .count	1
      .count	2\n
      """

  Scenario: Avoid caching unsucessful response
    Given the GitHub API server:
      """
      count = 0
      get('/count') {
        count += 1
        status 400 if count == 1
        json :count => count
      }
      """
    When I run `hub api -t count --cache 5`
    And I successfully run `hub api -t count --cache 5`
    And I successfully run `hub api -t count --cache 5`
    Then the output should contain exactly:
      """
      .count	2
      .count	2
      Error: HTTP 400 Bad Request
      .count	1\n
      """

  Scenario: Avoid caching response if the OAuth token changes
    Given the GitHub API server:
      """
      count = 0
      get('/count') {
        count += 1
        json :count => count
      }
      """
    When I successfully run `hub api -t count --cache 5`
    Given I am "octocat" on github.com with OAuth token "TOKEN2"
    When I successfully run `hub api -t count --cache 5`
    Then the output should contain exactly:
      """
      .count	1
      .count	2\n
      """
