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
             :direction => "desc"

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

  Scenario: List pull requests with requested reviewers
    Given the GitHub API server:
    """
    get('/repos/github/hub/pulls') {
      assert :per_page => "100",
             :page => :no,
             :sort => nil,
             :direction => "desc"

      json [
        { :number => 999,
          :title => "First",
          :state => "open",
          :base => {
            :ref => "master",
            :label => "github:master",
            :repo => { :owner => { :login => "github" } }
          },
          :head => { :ref => "patch-1", :label => "octocat:patch-1" },
          :user => { :login => "octocat" },
          :requested_reviewers => [
            { :login => "rey" },
          ],
          :requested_teams => [
            { :slug => "troopers" },
            { :slug => "cantina-band" },
          ]
        },
        { :number => 102,
          :title => "Second",
          :state => "open",
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-2", :label => "octocat:patch-2" },
          :user => { :login => "octocat" },
          :requested_reviewers => [
            { :login => "luke" },
            { :login => "jyn" },
          ]
        },
      ]
    }
    """
    When I successfully run `hub pr list -f "%sC%>(8)%i %rs%n"`
    Then the output should contain exactly:
      """
          #999 rey, github/troopers, github/cantina-band
          #102 luke, jyn\n
      """

  @keep-ansi-escape-sequences
  Scenario: List draft status
    Given the GitHub API server:
    """
    get('/repos/github/hub/pulls') {
      halt 400 unless env['HTTP_ACCEPT'] == 'application/vnd.github.shadow-cat-preview+json;charset=utf-8'

      json [
        { :number => 999,
          :state => "open",
          :draft => true,
          :merged_at => nil,
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-2", :label => "octocat:patch-2" },
          :user => { :login => "octocat" },
        },
        { :number => 102,
          :state => "open",
          :draft => false,
          :merged_at => nil,
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-1", :label => "octocat:patch-1" },
          :user => { :login => "octocat" },
        },
        { :number => 42,
          :state => "closed",
          :draft => false,
          :merged_at => "2018-12-11T10:50:33Z",
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-3", :label => "octocat:patch-3" },
          :user => { :login => "octocat" },
        },
        { :number => 8,
          :state => "closed",
          :draft => false,
          :merged_at => nil,
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-4", :label => "octocat:patch-4" },
          :user => { :login => "octocat" },
        },
      ]
    }
    """
    When I successfully run `hub pr list --format "%I %pC %pS %Creset%n" --color`
    Then its output should contain exactly:
      """
      999 \e[37m draft \e[m
      102 \e[32m open \e[m
      42 \e[35m merged \e[m
      8 \e[31m closed \e[m\n
      """
    When I successfully run `hub -c color.ui=always pr list --format "%I %pC %pS %Creset%n"`
    Then its output should contain exactly:
      """
      999 \e[37m draft \e[m
      102 \e[32m open \e[m
      42 \e[35m merged \e[m
      8 \e[31m closed \e[m\n
      """
    When I successfully run `hub -c color.ui=false pr list --format "%I %pC%pS%Creset%n" --color=auto`
    Then its output should contain exactly:
      """
      999 draft
      102 open
      42 merged
      8 closed\n
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

  Scenario: Filter by merged state
    Given the GitHub API server:
    """
    get('/repos/github/hub/pulls') {
      assert :state => "closed"

      json [
        { :number => 999,
          :title => "First",
          :state => "closed",
          :merged_at => "2018-12-11T10:50:33Z",
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-1", :label => "octocat:patch-1" },
          :user => { :login => "octocat" },
        },
        { :number => 102,
          :title => "Second",
          :state => "closed",
          :merged_at => nil,
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-2", :label => "octocat:patch-2" },
          :user => { :login => "octocat" },
        },
        { :number => 13,
          :title => "Third",
          :state => "closed",
          :merged_at => "2018-12-11T10:50:33Z",
          :base => { :ref => "master", :label => "github:master" },
          :head => { :ref => "patch-3", :label => "octocat:patch-3" },
          :user => { :login => "octocat" },
        },
      ]
    }
    """
    When I successfully run `hub pr list --state=merged`
    Then the output should contain exactly:
      """
          #999  First
           #13  Third\n
      """
