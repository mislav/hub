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
      will_paginate 1.2.0 (v1.2.0)

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
      will_paginate 1.2.0 (v1.2.0)

      ### Hello to my release

      Here is what's broken:
      - everything

      ## Downloads

      https://github.com/mislav/will_paginate/releases/download/v1.2.0/example.zip
      https://github.com/mislav/will_paginate/archive/v1.2.0.zip
      https://github.com/mislav/will_paginate/archive/v1.2.0.tar.gz\n
      """
