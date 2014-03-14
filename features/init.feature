Feature: hub init

  Scenario: Init with default template
    Given default template exists
    When I run `hub init`
    Then "git init" should be run
    And a file named "README.md" should exist

  Scenario: Init without default template
    Given default template exists
    When I run `hub init -c`
    Then "git init" should be run
    And a file named "README.md" should not exist

