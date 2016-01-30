Feature: hub release

  Background:
    Given I am in "git://github.com/mislav/will_paginate.git" git repo
    And I am "mislav" on github.com with OAuth token "OTOKEN"

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
    When I successfully run `hub release`
    Then the output should contain exactly:
      """
      v1.2.0-pre
      v1.0.2\n
      """
