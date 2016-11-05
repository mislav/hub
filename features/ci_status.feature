Feature: hub ci-status

  Background:
    Given I am in "git://github.com/michiels/pencilbox.git" git repo
    And I am "michiels" on github.com with OAuth token "OTOKEN"

  Scenario: Fetch commit SHA
    Given there is a commit named "the_sha"
    Given the remote commit state of "michiels/pencilbox" "the_sha" is "success"
    When I run `hub ci-status the_sha`
    Then the output should contain exactly "success\n"
    And the exit status should be 0

  Scenario: Fetch commit SHA with URL
    Given there is a commit named "the_sha"
    Given the remote commit state of "michiels/pencilbox" "the_sha" is "success"
    When I run `hub ci-status the_sha -v`
    Then the output should contain exactly:
      """
      ✔︎	continuous-integration/travis-ci/push	https://travis-ci.org/michiels/pencilbox/builds/1234567\n
      """
    And the exit status should be 0

  Scenario: Multiple statuses with verbose output
    Given there is a commit named "the_sha"
    Given the remote commit states of "michiels/pencilbox" "the_sha" are:
      """
      { :state => "pending",
        :statuses => [
          { :state => "success",
            :context => "continuous-integration/travis-ci/push",
            :target_url => "https://travis-ci.org/michiels/pencilbox/builds/1234567" },
          { :state => "pending",
            :context => "continuous-integration/travis-ci/merge",
            :target_url => nil },
          { :state => "failure",
            :context => "GitHub CLA",
            :target_url => "https://cla.github.com/michiels/pencilbox/accept/mislav" },
          { :state => "error",
            :context => "whatevs!" }
        ]
      }
      """
    When I run `hub ci-status -v the_sha`
    Then the output should contain exactly:
      """
      ✔︎	continuous-integration/travis-ci/push 	https://travis-ci.org/michiels/pencilbox/builds/1234567
      ●	continuous-integration/travis-ci/merge
      ✖︎	GitHub CLA                            	https://cla.github.com/michiels/pencilbox/accept/mislav
      ✖︎	whatevs!\n
      """
    And the exit status should be 2

  Scenario: Exit status 1 for 'error' and 'failure'
    Given the remote commit state of "michiels/pencilbox" "HEAD" is "error"
    When I run `hub ci-status`
    Then the exit status should be 1
    And the output should contain exactly "error\n"

  Scenario: Use HEAD when no sha given
    Given the remote commit state of "michiels/pencilbox" "HEAD" is "pending"
    When I run `hub ci-status`
    Then the exit status should be 2
    And the output should contain exactly "pending\n"

  Scenario: Exit status 3 for no statuses available
    Given there is a commit named "the_sha"
    Given the remote commit state of "michiels/pencilbox" "the_sha" is nil
    When I run `hub ci-status the_sha`
    Then the output should contain exactly "no status\n"
    And the exit status should be 3

  Scenario: Exit status 3 for no statuses available without URL
    Given there is a commit named "the_sha"
    Given the remote commit state of "michiels/pencilbox" "the_sha" is nil
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

  Scenario: Enterprise CI statuses
    Given the "origin" remote has url "git@git.my.org:michiels/pencilbox.git"
    And I am "michiels" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    Given there is a commit named "the_sha"
    Given the remote commit state of "git.my.org/michiels/pencilbox" "the_sha" is "success"
    When I successfully run `hub ci-status the_sha`
    Then the output should contain exactly "success\n"

  Scenario: If alias named ci-status exists, it should not be expanded.
    Given there is a commit named "the_sha"
    Given the remote commit state of "michiels/pencilbox" "the_sha" is "success"
    When I successfully run `git config --global alias.ci-status "ci-status -v"`
    When I run `hub ci-status the_sha`
    Then the output should contain exactly "success\n"
    And the exit status should be 0
