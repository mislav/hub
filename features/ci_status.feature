Feature: hub ci-status

  Background:
    Given I am in "git://github.com/krlmlr/R-pkg-template.git" git repo
    And I am "krlmlr" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch commit SHA without URL
    Given there is a commit named "the_sha"
    Given the remote commit state of "krlmlr/R-pkg-template" "the_sha" is "success"
    When I run `hub ci-status the_sha`
    Then the output should contain exactly "success\n"
    And the exit status should be 0

  Scenario: Fetch commit SHA with URL
    Given there is a commit named "the_sha"
    Given the remote commit states of "krlmlr/R-pkg-template" "the_sha" are:
      """
      [ { :state => 'success',
          :target_url => 'https://travis-ci.org/krlmlr/R-pkg-template/builds/12905375' } ]
      """
    When I run `hub ci-status -v the_sha`
    Then the output should contain "success: https://travis-ci.org/krlmlr/R-pkg-template/builds/12905375"
    And the exit status should be 0

  Scenario: Multiple statuses, latest is passing
    Given there is a commit named "the_sha"
    Given the remote commit states of "krlmlr/R-pkg-template" "the_sha" are:
      """
      [ { :state => 'success' },
        { :state => 'pending' }  ]
      """
    When I run `hub ci-status the_sha`
    Then the output should contain exactly "success\n"
    And the exit status should be 0

  Scenario: Exit status 1 for 'error' and 'failure'
    Given the remote commit state of "krlmlr/R-pkg-template" "HEAD" is "error"
    When I run `hub ci-status`
    Then the exit status should be 1
    And the output should contain exactly "error\n"

  Scenario: Use HEAD when no sha given
    Given the remote commit state of "krlmlr/R-pkg-template" "HEAD" is "pending"
    When I run `hub ci-status`
    Then the exit status should be 2
    And the output should contain exactly "pending\n"

  Scenario: Exit status 3 for no statuses available
    Given there is a commit named "the_sha"
    Given the remote commit state of "krlmlr/R-pkg-template" "the_sha" is nil
    When I run `hub ci-status the_sha`
    Then the output should contain exactly "no status\n"
    And the exit status should be 3

  Scenario: Exit status 3 for no statuses available without URL
    Given there is a commit named "the_sha"
    Given the remote commit state of "krlmlr/R-pkg-template" "the_sha" is nil
    When I run `hub ci-status -v the_sha`
    Then the output should contain exactly "no status\n"
    And the exit status should be 3

  Scenario: Abort with message when invalid ref given
    When I run `hub ci-status this-is-an-invalid-ref`
    Then the exit status should be 1
    And the output should contain exactly "Aborted: no revision could be determined from 'this-is-an-invalid-ref'\n"

  Scenario: Non-GitHub repo
    Given the "origin" remote has url "mygh:Manganeez/repo.git"
    When I run `hub ci-status`
    Then the stderr should contain "Aborted: the origin remote doesn't point to a GitHub repository.\n"
    And the exit status should be 1
