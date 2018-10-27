Feature: hub release

  Background:
    Given I am in "git://github.com/mislav/will_paginate.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

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
