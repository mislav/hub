Contributing to hub
===================

This project adheres to the [Open Code of Conduct][code-of-conduct]. By participating, you are expected to uphold this code.

[code-of-conduct]: http://todogroup.org/opencodeofconduct/#Hub/opensource@github.com

You will need:

1. Go 1.8+
1. Ruby 1.9+ with Bundler
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

1. Clone hub:
    `git clone https://github.com/github/hub.git && cd hub`
2. Verify that existing tests pass:
    `make test-all`
3. Create a topic branch:
    `git checkout -b feature`
4. **Make your changes.**
   (It helps a lot if you write tests first.)
5. Verify that the tests still pass.
6. Fork hub on GitHub (adds a remote named "YOUR-USER"):
    `make && bin/hub fork`
7. Push to your fork:
    `git push -u <YOUR-USER> HEAD`
8. Open a pull request describing your changes:
    `bin/hub pull-request`

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
