Feature: hub alias

  Scenario: bash instructions
    Given $SHELL is "/bin/bash"
    When I successfully run `hub alias`
    Then the output should contain exactly:
      """
      # Wrap git automatically by adding the following to ~/.bash_profile:

      eval "$(hub alias -s)"\n
      """

  Scenario: fish instructions
    Given $SHELL is "/usr/local/bin/fish"
    When I successfully run `hub alias`
    Then the output should contain exactly:
      """
      # Wrap git automatically by adding the following to ~/.config/fish/config.fish:

      eval (hub alias -s)\n
      """

  Scenario: zsh instructions
    Given $SHELL is "/bin/zsh"
    When I successfully run `hub alias`
    Then the output should contain exactly:
      """
      # Wrap git automatically by adding the following to ~/.zshrc:

      eval "$(hub alias -s)"\n
      """

  Scenario: bash code
    Given $SHELL is "/bin/bash"
    When I successfully run `hub alias -s`
    Then the output should contain exactly:
      """
      alias git=hub\n
      """

  Scenario: fish code
    Given $SHELL is "/usr/local/bin/fish"
    When I successfully run `hub alias -s`
    Then the output should contain exactly:
      """
      alias git=hub\n
      """

  Scenario: zsh code
    Given $SHELL is "/bin/zsh"
    When I successfully run `hub alias -s`
    Then the output should contain exactly:
      """
      alias git=hub\n
      """

  Scenario: unsupported shell
    Given $SHELL is "/bin/zwoosh"
    When I run `hub alias -s`
    Then the output should contain exactly:
      """
      hub alias: unsupported shell
      supported shells: bash zsh sh ksh csh fish\n
      """
    And the exit status should be 1
