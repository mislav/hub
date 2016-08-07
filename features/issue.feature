Feature: hub issue
  Background:
    Given I am in "git://github.com/github/hub.git" git repo
    And I am "cornwe19" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch issues
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :assignee => "Cornwe19"

      json [
        { :number => 102,
          :title => "First issue",
          :state => "open",
        },
        { :number => 13,
          :title => "Second issue",
          :state => "open",
        },
      ]
    }
    """
    When I run `hub issue -a Cornwe19`
    Then the output should contain exactly:
      """
          #102  First issue
           #13  Second issue\n
      """
    And the exit status should be 0

  Scenario: Custom format for issues list
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      json [
        { :number => 102,
          :title => "First issue",
          :state => "open",
          :user => { :login => "lascap" },
        },
        { :number => 13,
          :title => "Second issue",
          :state => "closed",
          :user => { :login => "mislav" },
        },
      ]
    }
    """
    When I run `hub issue -f "%in,%u%n" -a Cornwe19`
    Then the output should contain exactly:
      """
      102,lascap
      13,mislav\n
      """
    And the exit status should be 0

  Scenario: Create an issue
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "Not workie, pls fix",
               :body => "",
               :labels => nil

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create -m "Not workie, pls fix"`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """

  Scenario: Create an issue with labels
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "hello",
               :body => "",
               :labels => ["wont fix", "docs"]

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create -m "hello" -l "wont fix,docs"`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """
