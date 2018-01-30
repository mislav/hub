Feature: hub labels
  Background:
    Given I am in "git://github.com/github/hub.git" git repo
    And I am "testeroni" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch issues
    Given the GitHub API server:
    """
    get('/repos/github/hub/labels') {
      json [
        { :name => "bug",
          :color => "ff0000",
        },
        { :name => "feature",
          :color => "00ff00",
        },
      ]
    }
    """
    When I successfully run `hub issue labels`
    Then the output should contain exactly:
      """
      bug
      feature\n
      """
