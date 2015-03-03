Feature: hub push
  Background:
    Given I am in "git://github.com/mislav/coral.git" git repo

  Scenario: Normal push
    When I successfully run `hub push`
    Then the git command should be unchanged

  Scenario: Push current branch to multiple remotes
    Given I am on the "cool-feature" branch
    When I successfully run `hub push origin,staging`
    Then "git push origin cool-feature" should be run
    Then "git push staging cool-feature" should be run

  Scenario: Push explicit branch to multiple remotes
    When I successfully run `hub push origin,staging,qa cool-feature`
    Then "git push origin cool-feature" should be run
    Then "git push staging cool-feature" should be run
    Then "git push qa cool-feature" should be run

  Scenario: Push multiple refs to multiple remotes
    When I successfully run `hub push origin,staging master new-feature`
    Then "git push origin master new-feature" should be run
    Then "git push staging master new-feature" should be run
