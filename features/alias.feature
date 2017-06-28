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
      # Wrap git automatically by adding the following to ~/.config/fish/functions/git.fish:

      function git --wraps hub --description 'Alias for hub, which wraps git to provide extra functionality with GitHub.'
          hub $argv
      end\n
      """
  
  Scenario: rc instructions
    Given $SHELL is "/usr/local/bin/rc"
    When I successfully run `hub alias`
    Then the output should contain exactly:
      """
      # Wrap git automatically by adding the following to $home/lib/profile:

      eval `{hub alias -s}\n
      """

  Scenario: zsh instructions
    Given $SHELL is "/bin/zsh"
    When I successfully run `hub alias`
    Then the output should contain exactly:
      """
      # Wrap git automatically by adding the following to ~/.zshrc:

      eval "$(hub alias -s)"\n
      """

  Scenario: csh instructions
    Given $SHELL is "/bin/csh"
    When I successfully run `hub alias`
    Then the output should contain exactly:
      """
      # Wrap git automatically by adding the following to ~/.cshrc:

      eval "`hub alias -s`"\n
      """

  Scenario: tcsh instructions
    Given $SHELL is "/bin/tcsh"
    When I successfully run `hub alias`
    Then the output should contain exactly:
      """
      # Wrap git automatically by adding the following to ~/.tcshrc:

      eval "`hub alias -s`"\n
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
  
  Scenario: rc code
    Given $SHELL is "/usr/local/bin/rc"
    When I successfully run `hub alias -s`
    Then the output should contain exactly:
      """
      fn git { builtin hub $* }\n
      """  

  Scenario: zsh code
    Given $SHELL is "/bin/zsh"
    When I successfully run `hub alias -s`
    Then the output should contain exactly:
      """
      alias git=hub\n
      """

  Scenario: csh code
    Given $SHELL is "/bin/csh"
    When I successfully run `hub alias -s`
    Then the output should contain exactly:
      """
      alias git hub\n
      """

  Scenario: tcsh code
    Given $SHELL is "/bin/tcsh"
    When I successfully run `hub alias -s`
    Then the output should contain exactly:
      """
      alias git hub\n
      """

  Scenario: unsupported shell
    Given $SHELL is "/bin/zwoosh"
    When I run `hub alias -s`
    Then the output should contain exactly:
      """
      hub alias: unsupported shell
      supported shells: bash zsh sh ksh csh tcsh fish rc\n
      """
    And the exit status should be 1

  Scenario: unknown shell
    Given $SHELL is ""
    When I run `hub alias`
    Then the output should contain exactly:
      """
      Error: couldn't detect shell type. Please specify your shell with `hub alias <shell>`\n
      """
    And the exit status should be 1

  Scenario: unknown shell output
    Given $SHELL is ""
    When I run `hub alias -s`
    Then the output should contain exactly:
      """
      Error: couldn't detect shell type. Please specify your shell with `hub alias -s <shell>`\n
      """
    And the exit status should be 1
