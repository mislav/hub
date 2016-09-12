Feature: hub sync
  Background:
    Given I am in "dotfiles" git repo
    And I make a commit
    And the "origin" remote has url "git://github.com/lostisland/faraday.git"

  Scenario: Prunes remote branches
    When I successfully run `hub sync`
    Then the output should contain exactly ""
    And "git fetch --prune --quiet --progress origin" should be run

  Scenario: Fast-forwards currently checked out local branch
    Given I am on the "feature" branch pushed to "origin/feature"
    And I successfully run `git reset -q --hard HEAD^`
    When I successfully run `hub sync`
    Then the output should contain "Updated branch feature"
    And "git merge --ff-only --quiet refs/remotes/origin/feature" should be run

  Scenario: Fast-forwards other local branches in the background
    Given I am on the "feature" branch pushed to "origin/feature"
    And I successfully run `git reset -q --hard HEAD^`
    And I am on the "bugfix" branch pushed to "origin/bugfix"
    And I successfully run `git reset -q --hard HEAD^`
    And I successfully run `git checkout -q master`
    When I successfully run `hub sync`
    Then the output should contain "Updated branch feature"
    And the output should contain "Updated branch bugfix"

  Scenario: Refuses to update local branch which has diverged from upstream
    Given I am on the "feature" branch pushed to "origin/feature"
    And I make a commit with message "diverge"
    When I successfully run `hub sync`
    Then the stderr should contain exactly:
      """
      warning: `feature' seems to contain unpushed commits\n
      """

  Scenario: Deletes local branch that had its upstream deleted
    Given I am on the "feature" branch with upstream "origin/feature"
    And I successfully run `git checkout -q master`
    And I successfully run `git merge --no-ff --no-edit feature`
    And I successfully run `git update-ref refs/remotes/origin/master HEAD`
    And I successfully run `rm .git/refs/remotes/origin/feature`
    And I successfully run `git checkout -q feature`
    When I successfully run `hub sync`
    Then the output should contain "Deleted branch feature"

  Scenario: Refuses to delete local branch whose upstream was deleted but not merged to master
    Given I am on the "feature" branch with upstream "origin/feature"
    And I successfully run `rm .git/refs/remotes/origin/feature`
    And I successfully run `git update-ref refs/remotes/origin/master master`
    When I successfully run `hub sync`
    Then the stderr should contain exactly:
      """
      warning: `feature' was deleted on origin, but appears not merged into master\n
      """
