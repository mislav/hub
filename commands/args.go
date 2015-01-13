package commands

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/github/hub/cmd"

	flag "github.com/github/hub/Godeps/_workspace/src/github.com/ogier/pflag"
)

type Args struct {
	Executable  string
	ConfigParam map[string]string
	Command     string
	Params      []string
	beforeChain []*cmd.Cmd
	afterChain  []*cmd.Cmd
	Noop        bool
}

func (a *Args) Words() []string {
	aa := make([]string, 0)
	for _, p := range a.Params {
		if !strings.HasPrefix(p, "-") {
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

func (a *Args) Replace(executable, command string, params ...string) {
	a.Executable = executable
	a.Command = command
	a.Params = params
}

func (a *Args) Commands() []*cmd.Cmd {
	result := a.beforeChain
	result = append(result, a.ToCmd())
	result = append(result, a.afterChain...)

	return result
}

func (a *Args) ToCmd() *cmd.Cmd {
	c := cmd.New(a.Executable)
	args := make([]string, 0)

	if a.Command != "" {
		args = append(args, a.Command)
	}

	for _, arg := range a.Params {
		if arg != "" {
			args = append(args, arg)
		}
	}

	return c.WithArgs(args...)
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
	if i < 0 || (i != 0 && i > a.ParamsSize()-1) {
		panic(fmt.Sprintf("Index %d is out of bound", i))
	}

	newParams := []string{}
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
		globalFlag  flag.FlagSet
		version     bool
		help        bool
		noop        bool
		configParam mapValue = make(mapValue)
	)

	globalFlag.BoolVarP(&version, "version", "", false, "")
	globalFlag.BoolVarP(&help, "help", "", false, "")
	globalFlag.BoolVarP(&noop, "noop", "", false, "")
	globalFlag.VarP(configParam, "", "c", "")

	globalFlag.SetOutput(ioutil.Discard)
	globalFlag.Init("hub", flag.ContinueOnError)
	aa := make([]string, 0)
	err := globalFlag.Parse(args)
	if err == nil {
		aa = globalFlag.Args()
	} else {
		aa = args
	}

	if version {
		aa = append([]string{"version"}, aa...)
	}

	if help {
		aa = append([]string{"help"}, aa...)
	}

	var (
		command string
		params  []string
	)

	if len(aa) == 0 {
		params = []string{}
	} else {
		command = aa[0]
		params = aa[1:]
	}

	return &Args{
		Executable:  "git",
		ConfigParam: configParam,
		Command:     command,
		Params:      params,
		Noop:        noop,
		beforeChain: make([]*cmd.Cmd, 0),
		afterChain:  make([]*cmd.Cmd, 0),
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
