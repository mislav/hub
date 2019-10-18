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
      my content is here
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
    When I successfully run `hub gist show myhash --json`
    Then the output should contain exactly:
      """
      {"hub_gist1.txt":{"content":"my content is here","raw_url":""},"hub_gist2.txt":{"content":"more content is here","raw_url":""}}
      """

  Scenario: Fetch a gist with many files while specifying a single one
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
      my content is here
      """

  Scenario: Creates a gist
    Given the GitHub API server:
      """
      post('/gists') {
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
    When I successfully run `hub gist create testfile.txt`
    Then the output should contain exactly:
      """
      http://gists.github.com/somehash
      """

  Scenario: Creates a gist with multiple files
    Given the GitHub API server:
      """
      post('/gists') {
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
      http://gists.github.com/somehash
      """

