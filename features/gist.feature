Feature: hub gist
  Background:
    Given I am in "git://github.com/octocat/Hello-World.git" git repo
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
    When I successfully run `hub gist myhash`
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
    When I successfully run `hub gist myhash`
    Then the output should contain exactly:
      """
GIST: my gist (myhash)

==== BEGIN hub_gist1.txt ====>
my content is here
<=== END hub_gist1.txt =======
==== BEGIN hub_gist2.txt ====>
more content is here
<=== END hub_gist2.txt =======
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
    When I successfully run `hub gist myhash hub_gist1.txt`
    Then the output should contain exactly:
      """
      my content is here
      """

  Scenario: Fetch a gist with many files without heders
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
    When I successfully run `hub gist myhash --no-headers`
    # can't do "exactly" since the ordering of JSON hashes
    # is not gauranteed, so just check for both lines
    Then the output should contain "my content is here"
    And the stdout should contain "more content is here"

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
    When I successfully run `hub gist --file testfile.txt`
    Then the output should contain exactly:
      """
      http://gists.github.com/somehash
      """

