Feature: hub last-status
  Background:
    Given I am in "pencilbox" git repo
    And I am "michiels" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch commit SHA
    Given the "origin" remote has url "git://github.com/michiels/pencilbox.git"
    Given the GitHub API server:
      """
      get('/repos/michiels/pencilbox/statuses/the_sha') {
        json [ { :state => "success" } ]
      }
      """
    When I successfully run `hub last-status the_sha`
    Then the output should contain exactly "success\n"