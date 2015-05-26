# fish completion for hub
# Since git is aliased as hub, need to check if `command hub` creates any conflict with git command completions

# statement starting with 'hub'
function __fish_hub_needs_command
    set cmd (commandline -opc)
    if [ (count $cmd) -eq 1 -a $cmd[1] = 'hub' ]
        return 0
    end
    return 1
end

# statement starting with 'hub <command>'
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
complete -f -c hub -n '__fish_hub_needs_command' -a help -d 'Display enhanced git-help(1)'

# fork
complete -f -c hub -n '__fish_hub_needs_command' -a fork -d 'Fork the original project on GitHub and add a new remote for it under your username.'
complete -f -c hub -n '__fish_hub_using_command fork' -l no-remote -d 'Fork the original project with no remote'
