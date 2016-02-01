Feature: hub issue
  Background:
    Given I am in "git://github.com/github/hub.git" git repo
    And I am "cornwe19" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch issues
    Given the GitHub API server:
    """
    get('/repos/github/hub/issues') {
      json([
        { :number => 102,
          :title => "First issue",
          :html_url => "https://github.com/github/hub/issues/102",
          :assignee => {
            :login => "octokit"
          }
        },
        { :number => 103,
          :title => "Second issue",
          :html_url => "https://github.com/github/hub/issues/103",
          :assignee => {
            :login => "cornwe19"
          }
        }
      ])
    }
    """
    When I run `hub issue -a Cornwe19`
    Then the output should contain exactly:
      """
          103] Second issue ( https://github.com/github/hub/issues/103 )\n
      """
    And the exit status should be 0
