Feature: hub pull-request
  Background:
    Given I am in "dotfiles" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: Detached HEAD
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    And I am in detached HEAD
    When I run `hub pull-request`
    Then the stderr should contain "Aborted: not currently on any branch.\n"
    And the exit status should be 1

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

  Scenario: With Unicode characters
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 if request.content_charset != 'utf-8'
        halt 422 if params[:title] != 'ăéñøü'
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request ăéñøü`
    Then the output should contain exactly "the://url\n"

  Scenario: Non-existing base
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    Given the GitHub API server:
      """
      post('/repos/origin/coral/pulls') { 404 }
      """
    When I run `hub pull-request -b origin:master hereyougo`
    Then the exit status should be 1
    Then the stderr should contain:
      """
      Error creating pull request: Not Found (HTTP 404)
      Are you sure that github.com/origin/coral exists?
      """

  Scenario: Supplies User-Agent string to API calls
    Given the "origin" remote has url "git://github.com/mislav/coral.git"
    Given the GitHub API server:
      """
      post('/repos/mislav/coral/pulls') {
        halt 400 unless request.user_agent.include?('Hub')
        json :html_url => "the://url"
      }
      """
    When I successfully run `hub pull-request useragent`
    Then the output should contain exactly "the://url\n"
