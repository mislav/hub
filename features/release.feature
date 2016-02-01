Feature: hub release

  Background:
    Given I am in "git://github.com/mislav/will_paginate.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

  Scenario: List non-draft releases
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: true,
            prerelease: false,
          },
          { tag_name: 'v1.2.0-pre',
            name: 'will_paginate 1.2.0-pre',
            draft: false,
            prerelease: true,
          },
          { tag_name: 'v1.0.2',
            name: 'will_paginate 1.0.2',
            draft: false,
            prerelease: false,
          },
        ]
      }
      """
    When I successfully run `hub release`
    Then the output should contain exactly:
      """
      v1.2.0-pre
      v1.0.2\n
      """

  Scenario: List all releases
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: true,
            prerelease: false,
          },
          { tag_name: 'v1.2.0-pre',
            name: 'will_paginate 1.2.0-pre',
            draft: false,
            prerelease: true,
          },
          { tag_name: 'v1.0.2',
            name: 'will_paginate 1.0.2',
            draft: false,
            prerelease: false,
          },
        ]
      }
      """
    When I successfully run `hub release --include-drafts`
    Then the output should contain exactly:
      """
      v1.2.0
      v1.2.0-pre
      v1.0.2\n
      """

  Scenario: Repository not found when listing releases
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        status 404
        json message: "Not Found",
             documentation_url: "https://developer.github.com/v3"
      }
      """
    When I run `hub release`
    Then the stderr should contain exactly:
      """
      Error fetching releases: Not Found (HTTP 404)
      Not Found\n
      """
    And the exit status should be 1

  Scenario: Server error when listing releases
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        status 504
        '<html><title>Its fine</title></html>'
      }
      """
    When I run `hub release`
    Then the stderr should contain exactly:
      """
      Error fetching releases: invalid character '<' looking for beginning of value (HTTP 504)\n
      """
    And the exit status should be 1

  Scenario: Show specific release
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: true,
            prerelease: false,
            tarball_url: "https://github.com/mislav/will_paginate/archive/v1.2.0.tar.gz",
            zipball_url: "https://github.com/mislav/will_paginate/archive/v1.2.0.zip",
            assets: [
              { browser_download_url: "https://github.com/mislav/will_paginate/releases/download/v1.2.0/example.zip",
              },
            ],
            body: <<MARKDOWN
### Hello to my release

Here is what's broken:
- everything
MARKDOWN
          },
        ]
      }
      """
    When I successfully run `hub release show v1.2.0`
    Then the output should contain exactly:
      """
      will_paginate 1.2.0

      ### Hello to my release

      Here is what's broken:
      - everything\n
      """

  Scenario: Show specific release including downloads
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: true,
            prerelease: false,
            tarball_url: "https://github.com/mislav/will_paginate/archive/v1.2.0.tar.gz",
            zipball_url: "https://github.com/mislav/will_paginate/archive/v1.2.0.zip",
            assets: [
              { browser_download_url: "https://github.com/mislav/will_paginate/releases/download/v1.2.0/example.zip",
              },
            ],
            body: <<MARKDOWN
### Hello to my release

Here is what's broken:
- everything
MARKDOWN
          },
        ]
      }
      """
    When I successfully run `hub release show v1.2.0 --show-downloads`
    Then the output should contain exactly:
      """
      will_paginate 1.2.0

      ### Hello to my release

      Here is what's broken:
      - everything

      ## Downloads

      https://github.com/mislav/will_paginate/releases/download/v1.2.0/example.zip
      https://github.com/mislav/will_paginate/archive/v1.2.0.zip
      https://github.com/mislav/will_paginate/archive/v1.2.0.tar.gz\n
      """

  Scenario: Create a release
    Given the GitHub API server:
      """
      post('/repos/mislav/will_paginate/releases') {
        assert :draft => true,
               :tag_name => "v1.2.0",
               :target_commitish => "master",
               :name => "will_paginate 1.2.0: Instant Gratification Monkey",
               :body => ""

        status 201
        json :html_url => "https://github.com/mislav/will_paginate/releases/v1.2.0"
      }
      """
    When I successfully run `hub release create -dm "will_paginate 1.2.0: Instant Gratification Monkey" v1.2.0`
    Then the output should contain exactly:
      """
      https://github.com/mislav/will_paginate/releases/v1.2.0\n
      """

  Scenario: Create a release with assets
    Given the GitHub API server:
      """
      post('/repos/mislav/will_paginate/releases') {
        status 201
        json :html_url => "https://github.com/mislav/will_paginate/releases/v1.2.0",
             :upload_url => "https://api.github.com/uploads/assets{?name,label}"
      }
      post('/uploads/assets') {
        assert :name => 'hello-1.2.0.tar.gz',
               :label => 'Hello World'
        status 201
      }
      """
    And a file named "hello-1.2.0.tar.gz" with:
      """
      TARBALL
      """
    When I successfully run `hub release create -m "hello" v1.2.0 -a "./hello-1.2.0.tar.gz#Hello World"`
    Then the output should contain exactly:
      """
      https://github.com/mislav/will_paginate/releases/v1.2.0
      Attaching release asset `./hello-1.2.0.tar.gz'...\n
      """

  Scenario: Edit existing release
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { url: 'https://api.github.com/repos/mislav/will_paginate/releases/123',
            tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: true,
            prerelease: false,
            body: <<MARKDOWN
### Hello to my release

Here is what's broken:
- everything
MARKDOWN
          },
        ]
      }
      patch('/repos/mislav/will_paginate/releases/123') {
        assert :name => 'KITTENS EVERYWHERE',
               :draft => false,
               :prerelease => nil
        json({})
      }
      """
    Given the git commit editor is "vim"
    And the text editor adds:
      """
      KITTENS EVERYWHERE
      """
    When I successfully run `hub release edit --draft=false v1.2.0`
    Then there should be no output
