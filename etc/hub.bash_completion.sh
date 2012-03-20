# hub tab-completion script for bash.
# This script complements the completion script that ships with git.

# Check that git tab completion is available
if declare -F _git > /dev/null; then
  # Duplicate and rename the 'list_all_commands' function
  eval "$(declare -f __git_list_all_commands | \
        sed 's/__git_list_all_commands/__git_list_all_commands_without_hub/')"

  # Wrap the 'list_all_commands' function with extra hub commands
  __git_list_all_commands() {
    cat <<-EOF
alias
pull-request
fork
create
browse
compare
EOF
    __git_list_all_commands_without_hub
  }

  # Ensure cached commands are cleared
  __git_all_commands=""
fi
