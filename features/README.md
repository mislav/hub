# Cucumber features for hub

How to run all features:

```sh
make bin/cucumber
bin/cucumber
```

Because this can take a couple of minutes, you may want to only run select files
related to the functionality that you're developing:

```sh
bin/cucumber feature/api.feature
```

The Cucumber test suite requires a Ruby development environment. If you want to
avoid setting that up, you can run tests inside a Docker container:

```sh
script/docker feature/api.feature
```

## How it works

Each scenario is actually making real invocations to `hub` on the command-line
in the context of a real (dynamically created) git repository.

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

## How to write new tests

The best way to learn to write new tests is to study the existing scenarios for
commands that are similar to those that you want to add or change.

Since Cucumber tests are written in a natural language, you mostly don't need to
know Ruby to write new tests.
