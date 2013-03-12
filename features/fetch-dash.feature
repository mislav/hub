Feature: hub fetch
  Background:
    Given I am in "mrjob" git repo
    And the "origin" remote has url "git://github.com/Yelp/mrjob.git"
    And I am "mislav" on github.com with OAuth token "OTOKEN"
    And the GitHub API server:
      """
      get('/repos/mislav/mrjob') { json :private => false }
      """

  Scenario: Fetches when remote has dash in username
    When I successfully run `hub fetch ankit-maverick`
    Then "git fetch ankit-maverick" should be run
    And there should be no output
    And the url for "ankit-maverick" should be "git://github.com/ankit-maverick/mrjob.git"
