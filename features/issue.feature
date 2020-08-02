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
             :direction => "desc"

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

  Scenario: List limited number of issues
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      response.headers["Link"] = %(<https://api.github.com/repositories/12345/issues?per_page=100&page=2>; rel="next")
      assert :per_page => "3"
      json [
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
        { :number => 999,
          :title => "Third issue",
          :state => "open",
          :user => { :login => "octocat" },
        },
      ]
    }
    """
    When I successfully run `hub issue -L 2`
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
             :direction => "desc"

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

  Scenario: Fetch issues assigned to milestone by number
    Given the GitHub API server:
      """
      get('/repos/github/hub/issues') {
        assert :milestone => "12"
        json []
      }
      """
    When I successfully run `hub issue -M 12`

  Scenario: Fetch issues assigned to milestone by name
    Given the GitHub API server:
      """
      get('/repos/github/hub/milestones') {
        status 200
        json [
          { :number => 237, :title => "prerelease" },
          { :number => 1337, :title => "v1" },
          { :number => 41319, :title => "Hello World!" }
        ]
      }
      get('/repos/github/hub/issues') {
        assert :milestone => "1337"
        json []
      }
      """
    When I successfully run `hub issue -M v1`

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

  Scenario: Custom format with no-color labels
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      json [
        { :number => 102,
          :title => "First issue",
          :state => "open",
          :user => { :login => "morganwahl" },
          :labels => [
            { :name => 'Has Migration',
              :color => 'cfcfcf' },
            { :name => 'Maintenance Window',
              :color => '888888' },
          ]
        },
        { :number => 201,
          :title => "No labels",
          :state => "open",
          :user => { :login => "octocat" },
          :labels => []
        },
      ]
    }
    """
    When I successfully run `hub issue -f "%I: %L%n" --color=never`
    Then the output should contain exactly:
      """
      102: Has Migration, Maintenance Window
      201: \n
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
    When I successfully run `hub issue create -m "hello" -M 12 --assign mislav,josh -apcorpet`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """

  Scenario: Create an issue with milestone by name
    Given the GitHub API server:
      """
      get('/repos/github/hub/milestones') {
        status 200
        json [
          { :number => 237, :title => "prerelease" },
          { :number => 1337, :title => "v1" },
          { :number => 41319, :title => "Hello World!" }
        ]
      }
      post('/repos/github/hub/issues') {
        assert :milestone => 41319
        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create -m "hello" -M "hello world!"`
    Then the output should contain exactly:
      """
      https://github.com/github/hub/issues/1337\n
      """

  Scenario: Editing empty issue message
    Given the git commit editor is "vim"
    And the text editor adds:
      """
      hello

      my nice issue
      """
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "hello",
               :body => "my nice issue"

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create -m '' --edit`
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

  Scenario: Multiple issue templates
    Given the git commit editor is "vim"
    And the text editor adds:
      """
      hello
      """
    And a file named ".github/ISSUE_TEMPLATE/bug_report.md" with:
      """
      I want to report a bug
      """
    And a file named ".github/ISSUE_TEMPLATE/feature_request.md" with:
      """
      There is a feature that I need!
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

  Scenario: Multiple issue templates with default
    Given the git commit editor is "vim"
    And the text editor adds:
      """
      hello
      """
    And a directory named ".github/ISSUE_TEMPLATE"
    And a file named ".github/ISSUE_TEMPLATE.md" with:
      """
      The default template
      """
    Given the GitHub API server:
      """
      post('/repos/github/hub/issues') {
        assert :title => "hello",
               :body => "The default template"

        status 201
        json :html_url => "https://github.com/github/hub/issues/1337"
      }
      """
    When I successfully run `hub issue create`
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

  Scenario: Update an issue's title
    Given the GitHub API server:
      """
      patch('/repos/github/hub/issues/1337') {
        assert :title => "Not workie, pls fix",
               :body => "",
               :milestone => :no,
               :assignees => :no,
               :labels => :no,
               :state => :no
      }
      """
    Then I successfully run `hub issue update 1337 -m "Not workie, pls fix"`

  Scenario: Update an issue's state
    Given the GitHub API server:
      """
      patch('/repos/github/hub/issues/1337') {
        assert :title => :no,
               :labels => :no,
               :state => "closed"
      }
      """
    Then I successfully run `hub issue update 1337 -s closed`
    
  Scenario: Update an issue's labels
    Given the GitHub API server:
      """
      patch('/repos/github/hub/issues/1337') {
        assert :title => :no,
               :body => :no,
               :milestone => :no,
               :assignees => :no,
               :labels => ["bug", "important"]
      }
      """
    Then I successfully run `hub issue update 1337 -l bug,important`

  Scenario: Update an issue's milestone
    Given the GitHub API server:
      """
      patch('/repos/github/hub/issues/1337') {
        assert :title => :no,
               :body => :no,
               :milestone => 42,
               :assignees => :no,
               :labels => :no
      }
      """
    Then I successfully run `hub issue update 1337 -M 42`

  Scenario: Update an issue's milestone by name
    Given the GitHub API server:
      """
      get('/repos/github/hub/milestones') {
        status 200
        json [
          { :number => 237, :title => "prerelease" },
          { :number => 42, :title => "Hello World!" }
        ]
      }
      patch('/repos/github/hub/issues/1337') {
        assert :title => :no,
               :body => :no,
               :milestone => 42,
               :assignees => :no,
               :labels => :no
      }
      """
    Then I successfully run `hub issue update 1337 -M "hello world!"`

  Scenario: Update an issue's assignees
    Given the GitHub API server:
      """
      patch('/repos/github/hub/issues/1337') {
        assert :title => :no,
               :body => :no,
               :milestone => :no,
               :assignees => ["Cornwe19"],
               :labels => :no
      }
      """
    Then I successfully run `hub issue update 1337 -a Cornwe19`

  Scenario: Update an issue's title, labels, milestone, and assignees
    Given the GitHub API server:
      """
      patch('/repos/github/hub/issues/1337') {
        assert :title => "Not workie, pls fix",
               :body => "",
               :milestone => 42,
               :assignees => ["Cornwe19"],
               :labels => ["bug", "important"]
      }
      """
    Then I successfully run `hub issue update 1337  -m "Not workie, pls fix" -M 42 -l bug,important -a Cornwe19`

  Scenario: Clear existing issue labels, assignees, milestone
    Given the GitHub API server:
      """
      patch('/repos/github/hub/issues/1337') {
        assert :title => :no,
               :body => :no,
               :milestone => nil,
               :assignees => [],
               :labels => []
      }
      """
    Then I successfully run `hub issue update 1337 --milestone= --assign= --labels=`

  Scenario: Update an issue's title and body manually
    Given the git commit editor is "vim"
    And the text editor adds:
      """
      My new title
      """
    Given the GitHub API server:
      """
      get('/repos/github/hub/issues/1337') {
        json \
          :number => 1337,
          :title => "My old title",
          :body => "My old body"
      }
      patch('/repos/github/hub/issues/1337') {
        assert :title => "My new title",
               :body => "My old title\n\nMy old body",
               :milestone => :no,
               :assignees => :no,
               :labels => :no
      }
      """
    Then I successfully run `hub issue update 1337 --edit`

  Scenario: Update an issue's title and body via a file
    Given a file named "my-issue.md" with:
      """
      My new title
      
      My new body
      """
    Given the GitHub API server:
      """
      patch('/repos/github/hub/issues/1337') {
        assert :title => "My new title",
               :body => "My new body",
               :milestone => :no,
               :assignees => :no,
               :labels => :no
      }
      """
    Then I successfully run `hub issue update 1337 -F my-issue.md`

  Scenario: Update an issue without specifying fields to update
    When I run `hub issue update 1337`
    Then the exit status should be 1
    Then the stderr should contain "please specify fields to update"
    Then the stderr should contain "Usage: hub issue"

  Scenario: Fetch issue labels
    Given the GitHub API server:
    """
    get('/repos/github/hub/labels') {
      response.headers["Link"] = %(<https://api.github.com/repositories/12345/labels?per_page=100&page=2>; rel="next")
      assert :per_page => "100", :page => nil
      json [
        { :name => "Discuss",
          :color => "0000ff",
        },
        { :name => "bug",
          :color => "ff0000",
        },
        { :name => "feature",
          :color => "00ff00",
        },
      ]
    }
    get('/repositories/12345/labels') {
      assert :per_page => "100", :page => "2"
      json [
        { :name => "affects",
          :color => "ffffff",
        },
      ]
    }
    """
    When I successfully run `hub issue labels`
    Then the output should contain exactly:
      """
      affects
      bug
      Discuss
      feature\n
      """

  Scenario: Fetch single issue
    Given the GitHub API server:
      """
      get('/repos/github/hub/issues/102') {
        json \
          :number => 102,
          :state => "open",
          :body => "I want this feature",
          :title => "Feature request for hub issue show",
          :created_at => "2017-04-14T16:00:49Z",
          :user => { :login => "royels" },
          :assignees => [{:login => "royels"}],
          :comments => 1
      }
      get('/repos/github/hub/issues/102/comments') {
        json [
          { :body => "I am from the future",
            :created_at => "2011-04-14T16:00:49Z",
            :user => { :login => "octocat" }
          },
          { :body => "I did the thing",
            :created_at => "2013-10-30T22:20:00Z",
            :user => { :login => "hubot" }
          },
        ]
      }
      """
    When I successfully run `hub issue show 102`
    Then the output should contain exactly:
      """
      # Feature request for hub issue show

      * created by @royels on 2017-04-14 16:00:49 +0000 UTC
      * assignees: royels

      I want this feature

      ## Comments:

      ### comment by @octocat on 2011-04-14 16:00:49 +0000 UTC

      I am from the future

      ### comment by @hubot on 2013-10-30 22:20:00 +0000 UTC

      I did the thing\n
      """

  Scenario: Format single issue
    Given the GitHub API server:
      """
      get('/repos/github/hub/issues/102') {
        json \
          :number => 102,
          :state => "open",
          :body => "I want this feature",
          :title => "Feature request for hub issue show",
          :created_at => "2017-04-14T16:00:49Z",
          :user => { :login => "royels" },
          :assignees => [{:login => "royels"}],
          :comments => 1
      }
      get('/repos/github/hub/issues/102/comments') {
        json [
          { :body => "I am from the future",
            :created_at => "2011-04-14T16:00:49Z",
            :user => { :login => "octocat" }
          },
          { :body => "I did the thing",
            :created_at => "2013-10-30T22:20:00Z",
            :user => { :login => "hubot" }
          },
        ]
      }
      """
    When I successfully run `hub issue show 102 --format='%I %t%n%n%b%n'`
    Then the output should contain exactly:
      """
      102 Feature request for hub issue show

      I want this feature\n
      """

  Scenario: Format with literal % characters
    Given the GitHub API server:
      """
      get('/repos/github/hub/issues/102') {
        json \
          :number => 102,
          :state => "open",
          :title => "Feature request % hub",
          :user => { :login => "alexfornuto" }
      }
      get('/repos/github/hub/issues/102/comments') {
        json []
      }
      """
    When I successfully run `hub issue show 102 --format='%t%%t%%n%n'`
    Then the output should contain exactly:
      """
      Feature request % hub%t%n\n
      """

  Scenario: Did not supply an issue number
    When I run `hub issue show`
    Then the exit status should be 1
    Then the stderr should contain "Usage: hub issue"

  Scenario: Show error message if http code is not 200 for issues endpoint
    Given the GitHub API server:
      """
      get('/repos/github/hub/issues/102') {
        status 500
      }
      """
    When I run `hub issue show 102`
    Then the output should contain exactly:
      """
      Error fetching issue: Internal Server Error (HTTP 500)\n
      """

  Scenario: Show error message if http code is not 200 for comments endpoint
    Given the GitHub API server:
      """
      get('/repos/github/hub/issues/102') {
        json \
          :number => 102,
          :body => "I want this feature",
          :title => "Feature request for hub issue show",
          :created_at => "2017-04-14T16:00:49Z",
          :user => { :login => "royels" }
      }
      get('/repos/github/hub/issues/102/comments') {
        status 404
      }
      """
    When I run `hub issue show 102`
    Then the output should contain exactly:
      """
      Error fetching comments for issue: Not Found (HTTP 404)\n
      """
