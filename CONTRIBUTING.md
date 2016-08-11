Contributing to hub
===================

You will need:

1. Go 1.6 or better
1. Ruby 1.9+
2. git 1.8+
3. tmux & zsh (optional) - for running shell completion tests

We will also require you to sign the [GitHub Contributor License Agreement](https://cla.github.com/)
after you submit your first pull request to this project. The link to sign the
agreement will be presented to you in the web interface of the pull request.

## What makes a good hub feature

hub is a tool that wraps git to provide useful integration with GitHub. A new
feature is a good idea for hub if it improves some workflow for a GitHub user.

* A feature that encapsulates a git workflow *not specific* to GitHub is **not**
  a good fit for hub, since something like that is best implemented as an
  external script.
* If you're proposing to add a new custom command such as `hub foo`, please
  research if there's a possibility that such a custom command could conflict
  with other commands from popular 3rd party git projects.

## How to install dependencies and run tests

These instructions assume that _you already have hub installed_ and aliased as
`git` (see "Aliasing").

1. Clone hub:
    `git clone github/hub && cd hub`
1. Install necessary development dependencies:
    `script/bootstrap`
2. Verify that existing tests pass:
    `script/test`
3. Create a topic branch:
    `git checkout -b feature`
4. **Make your changes.**
   (It helps a lot if you write tests first.)
5. Verify that tests still pass:
    `script/test`
6. Fork hub on GitHub (adds a remote named "YOUR-USER"):
    `git fork`
7. Push to your fork:
    `git push <YOUR-USER> HEAD`
8. Open a pull request describing your changes:
    `git pull-request`

## How to write tests

The new test suite is written in Cucumber under `features/` directory. Each
scenario is actually making real invocations to `hub` on the command-line in the
context of a real (dynamically created) git repository.

Whenever a scenario requires talking to the GitHub API, a fake HTTP server is
spun locally to replace the real GitHub API. This is done so that the test suite
runs faster and is available offline as well. The fake API server is defined
as a Sinatra app inline in each scenario:

```
Given the GitHub API server:
  """
  post('/repos/github/hub/pulls') {
    status 200
  }
  """
```

The best way to learn to write new tests is to study the existing scenarios for
commands that are similar to those that you want to add or change.
