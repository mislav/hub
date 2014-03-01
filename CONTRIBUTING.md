Contributing to hub
===================

<i>**Warning:** in the near future, hub might be implemented
[entirely in Go instead of Ruby](https://github.com/github/hub/issues/475).
Keep that in mind, and don't contribute big features/refactorings to
the Ruby codebase, as such pull requests will be unlikely to get accepted.</i>

You will need:

1. Ruby 1.8.7+
2. git 1.8+
3. tmux & zsh (optional) - for running shell completion tests

## What makes a good hub feature

hub is a tool that wraps git to provide useful integration with GitHub. A new
feature is a good idea for hub only if it relates to *both git and GitHub*.

* A feature that adds GitHub Issues management is **not** a good fit for hub
  since it's not git-related. Please use [ghi](https://github.com/stephencelis/ghi)
  instead.
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
1. Ensure Bundler is installed:  
    `which bundle || gem install bundler`
1. Install development dependencies:  
    `bundle install`
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

## How hub works

1.  [Runner](lib/hub/runner.rb#files) handles the command-line invocation;

2.  [Args](lib/hub/args.rb#files) wraps ARGV for easy access;

3.  [Commands](lib/hub/commands.rb#files) dispatches each command to the
    appropriate method, e.g. `hub pull-request` runs the `pull_request`
    method. Each method processes args as needed, using Context and GitHubAPI
    in the process;

4.  [Context](lib/hub/context.rb#files) handles inspecting the current
    environment and git repository;

5.  [GitHubAPI](lib/hub/github_api.rb#files) handles GitHub API authentication
    and communication;

6.  And finally, Runner receives the resulting arguments to execute in the
    shell by forwarding them to `git`.

## How to write tests

The old test suite for hub was written in test/unit and some legacy tests can
still be found in the `test/` directory. Unless you have a need for writing
super-isolated unit tests, **do not add** any more tests to this suite.

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
