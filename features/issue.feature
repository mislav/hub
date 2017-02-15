Feature: hub issue
  Background:
    Given I am in "git://github.com/github/hub.git" git repo
    And I am "cornwe19" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch issues
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :assignee => "Cornwe19",
             :sort => nil,
             :direction => nil

      json [
        { :number => 999,
          :title => "First pull",
          :state => "open",
          :user => { :login => "octocat" },
          :pull_request => { },
        },
        { :number => 102,
          :title => "First issue",
          :state => "open",
          :user => { :login => "octocat" },
        },
        { :number => 13,
          :title => "Second issue",
          :state => "open",
          :user => { :login => "octocat" },
        },
      ]
    }
    """
    When I successfully run `hub issue -a Cornwe19`
    Then the output should contain exactly:
      """
          #102  First issue
           #13  Second issue\n
      """

  Scenario: Fetch issues and pull requests
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :assignee => "Cornwe19",
             :sort => nil,
             :direction => nil

      json [
        { :number => 999,
          :title => "First pull",
          :state => "open",
          :user => { :login => "octocat" },
          :pull_request => { },
        },
        { :number => 102,
          :title => "First issue",
          :state => "open",
          :user => { :login => "octocat" },
        },
        { :number => 13,
          :title => "Second issue",
          :state => "open",
          :user => { :login => "octocat" },
        },
      ]
    }
    """
    When I successfully run `hub issue -a Cornwe19 --include-pulls`
    Then the output should contain exactly:
      """
          #999  First pull
          #102  First issue
           #13  Second issue\n
      """

  Scenario: Fetch issues not assigned to any milestone
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :milestone => "none"
      json []
    }
    """
    When I successfully run `hub issue -M none`

  Scenario: Fetch issues created by a given user
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :creator => "octocat"
      json []
    }
    """
    When I successfully run `hub issue -c octocat`

  Scenario: Fetch issues mentioning a given user
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :mentioned => "octocat"
      json []
    }
    """
    When I successfully run `hub issue -@ octocat`

  Scenario: Fetch issues with certain labels
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :labels => "foo,bar"
      json []
    }
    """
    When I successfully run `hub issue -l foo,bar`

  Scenario: Fetch issues updated after a certain date and time
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :since => "2016-08-18T09:11:32Z"
      json []
    }
    """
    When I successfully run `hub issue -d 2016-08-18T09:11:32Z`

  Scenario: Fetch issues sorted by number of comments ascending
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :sort => "comments"
      assert :direction => "asc"

      json []
    }
    """
    When I successfully run `hub issue -o comments -^`

  Scenario: Fetch issues across multiple pages
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :per_page => "100", :page => :no
      response.headers["Link"] = %(<https://api.github.com/repositories/12345?per_page=100&page=2>; rel="next")
      json [
        { :number => 102,
          :title => "First issue",
          :state => "open",
          :user => { :login => "octocat" },
        },
      ]
    }

    get('/repositories/12345') {
      assert :per_page => "100"
      if params[:page] == "2"
        response.headers["Link"] = %(<https://api.github.com/repositories/12345?per_page=100&page=3>; rel="next")
        json [
          { :number => 13,
            :title => "Second issue",
            :state => "open",
            :user => { :login => "octocat" },
          },
          { :number => 103,
            :title => "Issue from 2nd page",
            :state => "open",
            :user => { :login => "octocat" },
          },
        ]
      elsif params[:page] == "3"
        json [
          { :number => 21,
            :title => "Even more issuez",
            :state => "open",
            :user => { :login => "octocat" },
          },
        ]
      else
        status 400
      end
    }
    """
    When I successfully run `hub issue`
    Then the output should contain exactly:
      """
          #102  First issue
           #13  Second issue
          #103  Issue from 2nd page
           #21  Even more issuez\n
      """

  Scenario: Custom format for issues list
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      assert :assignee => 'Cornwe19'
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
    When I successfully run `hub issue -f "%I,%au%n" -a Cornwe19`
    Then the output should contain exactly:
      """
      102,lascap
      13,mislav\n
      """

  Scenario: List all assignees
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      json [
        { :number => 102,
          :title => "First issue",
          :state => "open",
          :user => { :login => "octocat" },
          :assignees => [
            { :login => "mislav" },
            { :login => "lascap" },
          ]
        },
        { :number => 13,
          :title => "Second issue",
          :state => "closed",
          :user => { :login => "octocat" },
          :assignees => [
            { :login => "keenahn" },
          ]
        },
      ]
    }
    """
    When I successfully run `hub issue -f "%I:%as%n"`
    Then the output should contain exactly:
      """
      102:mislav, lascap
      13:keenahn\n
      """

  Scenario: Create an issue
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "Not workie, pls fix",
               :body => "",
               :labels => :no

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create -m "Not workie, pls fix"`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """

  Scenario: Create an issue and open in browser
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        status 201
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub issue create -o -m hello`
    Then the output should contain exactly ""
    Then "open the://url" should be run

  Scenario: Create an issue with labels
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "hello",
               :body => "",
               :milestone => :no,
               :assignees => :no,
               :labels => ["wont fix", "docs", "nope"]

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create -m "hello" -l "wont fix,docs" -lnope`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """

  Scenario: Create an issue with milestone and assignees
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "hello",
               :body => "",
               :milestone => 12,
               :assignees => ["mislav", "josh", "pcorpet"],
               :labels => :no

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create -m "hello" -M 12 -a mislav,josh -apcorpet`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """

  Scenario: Issue template
    Given the git commit editor is "vim"
    And the text editor adds:
      """
      hello
      """
    And a file named "issue_template.md" with:
      """
      my nice issue template
      """
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "hello",
               :body => "my nice issue template"

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """

  Scenario: Issue template from a subdirectory
    Given the git commit editor is "vim"
    And the text editor adds:
      """
      hello
      """
    And a file named ".github/issue_template.md" with:
      """
      my nice issue template
      """
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "hello",
               :body => "my nice issue template"

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    Given a directory named "subdir"
    When I cd to "subdir"
    And I successfully run `hub issue create`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """

  Scenario: A file named ".github"
    Given the git commit editor is "vim"
    And the text editor adds:
      """
      hello
      """
    And a file named ".github" with:
      """
      this is ignored
      """
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "hello",
               :body => ""

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """
