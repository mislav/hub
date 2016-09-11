package commands

import (
	"fmt"
	"strings"

	"github.com/github/hub/cmd"
)

type Args struct {
	Executable  string
	GlobalFlags []string
	Command     string
	ProgramPath string
	Params      []string
	beforeChain []*cmd.Cmd
	afterChain  []*cmd.Cmd
	Noop        bool
	Terminator  bool
	noForward   bool
	Callbacks   []func() error
}

func (a *Args) Words() []string {
	aa := make([]string, 0)
	for _, p := range a.Params {
		if !looksLikeFlag(p) {
			aa = append(aa, p)
		}
	}

	return aa
}

func (a *Args) Before(command ...string) {
	a.beforeChain = append(a.beforeChain, cmd.NewWithArray(command))
}

func (a *Args) After(command ...string) {
	a.afterChain = append(a.afterChain, cmd.NewWithArray(command))
}

func (a *Args) AfterFn(fn func() error) {
	a.Callbacks = append(a.Callbacks, fn)
}

func (a *Args) NoForward() {
	a.noForward = true
}

func (a *Args) Replace(executable, command string, params ...string) {
	a.Executable = executable
	a.Command = command
	a.Params = params
	a.GlobalFlags = []string{}
	a.noForward = false
}

func (a *Args) Commands() []*cmd.Cmd {
	result := a.beforeChain

	if !a.noForward {
		result = append(result, a.ToCmd())
	}

	result = append(result, a.afterChain...)
	return result
}

func (a *Args) ToCmd() *cmd.Cmd {
	c := cmd.New(a.Executable)
	c.WithArgs(a.GlobalFlags...)

	if a.Command != "" {
		c.WithArg(a.Command)
	}

	for _, arg := range a.Params {
		if arg != "" {
			c.WithArg(arg)
		}
	}

	return c
}

func (a *Args) GetParam(i int) string {
	return a.Params[i]
}

func (a *Args) FirstParam() string {
	if a.ParamsSize() == 0 {
		panic(fmt.Sprintf("Index 0 is out of bound"))
	}

	return a.Params[0]
}

func (a *Args) LastParam() string {
	if a.ParamsSize()-1 < 0 {
		panic(fmt.Sprintf("Index %d is out of bound", a.ParamsSize()-1))
	}

	return a.Params[a.ParamsSize()-1]
}

func (a *Args) HasSubcommand() bool {
	return !a.IsParamsEmpty() && a.Params[0][0] != '-'
}

func (a *Args) InsertParam(i int, items ...string) {
	if i < 0 {
		panic(fmt.Sprintf("Index %d is out of bound", i))
	}

	if i > a.ParamsSize() {
		i = a.ParamsSize()
	}

	newParams := make([]string, 0)
	newParams = append(newParams, a.Params[:i]...)
	newParams = append(newParams, items...)
	newParams = append(newParams, a.Params[i:]...)

	a.Params = newParams
}

func (a *Args) RemoveParam(i int) string {
	newParams, item := removeItem(a.Params, i)
	a.Params = newParams

	return item
}

func (a *Args) ReplaceParam(i int, item string) {
	if i < 0 || i > a.ParamsSize()-1 {
		panic(fmt.Sprintf("Index %d is out of bound", i))
	}

	a.Params[i] = item
}

func (a *Args) IndexOfParam(param string) int {
	for i, p := range a.Params {
		if p == param {
			return i
		}
	}

	return -1
}

func (a *Args) ParamsSize() int {
	return len(a.Params)
}

func (a *Args) IsParamsEmpty() bool {
	return a.ParamsSize() == 0
}

func (a *Args) PrependParams(params ...string) {
	a.Params = append(params, a.Params...)
}

func (a *Args) AppendParams(params ...string) {
	a.Params = append(a.Params, params...)
}

func (a *Args) HasFlags(flags ...string) bool {
	for _, f := range flags {
		if i := a.IndexOfParam(f); i != -1 {
			return true
		}
	}

	return false
}

func NewArgs(args []string) *Args {
	var (
		command     string
		params      []string
		noop        bool
		globalFlags []string
	)

	slurpGlobalFlags(&args, &globalFlags)
	noop = removeValue(&globalFlags, noopFlag)

	if len(args) == 0 {
		params = []string{}
	} else {
		command = args[0]
		params = args[1:]
	}

	return &Args{
		Executable:  "git",
		GlobalFlags: globalFlags,
		Command:     command,
		Params:      params,
		Noop:        noop,
		beforeChain: make([]*cmd.Cmd, 0),
		afterChain:  make([]*cmd.Cmd, 0),
	}
}

const (
	noopFlag    = "--noop"
	versionFlag = "--version"
	helpFlag    = "--help"
	configFlag  = "-c"
	chdirFlag   = "-C"
	flagPrefix  = "-"
)

func looksLikeFlag(value string) bool {
	return strings.HasPrefix(value, flagPrefix)
}

func slurpGlobalFlags(args *[]string, globalFlags *[]string) {
	slurpNextValue := false
	commandIndex := 0

	for i, arg := range *args {
		if slurpNextValue {
			commandIndex = i + 1
			slurpNextValue = false
		} else if arg == versionFlag || arg == helpFlag || !looksLikeFlag(arg) {
			break
		} else {
			commandIndex = i + 1
			if arg == configFlag || arg == chdirFlag {
				slurpNextValue = true
			}
		}
	}

	if commandIndex > 0 {
		aa := *args
		*globalFlags = aa[0:commandIndex]
		*args = aa[commandIndex:]
	}
}

func removeItem(slice []string, index int) (newSlice []string, item string) {
	if index < 0 || index > len(slice)-1 {
		panic(fmt.Sprintf("Index %d is out of bound", index))
	}

	item = slice[index]
	newSlice = append(slice[:index], slice[index+1:]...)

	return newSlice, item
}

func removeValue(slice *[]string, value string) (found bool) {
	aa := *slice
	for i := len(aa) - 1; i >= 0; i-- {
		arg := aa[i]
		if arg == value {
			found = true
			*slice, _ = removeItem(*slice, i)
		}
	}
	return found
}
