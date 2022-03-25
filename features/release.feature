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

  Scenario: List limited number of releases
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        response.headers["Link"] = %(<https://api.github.com/repositories/12345?per_page=100&page=2>; rel="next")
        assert :per_page => "3"
        json [
          { tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: false,
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
    When I successfully run `hub release -L 2`
    Then the output should contain exactly:
      """
      v1.2.0
      v1.2.0-pre\n
      """

  Scenario: Pretty-print releases
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
            draft: true,
            prerelease: false,
            created_at: '2018-02-27T19:35:32Z',
            published_at: '2018-04-01T19:35:32Z',
            assets: [
              {browser_download_url: 'the://url', label: ''},
            ],
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
    When I successfully run `hub release --include-drafts --format='%t (%S)%n'`
    Then the output should contain exactly:
      """
      will_paginate 1.2.0 (draft)
      will_paginate 1.2.0-pre (pre-release)
      will_paginate 1.0.2 ()\n
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

  Scenario: Format specific release
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
    When I successfully run `hub release show v1.2.0 --format='%t (%T)%n%as%n%n%b%n'`
    Then the output should contain exactly:
      """
      will_paginate 1.2.0 (v1.2.0)
      https://github.com/mislav/will_paginate/releases/download/v1.2.0/example.zip	

      ### Hello to my release

      Here is what's broken:
      - everything\n\n
      """

  Scenario: Show release no tag
    When I run `hub release show`
    Then the exit status should be 1
    Then the stderr should contain "hub release show"

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

  Scenario: Create a release from file
    Given the GitHub API server:
      """
      post('/repos/mislav/will_paginate/releases') {
        assert :name => "Epic New Version",
               :body => "body\ngoes\n\nhere"

        status 201
        json :html_url => "https://github.com/mislav/will_paginate/releases/v1.2.0"
      }
      """
    And a file named "message.txt" with:
      """
      Epic New Version

      body
      goes

      here
      """
    When I successfully run `hub release create -F message.txt v1.2.0`
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
             :upload_url => "https://uploads.github.com/uploads/assets{?name,label}"
      }
      post('/uploads/assets', :host_name => 'uploads.github.com') {
        halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN'
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
      Attaching 1 asset...\n
      """

  Scenario: Retry attaching assets on 5xx errors
    Given the GitHub API server:
      """
      attempt = 0
      post('/repos/mislav/will_paginate/releases') {
        status 201
        json :html_url => "https://github.com/mislav/will_paginate/releases/v1.2.0",
             :upload_url => "https://uploads.github.com/uploads/assets{?name,label}"
      }
      post('/uploads/assets', :host_name => 'uploads.github.com') {
        attempt += 1
        halt 400 unless request.body.read.to_s == "TARBALL"
        halt 502 if attempt == 1
        status 201
      }
      """
    And a file named "hello-1.2.0.tar.gz" with:
      """
      TARBALL
      """
    When I successfully run `hub release create -m "hello" v1.2.0 -a hello-1.2.0.tar.gz`
    Then the output should contain exactly:
      """
      https://github.com/mislav/will_paginate/releases/v1.2.0
      Attaching 1 asset...\n
      """

  Scenario: Create a release with some assets failing
    Given the GitHub API server:
      """
      post('/repos/mislav/will_paginate/releases') {
        status 201
        json :tag_name => "v1.2.0",
             :html_url => "https://github.com/mislav/will_paginate/releases/v1.2.0",
             :upload_url => "https://uploads.github.com/uploads/assets{?name,label}"
      }
      post('/uploads/assets', :host_name => 'uploads.github.com') {
        halt 422 if params[:name] == "two"
        status 201
      }
      """
    And a file named "one" with:
      """
      ONE
      """
    And a file named "two" with:
      """
      TWO
      """
    And a file named "three" with:
      """
      THREE
      """
    When I run `hub release create -m "m" v1.2.0 -a one -a two -a three`
    Then the exit status should be 1
    Then the stderr should contain exactly:
      """
      Attaching 3 assets...
      The release was created, but attaching 2 assets failed. You can retry with:
      hub release edit v1.2.0 -m '' -a two -a three
      
      Error uploading release asset: Unprocessable Entity (HTTP 422)\n
      """

  Scenario: Create a release with nonexistent asset
    When I run `hub release create -m "hello" v1.2.0 -a "idontexis.tgz"`
    Then the exit status should be 1
    Then the stderr should contain exactly:
      """
      open idontexis.tgz: no such file or directory\n
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

  Scenario: Create release no tag
    When I run `hub release create -m hello`
    Then the exit status should be 1
    Then the stderr should contain "hub release create"

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
    Then the output should not contain anything

  Scenario: Edit existing release when there is a fork
    Given the "doge" remote has url "git://github.com/doge/will_paginate.git"
    And I am on the "feature" branch with upstream "doge/feature"
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { url: 'https://api.github.com/repos/mislav/will_paginate/releases/123',
            tag_name: 'v1.2.0',
          },
        ]
      }
      patch('/repos/mislav/will_paginate/releases/123') {
        json({})
      }
      """
    When I successfully run `hub release edit -m "" v1.2.0`
    Then the output should not contain anything

  Scenario: Edit existing release no title
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            name: 'will_paginate 1.2.0',
          },
        ]
      }
      """
    And a file named "message.txt" with:
      """
      """
    When I run `hub release edit v1.2.0 -F message.txt`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Aborting editing due to empty release title\n
      """

  Scenario: Edit existing release by uploading assets
    Given the GitHub API server:
      """
      deleted = false
      get('/repos/mislav/will_paginate/releases') {
        json [
          { url: 'https://api.github.com/repos/mislav/will_paginate/releases/123',
            upload_url: 'https://uploads.github.com/uploads/assets{?name,label}',
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
      post('/uploads/assets', :host_name => 'uploads.github.com') {
        halt 422 unless deleted
        halt 401 unless request.env['HTTP_AUTHORIZATION'] == 'token OTOKEN'
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
      Attaching 1 asset...\n
      """

  Scenario: Edit release no tag
    When I run `hub release edit -m hello`
    Then the exit status should be 1
    Then the stderr should contain "hub release edit"

  Scenario: Download a release asset
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
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

  Scenario: Download release assets that match pattern
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            assets: [
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9876',
                name: 'hello-amd32-1.2.0.tar.gz',
              },
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9877',
                name: 'hello-amd64-1.2.0.tar.gz',
              },
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9878',
                name: 'hello-x86-1.2.0.tar.gz',
              },
            ],
          },
        ]
      }
      get('/repos/mislav/will_paginate/assets/9876') { "TARBALL" }
      get('/repos/mislav/will_paginate/assets/9877') { "TARBALL" }
      """
      When I successfully run `hub release download v1.2.0 --include '*amd*'`
      Then the output should contain exactly:
        """
        Downloading hello-amd32-1.2.0.tar.gz ...
        Downloading hello-amd64-1.2.0.tar.gz ...\n
        """
      And the file "hello-x86-1.2.0.tar.gz" should not exist

  Scenario: Glob pattern allows exact match
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            assets: [
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9876',
                name: 'hello-amd32-1.2.0.tar.gz',
              },
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9877',
                name: 'hello-amd64-1.2.0.tar.gz',
              },
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9878',
                name: 'hello-x86-1.2.0.tar.gz',
              },
            ],
          },
        ]
      }
      get('/repos/mislav/will_paginate/assets/9876') { "ASSET_TARBALL" }
      """
      When I successfully run `hub release download v1.2.0 --include hello-amd32-1.2.0.tar.gz`
      Then the output should contain exactly:
        """
        Downloading hello-amd32-1.2.0.tar.gz ...\n
        """
      And the file "hello-amd32-1.2.0.tar.gz" should contain exactly:
        """
        ASSET_TARBALL
        """
      And the file "hello-amd64-1.2.0.tar.gz" should not exist
      And the file "hello-x86-1.2.0.tar.gz" should not exist

  Scenario: Advanced glob pattern
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            assets: [
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9876',
                name: 'hello-amd32-1.2.0.tar.gz',
              },
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9876',
                name: 'hello-amd32-1.2.1.tar.gz',
              },
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9876',
                name: 'hello-amd32-1.2.2.tar.gz',
              },
            ],
          },
        ]
      }
      get('/repos/mislav/will_paginate/assets/9876') { "ASSET_TARBALL" }
      """
      When I successfully run `hub release download v1.2.0 --include '*-amd32-?.?.[01].tar.gz'`
      Then the output should contain exactly:
        """
        Downloading hello-amd32-1.2.0.tar.gz ...
        Downloading hello-amd32-1.2.1.tar.gz ...\n
        """

  Scenario: No matches for download pattern
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        json [
          { tag_name: 'v1.2.0',
            assets: [
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9876',
                name: 'hello-amd32-1.2.0.tar.gz',
              },
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9876',
                name: 'hello-amd32-1.2.1.tar.gz',
              },
              { url: 'https://api.github.com/repos/mislav/will_paginate/assets/9876',
                name: 'hello-amd32-1.2.2.tar.gz',
              },
            ],
          },
        ]
      }
      """
      When I run `hub release download v1.2.0 --include amd32`
      Then the exit status should be 1
      Then the stderr should contain exactly:
        """
        the `--include` pattern did not match any available assets:
        hello-amd32-1.2.0.tar.gz
        hello-amd32-1.2.1.tar.gz
        hello-amd32-1.2.2.tar.gz\n
        """

  Scenario: Download release no tag
    When I run `hub release download`
    Then the exit status should be 1
    Then the stderr should contain "hub release download"

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
    Then the output should not contain anything

  Scenario: Release not found
    Given the GitHub API server:
      """
      get('/repos/mislav/will_paginate/releases') {
        assert :per_page => "100"
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
    When I run `hub release delete v2.0`
    Then the exit status should be 1
    And the stderr should contain exactly:
      """
      Unable to find release with tag name `v2.0'\n
      """

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
