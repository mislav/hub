Feature: hub assignees
  Background:
    Given I am in "git://github.com/github/hub.git" git repo
    And I am "cornwe19" on github.com with OAuth token "OTOKEN"

  Scenario: list assignees
    Given the GitHub API server:
    """
    get('/repos/github/hub/assignees') {
      json [
        { :login => "pcorpet",
          :id => 7937848,
          :type => "User",
        },
        { :login => "octokat",
          :id => 999,
          :type => "User",
        },
      ]
    }
    """
    When I successfully run `hub assignees`
    Then the output should contain exactly:
      """
      pcorpet
      octokat

      """
