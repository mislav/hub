Feature: hub init

  Scenario: Global template exists
    Given template at "#{ENV['HUB_TEMPLATE']}" exists
    When I run `hub init`
    Then "git init" should be run
    And "#{ENV['HUB_TEMPLATE']}" should be copied here 

  #Scenario: Global template does not exist
    #Given default template does not exist 
    #When I run `hub init`
    #Then "git init" should be run
