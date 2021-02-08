Feature: hub pr merge
  Background:
    Given I am in "git://github.com/friederbluemle/hub.git" git repo
    And I am "friederbluemle" on github.com with OAuth token "OTOKEN"

  Scenario: Default merge
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        assert :merge_method => "merge",
               :commit_title => :no,
               :commit_message => :no,
               :sha => :no

        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }
      """
    When I successfully run `hub pr merge 12`
    Then the output should contain exactly ""

  Scenario: Squash merge
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        assert :merge_method => "squash",
               :commit_title => :no,
               :commit_message => :no,
               :sha => :no

        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }
      """
    When I successfully run `hub pr merge --squash 12`
    Then the output should contain exactly ""

  Scenario: Merge with rebase
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        assert :merge_method => "rebase",
               :commit_title => :no,
               :commit_message => :no,
               :sha => :no

        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }
      """
    When I successfully run `hub pr merge --rebase 12`
    Then the output should contain exactly ""

  Scenario: Merge with title
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        assert :commit_title => "mytitle",
               :commit_message => ""

        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }
      """
    When I successfully run `hub pr merge 12 -m mytitle`
    Then the output should contain exactly ""

  Scenario: Merge with title and body
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        assert :commit_title => "mytitle",
               :commit_message => "msg1\n\nmsg2"

        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }
      """
    When I successfully run `hub pr merge 12 -m mytitle -m msg1 -m msg2`
    Then the output should contain exactly ""

  Scenario: Merge with title and body from file
    Given a file named "msg.txt" with:
      """
      mytitle

      msg1

      msg2
      """
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        assert :commit_title => "mytitle",
               :commit_message => "msg1\n\nmsg2"

        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }
      """
    When I successfully run `hub pr merge 12 -F msg.txt`
    Then the output should contain exactly ""

  Scenario: Merge with head SHA
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        assert :sha => "MYSHA"

        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }
      """
    When I successfully run `hub pr merge 12 --head-sha MYSHA`
    Then the output should contain exactly ""

  Scenario: Delete branch
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }

      get('/repos/friederbluemle/hub/pulls/12'){
        json \
          :number => 12,
          :state => "merged",
          :base => {
            :ref => "main",
            :label => "friederbluemle:main",
            :repo => { :owner => { :login => "friederbluemle" } }
          },
          :head => {
            :ref => "patch-1",
            :label => "friederbluemle:patch-1",
            :repo => { :owner => { :login => "friederbluemle" } }
          }
      }

      delete('/repos/friederbluemle/hub/git/refs/heads/patch-1'){
        status 204
      }
      """
    When I successfully run `hub pr merge -d 12`
    Then the output should contain exactly ""

  Scenario: Delete already deleted branch
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }

      get('/repos/friederbluemle/hub/pulls/12'){
        json \
          :number => 12,
          :state => "merged",
          :base => {
            :ref => "main",
            :label => "friederbluemle:main",
            :repo => { :owner => { :login => "friederbluemle" } }
          },
          :head => {
            :ref => "patch-1",
            :label => "friederbluemle:patch-1",
            :repo => { :owner => { :login => "friederbluemle" } }
          }
      }

      delete('/repos/friederbluemle/hub/git/refs/heads/patch-1'){
        status 422
        json :message => "Invalid branch name"
      }
      """
    When I successfully run `hub pr merge -d 12`
    Then the output should contain exactly ""

  Scenario: Delete branch on cross-repo PR
    Given the GitHub API server:
      """
      put('/repos/friederbluemle/hub/pulls/12/merge'){
        json :merged => true,
          :sha => "MERGESHA",
          :message => "All done!"
      }

      get('/repos/friederbluemle/hub/pulls/12'){
        json \
          :number => 12,
          :state => "merged",
          :base => {
            :ref => "main",
            :label => "friederbluemle:main",
            :repo => { :owner => { :login => "friederbluemle" } }
          },
          :head => {
            :ref => "patch-1",
            :label => "monalisa:patch-1",
            :repo => { :owner => { :login => "monalisa" } }
          }
      }
      """
    When I successfully run `hub pr merge -d 12`
    Then the output should contain exactly ""
