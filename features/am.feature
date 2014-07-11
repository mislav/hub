Feature: hub am
  Scenario: Enterprise repo
    Given I am in "git://git.my.org/mislav/dotfiles.git" git repo
    And I am "mislav" on git.my.org with OAuth token "FITOKEN"
    And "git.my.org" is a whitelisted Enterprise host
    Given the GitHub API server:
      """
      get('/api/v3/repos/mislav/dotfiles/pulls/387') {
        halt 400 unless request.env['HTTP_ACCEPT'] == 'application/vnd.github.v3.patch'
        <<PATCH
From 7eb75a26ee8e402aad79fcf36a4c1461e3ec2592 Mon Sep 17 00:00:00 2001
From: Mislav <mislav.marohnic@gmail.com>
Date: Tue, 24 Jun 2014 11:07:05 -0700
Subject: [PATCH] Create a README
---
diff --git a/README.md b/README.md
new file mode 100644
index 0000000..ce01362
--- /dev/null
+++ b/README.md
+hello
-- 
1.9.3
PATCH
      }
      """
    When I successfully run `hub am -q -3 https://git.my.org/mislav/dotfiles/pull/387`
    And I successfully run `git log -1 --format=%s`
    Then the output should contain exactly "Create a README\n"
