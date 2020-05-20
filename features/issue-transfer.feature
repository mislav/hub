Feature: hub issue transfer
  Background:
    Given I am in "git://github.com/octocat/hello-world.git" git repo
    And I am "srafi1" on github.com with OAuth token "OTOKEN"

  Scenario: Transfer issue
    Given the GitHub API server:
    """
    count = 0
    post('/graphql') {
      halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN'
      count += 1
      case count
      when 1
        assert :query => /\A\s*query\(/,
          :variables => {
            :issue => 123,
            :sourceOwner => "octocat",
            :sourceRepo => "hello-world",
            :targetOwner => "octocat",
            :targetRepo => "spoon-knife",
          }
        json :data => {
          :source => { :issue => { :id => "ISSUE-ID" } },
          :target => { :id => "REPO-ID" },
        }
      when 2
        assert :query => /\A\s*mutation\(/,
          :variables => {
            :issue => "ISSUE-ID",
            :repo => "REPO-ID",
          }
        json :data => {
          :transferIssue => { :issue => { :url => "the://url" } }
        }
      else
        status 400
        json :message => "request not stubbed"
      end
    }
    """
    When I successfully run `hub issue transfer 123 spoon-knife`
    Then the output should contain exactly "the://url\n"

  Scenario: Transfer to another owner
    Given the GitHub API server:
    """
    count = 0
    post('/graphql') {
      count += 1
      case count
      when 1
        assert :variables => {
          :targetOwner => "monalisa",
          :targetRepo => "playground",
        }
        json :data => {}
      when 2
        json :errors => [
          { :message => "New repository must have the same owner as the current repository" },
        ]
      else
        status 400
        json :message => "request not stubbed"
      end
    }
    """
    When I run `hub issue transfer 123 monalisa/playground`
    Then the exit status should be 1
    Then the stderr should contain exactly:
      """
      API error: New repository must have the same owner as the current repository\n
      """
