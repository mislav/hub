function __fish_hub_needs_command
  set cmd (commandline -opc)
  if [ (count $cmd) -eq 1 ]
    return 0
  else
    return 1
  end
end

function  __fish_hub_using_command
  set cmd (commandline -opc)
  if [ (count $cmd) -gt 1 ]
    if [ $argv[1] = $cmd[2] ]
      return 0
    end
  end
  return 1
end

complete -f -c hub -n '__fish_hub_needs_command' -a alias -d "show shell instructions for wrapping git"
complete -f -c hub -n '__fish_hub_needs_command' -a browse -d "browse the project on GitHub"
complete -f -c hub -n '__fish_hub_needs_command' -a compare -d "lookup commit in GitHub Status API"
complete -f -c hub -n '__fish_hub_needs_command' -a create -d "create new repo on GitHub for the current project"
complete -f -c hub -n '__fish_hub_needs_command' -a delete -d "delete a GitHub repo"
complete -f -c hub -n '__fish_hub_needs_command' -a fork -d "fork origin repo on GitHub"
complete -f -c hub -n '__fish_hub_needs_command' -a pull-request -d "open a pull request on GitHub"
complete -f -c hub -n '__fish_hub_needs_command' -a pr -d "list or checkout a GitHub release"
complete -f -c hub -n '__fish_hub_needs_command' -a issue -d "list or create a GitHub issue"
complete -f -c hub -n '__fish_hub_needs_command' -a release -d "list or create a GitHub release"
complete -f -c hub -n '__fish_hub_needs_command' -a ci-status -d "display GitHub Status information for a commit"
complete -f -c hub -n '__fish_hub_needs_command' -a sync -d "update local branches from upstream"

# alias
complete -f -c hub -n ' __fish_hub_using_command alias' -a 'bash zsh sh ksh csh fish' -d "output shell script suitable for eval"
# pull-request
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s f -d "Skip the check for unpushed commits"
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s -m -d "Use the first line of <MESSAGE> as pull request title, and the rest as pull request description"
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s F -d "Read the pull request title and description from <FILE>"
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s o -d "Open the new pull request in a web browser"
complete -f -c hub -n ' __fish_hub_using_command pull-request' -l browse -d "Open the new pull request in a web browser"
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s p -d "Push the current branch to <HEAD> before creating the pull request"
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s b -d 'The base branch in "[OWNER:]BRANCH" format'
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s h -d 'The head branch in "[OWNER:]BRANCH" format'
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s a -d 'A comma-separated list of GitHub handles to assign to this pull request'
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s M -d "The milestone name to add to this pull request. Passing the milestone number is deprecated."
complete -f -c hub -n ' __fish_hub_using_command pull-request' -s l -d "Add a comma-separated list of labels to this pull request"
# fork
complete -f -c hub -n ' __fish_hub_using_command fork' -l no-remote -d "Skip adding a git remote for the fork"
# browse
complete -f -c hub -n ' __fish_hub_using_command browse' -s u -d "Print the URL instead of opening it"
complete -f -c hub -n ' __fish_hub_using_command browse' -s c -d "Put the URL in clipboard instead of opening it"
complete -f -c hub -n ' __fish_hub_using_command browse' -a '-- commits' -d 'commits'
complete -f -c hub -n ' __fish_hub_using_command browse' -a '-- contributors' -d 'contributors'
complete -f -c hub -n ' __fish_hub_using_command browse' -a '-- issues' -d 'issues'
complete -f -c hub -n ' __fish_hub_using_command browse' -a '-- pulls' -d 'pull requests'
complete -f -c hub -n ' __fish_hub_using_command browse' -a '-- wiki' -d 'wiki'
# compare
complete -f -c hub -n ' __fish_hub_using_command compare' -s u -d 'Print the URL instead of opening it'
# create
complete -f -c hub -n ' __fish_hub_using_command create' -s o -d "Open the new repository in a web browser"
complete -f -c hub -n ' __fish_hub_using_command create' -l browse -d "Open the new repository in a web browser"
complete -f -c hub -n ' __fish_hub_using_command create' -s p -d "Create a private repository"
complete -f -c hub -n ' __fish_hub_using_command create' -s c -d "Put the URL of the new repository to clipboard instead of printing it."
complete -f -c hub -n ' __fish_hub_using_command create' -l copy -d "Put the URL of the new repository to clipboard instead of printing it."
# delete
complete -f -c hub -n ' __fish_hub_using_command delete' -s y -d "Skip the confirmation prompt"
complete -f -c hub -n ' __fish_hub_using_command delete' -l yes -d "Skip the confirmation prompt"
# ci-status
complete -f -c hub -n ' __fish_hub_using_command ci-status' -s v -d "Print detailed report of all status checks and their URLs"
