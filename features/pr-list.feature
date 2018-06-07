Feature: hub pr list
  Background:
    Given I am in "git://github.com/github/hub.git" git repo
    And I am "defunkt" on github.com with OAuth token "OTOKEN"

  Scenario: List pulls
    Given the GitHub API server:
    """
    get('/repos/github/hub/pulls') {
      assert :per_page => "100",
             :page => :no,
             :sort => nil,
             :direction => nil

      response.headers["Link"] = %(<https://api.github.com/repositories/12345?per_page=100&page=2>; rel="next")

      json [
        { :number => 999,
          :title => "First",
          :state => "open",
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-1", :label => "octocat:patch-1" },
          :user => { :login => "octocat" },
        },
        { :number => 102,
          :title => "Second",
          :state => "open",
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-2", :label => "octocat:patch-2" },
          :user => { :login => "octocat" },
        },
        { :number => 13,
          :title => "Third",
          :state => "open",
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-3", :label => "octocat:patch-3" },
          :user => { :login => "octocat" },
        },
      ]
    }

    get('/repositories/12345') {
      assert :per_page => "100",
             :page => "2"

      json [
        { :number => 7,
          :title => "Fourth",
          :state => "open",
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-4", :label => "octocat:patch-4" },
          :user => { :login => "octocat" },
        },
      ]
    }
    """
    When I successfully run `hub pr list`
    Then the output should contain exactly:
      """
          #999  First
          #102  Second
           #13  Third
            #7  Fourth\n
      """

  Scenario: Sort by number of comments ascending
    Given the GitHub API server:
    """
    get('/repos/github/hub/pulls') {
      assert :sort => "comments",
             :direction => "asc"

      json []
    }
    """
    When I successfully run `hub pr list -o comments -^`
    Then the output should contain exactly ""

  Scenario: Filter by base and head
    Given the GitHub API server:
    """
    get('/repos/github/hub/pulls') {
      assert :base => "develop",
             :head => "github:patch-1"

      json []
    }
    """
    When I successfully run `hub pr list -b develop -h patch-1`
    Then the output should contain exactly ""

  Scenario: Filter by head with owner
    Given the GitHub API server:
    """
    get('/repos/github/hub/pulls') {
      assert :head => "mislav:patch-1"

      json []
    }
    """
    When I successfully run `hub pr list -h mislav:patch-1`
    Then the output should contain exactly ""
