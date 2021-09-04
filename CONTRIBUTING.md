Contributing to hub
===================

Contributions to this project are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) to the public under the [project's open source license](LICENSE).

This project adheres to a [Code of Conduct][code-of-conduct]. By participating, you are expected to uphold this code.

[code-of-conduct]: ./CODE_OF_CONDUCT.md

You will need:

1. Go 1.11+
1. Ruby 1.9+ with Bundler
2. git 1.8+
3. tmux & zsh (optional) - for running shell completion tests

If setting up either Go or Ruby for development proves to be a pain, you can
run the test suite in a prepared Docker container via `script/docker`.

## What makes a good hub feature

hub is a tool that wraps git to provide useful integration with GitHub. A new
feature is a good idea for hub if it improves some workflow for a GitHub user.

* A feature that encapsulates a git workflow *not specific* to GitHub is **not**
  a good fit for hub, since something like that is best implemented as an
  external script.
* If you're proposing to add a new custom command such as `hub foo`, please
  research if there's a possibility that such a custom command could conflict
  with other commands from popular 3rd party git projects.
* If your contribution fixes a security vulnerability, please refer to the [SECURITY.md](./.github/SECURITY.md) security policy file

## How to install dependencies and run tests

1. [Clone the project](./README.md#source)
2. Verify that existing tests pass:
    `make test-all`
3. Create a topic branch:
    `git checkout -b feature`
4. **Make your changes.**
   (It helps a lot if you write tests first.)
5. Verify that the tests still pass.
6. Fork the project on GitHub:
    `make && bin/hub fork --remote-name=origin`
7. Push to your fork:
    `git push -u origin HEAD`
8. Open a pull request describing your changes:
    `bin/hub pull-request`

Vendored Go dependencies are managed with [`go mod`](https://github.com/golang/go/wiki/Modules).
Check `go help mod` for information on how to add or update a vendored
dependency.

## How to write tests

Go unit tests are in `*_test.go` files and are runnable with `make test`. These
run really fast (under 10s).

However, most hub functionality is exercised through integration-style tests
written in Cucumber. See [Features](./features) for more info.
