Feature: hub pr show
    Background:
        Given I am in "git://github.com/ashemesh/hub.git" git repo
        And I am "ashemesh" on github.com with OAuth token "OTOKEN"

    Scenario: Open Current Branch Pull Request
        Given I am on the "topic" branch
        Given the GitHub API server:
            """
            get('/repos/ashemesh/hub/pulls'){
                assert  :state => "open",
                        :head => "ashemesh:topic"
                json [
                    {
                        :html_url => "https://github.com/ashemesh/hub/pull/1",
                    },
                ]
            }
            """
        When I successfully run `hub pr show`
        Then "open https://github.com/ashemesh/hub/pull/1" should be run

    Scenario: Open Current Branch Pull Request In Upstream
    Given the "upstream" remote has url "git@github.com:github/hub.git" 
    And I am on the "topic" branch
    Given the GitHub API server:
        """
        get('/repos/github/hub/pulls'){
            assert  :state => "open",
                    :head => "ashemesh:topic"
            json [
                {
                    :html_url => "https://github.com/github/hub/pull/1",
                },
            ]
        }
        """
    When I successfully run `hub pr show`
    Then "open https://github.com/github/hub/pull/1" should be run

    Scenario: Open Pull Request For Given Head
        Given the GitHub API server:
            """
            get('/repos/ashemesh/hub/pulls'){
                assert  :state => "open",
                        :head => "ashemesh:topic"
                json [
                    {
                        :html_url => "https://github.com/ashemesh/hub/pull/1",
                    },
                ]
            }
            """
        When I successfully run `hub pr show --head topic`
        Then "open https://github.com/ashemesh/hub/pull/1" should be run

    Scenario: Open Pull Request In Fork Repository Given Head
        Given the GitHub API server:
            """
            get('/repos/ashemesh/hub'){
                json :html_url => "https://github.com/ashemesh/hub",
                     :parent => { 
                                    :html_url => "https://github.com/github/hub",
                                    :name => "hub",
                                    :owner => { :login => "github" }
                                }
            }
            get('/repos/github/hub/pulls'){
                assert  :state => "open",
                        :head => "ashemesh:topic"
                json [
                    {
                        :html_url => "https://github.com/github/hub/pull/1",
                    },
                ]
            }
            """
        When I successfully run `hub pr show --head ashemesh:topic`
        Then "open https://github.com/github/hub/pull/1" should be run
