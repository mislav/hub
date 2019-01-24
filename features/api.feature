Feature: hub api
  Background:
    Given I am "octokitten" on github.com with OAuth token "OTOKEN"

  Scenario: GET resource
    Given the GitHub API server:
      """
      get('/hello/world') {
        json :name => "Ed"
      }
      """
    When I successfully run `hub api hello/world`
    Then the output should contain exactly:
      """
      {"name":"Ed"}\n
      """

  Scenario: GET query string
    Given the GitHub API server:
      """
      get('/hello/world') {
        json :name => params[:name]
      }
      """
    When I successfully run `hub api -XGET -F name=Ed hello/world`
    Then the output should contain exactly:
      """
      {"name":"Ed"}\n
      """

  Scenario: GET full URL
    Given the GitHub API server:
      """
      get('/hello/world') {
        json :name => "Faye"
      }
      """
    When I successfully run `hub api https://api.github.com/hello/world`
    Then the output should contain exactly:
      """
      {"name":"Faye"}\n
      """

  Scenario: POST fields
    Given the GitHub API server:
      """
      post('/hello/world') {
        json :name => params[:name],
             :value => params[:a],
             :params => params.size
      }
      """
    When I successfully run `hub api -t -f name=@hubot -F a=b=c hello/world`
    Then the output should contain exactly:
      """
      .name	@hubot
      .value	b=c
      .params	2\n
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
