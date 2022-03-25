@cache_clear
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
      {"name":"Ed"}
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
      {"name":"Ed"}
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
    Then the exit status should be 22
    And the stdout should contain exactly:
      """
      {"name":"Ed"}
      """
    And the stderr should contain exactly ""

  Scenario: Non-success response flat output
    Given the GitHub API server:
      """
      get('/hello/world') {
        status 400
        json :name => "Ed"
      }
      """
    When I run `hub api -t hello/world`
    Then the exit status should be 22
    And the stdout should contain exactly:
      """
      .name	Ed\n
      """
    And the stderr should contain exactly ""

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
    Then the exit status should be 22
    And the stdout should contain exactly:
      """
      Something went wrong
      """
    And the stderr should contain exactly ""

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
      {"bool":"false","name":"Ed","num":"12","void":""}
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
      {"name":"Faye"}
      """

  Scenario: Paginate REST
    Given the GitHub API server:
      """
      get('/comments') {
        assert :per_page => "6"
        page = (params[:page] || 1).to_i
        response.headers["Link"] = %(<#{request.url}&page=#{page+1}>; rel="next") if page < 3
        json [{:page => page}]
      }
      """
    When I successfully run `hub api --paginate comments?per_page=6`
    Then the output should contain exactly:
      """
      [{"page":1}]
      [{"page":2}]
      [{"page":3}]
      """

  Scenario: Paginate GraphQL
    Given the GitHub API server:
      """
      post('/graphql') {
        variables = params[:variables] || {}
        page = (variables["endCursor"] || 1).to_i
        json :data => {
          :pageInfo => {
            :hasNextPage => page < 3,
            :endCursor => (page+1).to_s
          }
        }
      }
      """
    When I successfully run `hub api --paginate graphql -f query=QUERY`
    Then the output should contain exactly:
      """
      {"data":{"pageInfo":{"hasNextPage":true,"endCursor":"2"}}}
      {"data":{"pageInfo":{"hasNextPage":true,"endCursor":"3"}}}
      {"data":{"pageInfo":{"hasNextPage":false,"endCursor":"4"}}}
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
      {"name":"Jet"}
      """

  Scenario: Request headers
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
      {"accept":"text/json","foo":"bar"}
      """

  Scenario: Response headers
    Given the GitHub API server:
      """
      get('/hello/world') {
        json({})
      }
      """
    When I successfully run `hub api hello/world -i`
    Then the output should contain "HTTP/1.1 200 OK"
    And the output should contain "Content-Length: 2"

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
      {"bool":false,"name":"@hubot","num":12,"void":null}
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
      {"bool":"false","num":"12"}
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
      .query	query {\\n  repository\\n}\\n\n
      """

  Scenario: POST body from file
    Given the GitHub API server:
      """
      post('/create') {
        params[:obj].inspect
      }
      """
    Given a file named "payload.json" with:
      """
      {"obj": ["one", 2, null]}
      """
    When I successfully run `hub api create --input payload.json`
    Then the output should contain exactly:
      """
      ["one", 2, nil]
      """

  Scenario: POST body from stdin
    Given the GitHub API server:
      """
      post('/create') {
        params[:obj].inspect
      }
      """
    When I run `hub api create --input -` interactively
    And I pass in:
      """
      {"obj": {"name": "Ein", "datadog": true}}
      """
    Then the output should contain exactly:
      """
      {"name"=>"Ein", "datadog"=>true}
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
      {"name":"Jet","size":2}
      """

  Scenario: Enterprise GraphQL
    Given I am "octokitten" on git.my.org with OAuth token "FITOKEN"
    Given the GitHub API server:
      """
      post('/api/graphql', :host_name => 'git.my.org') {
        halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token FITOKEN'
        json :name => "Ed"
      }
      """
    And $GITHUB_HOST is "git.my.org"
    When I successfully run `hub api graphql -f query=QUERY`
    Then the output should contain exactly:
      """
      {"name":"Ed"}
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
      {"commits":12}
      """

  Scenario: Multiple string interpolation
    Given I am in "git://github.com/octocat/Hello-World.git" git repo
    Given the GitHub API server:
      """
      get('/repos/octocat/Hello-World/pulls') {
        json(params)
      }
      """
    When I successfully run `hub api repos/{owner}/{repo}/pulls?head={owner}:{repo}`
    Then the output should contain exactly:
      """
      {"head":"octocat:Hello-World"}
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
      repository(owner: "{owner}", name: "{repo}", nameWithOwner: "{owner}/{repo}")
      """
    Then the output should contain exactly:
      """
      .query	repository(owner: "octocat", name: "Hello-World", nameWithOwner: "octocat/Hello-World")\\n\n
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
    When I run `hub api -t 'count?a=1&b=2' --cache 5`
    Then it should pass with ".count	1"
    When I run `hub api -t 'count?b=2&a=1' --cache 5`
    Then it should pass with ".count	1"

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
    When I run `hub api -t graphql -F query=Q1 --cache 5`
    Then it should pass with ".count	1"
    When I run `hub api -t graphql -F query=Q1 --cache 5`
    Then it should pass with ".count	1"
    When I run `hub api -t graphql -F query=Q2 --cache 5`
    Then it should pass with ".count	2"

  Scenario: Cache client error response
    Given the GitHub API server:
      """
      count = 0
      get('/count') {
        count += 1
        status 404 if count == 1
        json :count => count
      }
      """
    When I run `hub api -t count --cache 5`
    Then it should fail with ".count	1"
    When I run `hub api -t count --cache 5`
    Then it should fail with ".count	1"
    And the exit status should be 22

  Scenario: Avoid caching server error response
    Given the GitHub API server:
      """
      count = 0
      get('/count') {
        count += 1
        status 500 if count == 1
        json :count => count
      }
      """
    When I run `hub api -t count --cache 5`
    Then it should fail with ".count	1"
    When I run `hub api -t count --cache 5`
    Then it should pass with ".count	2"
    When I run `hub api -t count --cache 5`
    Then it should pass with ".count	2"

  Scenario: Avoid caching response if the OAuth token changes
    Given the GitHub API server:
      """
      count = 0
      get('/count') {
        count += 1
        json :count => count
      }
      """
    When I run `hub api -t count --cache 5`
    Then it should pass with ".count	1"
    Given I am "octocat" on github.com with OAuth token "TOKEN2"
    When I run `hub api -t count --cache 5`
    Then it should pass with ".count	2"

  Scenario: Honor rate limit with pagination
    Given the GitHub API server:
      """
      get('/hello') {
        page = (params[:page] || 1).to_i
        if page < 2
          response.headers['X-Ratelimit-Remaining'] = '0'
          response.headers['X-Ratelimit-Reset'] = Time.now.utc.to_i.to_s
          response.headers['Link'] = %(</hello?page=#{page+1}>; rel="next")
        end
        json [{}]
      }
      """
    When I successfully run `hub api --obey-ratelimit --paginate hello`
    Then the stderr should contain "API rate limit exceeded; pausing until "

  Scenario: Succumb to rate limit with pagination
    Given the GitHub API server:
      """
      get('/hello') {
        page = (params[:page] || 1).to_i
        response.headers['X-Ratelimit-Remaining'] = '0'
        response.headers['X-Ratelimit-Reset'] = Time.now.utc.to_i.to_s
        if page == 2
          status 403
          json :message => "API rate limit exceeded"
        else
          response.headers['Link'] = %(</hello?page=#{page+1}>; rel="next")
          json [{page:page}]
        end
      }
      """
    When I run `hub api --paginate -t hello`
    Then the exit status should be 22
    And the stderr should not contain "API rate limit exceeded"
    And the stdout should contain exactly:
      """
      .[0].page	1
      .message	API rate limit exceeded\n
      """

  Scenario: Honor rate limit for 403s
    Given the GitHub API server:
      """
      count = 0
      get('/hello') {
        count += 1
        if count == 1
          response.headers['X-Ratelimit-Remaining'] = '0'
          response.headers['X-Ratelimit-Reset'] = Time.now.utc.to_i.to_s
          halt 403
        end
        json [{}]
      }
      """
    When I successfully run `hub api --obey-ratelimit hello`
    Then the stderr should contain "API rate limit exceeded; pausing until "

  Scenario: 403 unrelated to rate limit
    Given the GitHub API server:
      """
      get('/hello') {
        response.headers['X-Ratelimit-Remaining'] = '1'
        status 403
      }
      """
    When I run `hub api --obey-ratelimit hello`
    Then the exit status should be 22
    Then the stderr should not contain "API rate limit exceeded"

  Scenario: Warn about insufficient OAuth scopes
    Given the GitHub API server:
      """
      get('/hello') {
        response.headers['X-Accepted-Oauth-Scopes'] = 'repo, admin'
        response.headers['X-Oauth-Scopes'] = 'public_repo'
        status 403
        json({})
      }
      """
    When I run `hub api hello`
    Then the exit status should be 22
    And the output should contain exactly:
      """
      {}
      Your access token may have insufficient scopes. Visit http://github.com/settings/tokens
      to edit the 'hub' token and enable one of the following scopes: admin, repo\n
      """

  Scenario: Print the SSO challenge to stderr
    Given the GitHub API server:
      """
      get('/orgs/acme') {
        response.headers['X-GitHub-SSO'] = 'required; url=http://example.com?auth=HASH'
        status 403
        json({})
      }
      """
    When I run `hub api orgs/acme`
    Then the exit status should be 22
    And the stderr should contain exactly:
      """

      You must authorize your token to access this organization:
      http://example.com?auth=HASH\n
      """
