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

  Scenario: List non-prerelease releases
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
    When I successfully run `hub release --exclude-prereleases`
    Then the output should contain exactly:
      """
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

  Scenario: Fetch releases across multiple pages
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        assert :per_page => "100", :page => :no
        response.headers["Link"] = %(<https://api.github.com/repositories/12345?per_page=100&page=2>; rel="next")
        json [
          { tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: false,
            prerelease: false,
          },
        ]
      }

      get('/repositories/12345') {
        assert :per_page => "100"
        if params[:page] == "2"
          response.headers["Link"] = %(<https://api.github.com/repositories/12345?per_page=100&page=3>; rel="next")
          json [
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
        elsif params[:page] == "3"
          json [
            { tag_name: 'v1.0.0',
              name: 'will_paginate 1.0.0',
              draft: false,
              prerelease: true,
            },
          ]
        else
          status 400
        end
      }
      """
      When I successfully run `hub release`
      Then the output should contain exactly:
      """
      v1.2.0
      v1.2.0-pre
      v1.0.2
      v1.0.0\n
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
               :target_commitish => "",
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

  Scenario: Create a release with target commitish
    Given the GitHub API server:
      """
      post('/repos/mislav/will_paginate/releases') {
        assert :tag_name => "v1.2.0",
               :target_commitish => "my-branch"

        status 201
        json :html_url => "https://github.com/mislav/will_paginate/releases/v1.2.0"
      }
      """
    When I successfully run `hub release create -m hello v1.2.0 -t my-branch`
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

  Scenario: Open new release in web browser
    Given the GitHub API server:
      """
      post('/repos/mislav/will_paginate/releases') {
        status 201
        json :html_url => "https://github.com/mislav/will_paginate/releases/v1.2.0"
      }
      """
    When I successfully run `hub release create -o -m hello v1.2.0`
    Then the output should contain exactly ""
    And "open https://github.com/mislav/will_paginate/releases/v1.2.0" should be run

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

  Scenario: Edit existing release by uploading assets
    Given the GitHub API server:
      """
      deleted = false
      get('/repos/mislav/will_paginate/releases') {
        json [
          { url: 'https://api.github.com/repos/mislav/will_paginate/releases/123',
            upload_url: 'https://api.github.com/uploads/assets{?name,label}',
            tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: true,
            prerelease: false,
            assets: [
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/456',
                name: 'hello-1.2.0.tar.gz',
              },
            ],
          },
        ]
      }
      delete('/repos/mislav/will_paginate/assets/456') {
        deleted = true
        status 204
      }
      post('/uploads/assets') {
        halt 422 unless deleted
        assert :name => 'hello-1.2.0.tar.gz',
               :label => nil
        status 201
      }
      """
    And a file named "hello-1.2.0.tar.gz" with:
      """
      TARBALL
      """
    When I successfully run `hub release edit -m "" v1.2.0 -a hello-1.2.0.tar.gz`
    Then the output should contain exactly:
      """
      Attaching release asset `hello-1.2.0.tar.gz'...\n
      """

    Scenario: Download a release asset.
      Given the GitHub API server:
        """
        get('/repos/mislav/will_paginate/releases') {
          json [
            { url: 'https://api.github.com/repos/mislav/will_paginate/releases/123',
              upload_url: 'https://api.github.com/uploads/assets{?name,label}',
              tag_name: 'v1.2.0',
              name: 'will_paginate 1.2.0',
              draft: true,
              prerelease: false,
              assets: [
                { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9876',
                  name: 'hello-1.2.0.tar.gz',
                },
              ],
            },
          ]
        }
        get('/repos/mislav/will_paginate/assets/9876') {
          halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN'
          halt 415 unless request.accept?('application/octet-stream')
          status 302
          headers['Location'] = 'https://github-cloud.s3.amazonaws.com/releases/12204602/22ea221a-cf2f-11e2-222a-b3a3c3b3aa3a.gz'
          ""
        }
        get('/releases/12204602/22ea221a-cf2f-11e2-222a-b3a3c3b3aa3a.gz', :host_name => 'github-cloud.s3.amazonaws.com') {
          halt 400 unless request.env['HTTP_AUTHORIZATION'].nil?
          halt 415 unless request.accept?('application/octet-stream')
          headers['Content-Type'] = 'application/octet-stream'
          "ASSET_TARBALL"
        }
        """
        When I successfully run `hub release download v1.2.0`
        Then the output should contain exactly:
          """
          Downloading hello-1.2.0.tar.gz ...\n
          """
        And the file "hello-1.2.0.tar.gz" should contain exactly:
          """
          ASSET_TARBALL
          """

  Scenario: Delete a release
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
          json [
            { url: 'https://api.github.com/repos/mislav/will_paginate/releases/123',
              tag_name: 'v1.2.0',
            },
          ]
      }

      delete('/repos/mislav/will_paginate/releases/123') {
        status 204
      }
      """
    When I successfully run `hub release delete v1.2.0`
    Then there should be no output

  Scenario: Enterprise list releases
    Given the "origin" remote has url "git@git.my.org:mislav/will_paginate.git"
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    Given the GitHub API server:
      """
      get('/api/v3/repos/mislav/will_paginate/releases', :host_name => 'git.my.org') {
        json [
          { tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: false,
            prerelease: false,
          },
        ]
      }
      """
    When I successfully run `hub release`
    Then the output should contain exactly:
      """
      v1.2.0\n
      """
