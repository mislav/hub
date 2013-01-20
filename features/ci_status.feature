Feature: hub ci-status

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
    When I run `hub ci-status the_sha`
    Then the output should contain exactly "success\n"
    And the exit status should be 0

  Scenario: Exit status 1 for 'error' and 'failure'
    Given the "origin" remote has url "git://github.com/michiels/pencilbox.git"
    Given a HEAD commit with GitHub status "error"
    When I run `hub ci-status`
    Then the exit status should be 1

  Scenario: Use HEAD when no sha given
    Given the "origin" remote has url "git://github.com/michiels/pencilbox.git"
    Given a HEAD commit with GitHub status "pending"
    When I run `hub ci-status`
    Then the exit status should be 2

  Scenario: Exit status 3 for no statuses available
    Given the "origin" remote has url "git://github.com/michiels/pencilbox.git"
    Given the GitHub API server:
      """
      get('/repos/michiels/pencilbox/statuses/the_sha') {
        json [ ]
      }
      """
    When I run `hub ci-status the_sha`
    Then the output should contain exactly "no status\n"
    And the exit status should be 3

  Scenario: Non-GitHub repo
    Given the "origin" remote has url "mygh:Manganeez/repo.git"
    When I run `hub ci-status`
    Then the stderr should contain "Aborted: the origin remote doesn't point to a GitHub repository.\n"
    And the exit status should be 1