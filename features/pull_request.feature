Feature: hub pull-request
  Background:
    Given I am in "dotfiles" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Non-GitHub repo
    Given the "origin" remote has url "mygh:Manganeez/repo.git"
    When I run `hub pull-request`
    Then the stderr should contain "Aborted: the origin remote doesn't point to a GitHub repository.\n"
    And the exit status should be 1

  Scenario: Create pull request respecting "insteadOf" configuration
    Given the "origin" remote has url "mygh:Manganeez/repo.git"
    When I successfully run `git config url."git@github.com:".insteadOf mygh:`
    Given the GitHub API server:
      """
      post('/repos/Manganeez/repo/pulls') {
        { :base  => 'master',
          :head  => 'mislav:master',
          :title => 'hereyougo'
        }.each do |param, value|
          if params[param] != value
            halt 422, json(
              :message => "expected %s to be %s; got %s" % [
                param.inspect,
                value.inspect,
                params[param].inspect
              ]
            )
          end
        end
        json :html_url => "https://github.com/Manganeez/repo/pull/12"
      }
      """
    When I successfully run `hub pull-request hereyougo`
    Then the output should contain exactly "https://github.com/Manganeez/repo/pull/12\n"
