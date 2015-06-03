# fish completion for hub
# Since git is aliased as hub, need to check if `command hub` creates any conflict with git command completions

# TODO
# Snippet type thing. 'hub create [name]' where [name] should be in grey and suppress the smart suggesstions from fish.
# 'hub help <tab>' should show a list of commands that help is available for.
# '__fish_hub_suppress_files' function.

function __fish_hub_suppress_files
end

# statement starting with 'hub'
function __fish_hub_needs_command
    set cmd (commandline -opc)
    if [ (count $cmd) -eq 1 -a $cmd[1] = 'hub' ]
        return 0
    end
    return 1
end

# statement starting with 'hub COMMAND'
function __fish_hub_using_command
    set cmd (commandline -opc)
    if [ (count $cmd) -gt 1 ]
        if [ $argv[1] = $cmd[2] ]
            return 0
        end
    end
    return 1
end

# help
# show '--help' for every command
complete -f -c hub -n '__fish_hub_needs_command' -a help -d 'Display enhanced git-help(1)'
complete -f -c hub -n 'not __fish_hub_needs_command' -l help -d 'Display enhanced git-help(1)'
#complete -f -c hub -n 'not __fish_hub_needs_command' -l help -d 'Display enhanced git-help(1)'

# create
complete -f -c hub -n '__fish_hub_needs_command' -a create -d 'Create repository on Github and add Github as remote'
complete -f -c hub -n '__fish_hub_using_command create' -s p -d 'Create private repository'
complete -f -c hub -n '__fish_hub_using_command create' -s d -d 'Set description of repository'
complete -f -c hub -n '__fish_hub_using_command create' -s h -d 'Set homepage of repository'

# browse
complete -f -c hub -n '__fish_hub_needs_command' -a browse -d 'Open a GitHub page in the default browser'
complete -f -c hub -n '__fish_hub_using_command browse' -s u -d 'Output the URL rather than opening the browser'

# compare
complete -f -c hub -n '__fish_hub_needs_command' -a compare -d 'Open compare page on GitHub'
complete -f -c hub -n '__fish_hub_using_command compare' -s u -d 'Output the URL rather than opening the browser'

# fork
complete -f -c hub -n '__fish_hub_needs_command' -a fork -d 'Fork the original project on GitHub under your username'
complete -f -c hub -n '__fish_hub_using_command fork' -l no-remote -d 'Fork the original project with no remote'

# pull-request
complete -f -c hub -n '__fish_hub_needs_command' -a pull-request -d 'Open a pull request on GitHub'
complete -f -c hub -n '__fish_hub_using_command pull-request' -s o -l browse -d 'Open pull-request page on GitHub'
complete -f -c hub -n '__fish_hub_using_command pull-request' -s f -d 'Skip checking local commits that are not yet pushed'
complete -f -c hub -n '__fish_hub_using_command pull-request' -s m -d 'Message for pull-request'
complete -c hub -n '__fish_hub_using_command pull-request' -s f -d 'Message for pull-request from file'
complete -f -c hub -n '__fish_hub_using_command pull-request' -s i -d 'Provide issue number/issue URL for attaching to an existing pull-request'
complete -f -c hub -n '__fish_hub_using_command pull-request' -s b -d 'Specify BASE for pull-request'
complete -f -c hub -n '__fish_hub_using_command pull-request' -s h -d 'Specify HEAD for pull-request'

# ci-status
complete -f -c hub -n '__fish_hub_needs_command' -a ci-status -d 'Looks up the SHA for commit in GitHub Status API and displays the latest status'
complete -f -c hub -n '__fish_hub_needs_command ci-status' -s v -d 'Print the URL to CI build results'
