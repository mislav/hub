package commands

import (
	"fmt"
	"github.com/jingweno/gh/cmd"
)

type Args struct {
	args        []string
	beforeChain []*cmd.Cmd
	afterChain  []*cmd.Cmd
}

func (a *Args) Before(command ...string) {
	a.beforeChain = append(a.beforeChain, cmd.NewWithArray(command))
}

func (a *Args) After(command ...string) {
	a.afterChain = append(a.afterChain, cmd.NewWithArray(command))
}

func (a *Args) Commands() []*cmd.Cmd {
	result := a.beforeChain
	result = append(result, a.ToCmd())
	result = append(result, a.afterChain...)

	return result
}

func (a *Args) ToCmd() *cmd.Cmd {
	return cmd.New("git").WithArgs(a.Array()...)
}

func (a *Args) Get(i int) string {
	return a.args[i]
}

func (a *Args) First() string {
	return a.args[0]
}

func (a *Args) Last() string {
	return a.args[a.Size()-1]
}

func (a *Args) Rest() []string {
	return a.args[1:]
}

func (a *Args) Remove(i int) string {
	newArgs, item := removeItem(a.args, i)
	a.args = newArgs

	return item
}

func (a *Args) Replace(i int, item string) {
	if i > a.Size()-1 {
		panic(fmt.Sprintf("Index %d is out of bound", i))
	}

	a.args[i] = item
}

func (a *Args) IndexOf(arg string) int {
	for i, aa := range a.args {
		if aa == arg {
			return i
		}
	}

	return -1
}

func (a *Args) Size() int {
	return len(a.args)
}

func (a *Args) IsEmpty() bool {
	return a.Size() == 0
}

func (a *Args) Array() []string {
	return a.args
}

func (a *Args) Append(args ...string) {
	a.args = append(a.args, args...)
}

func (a *Args) Prepend(args ...string) {
	a.args = append(args, a.args...)
}

func NewArgs(args []string) *Args {
	return &Args{args, make([]*cmd.Cmd, 0), make([]*cmd.Cmd, 0)}
}

func removeItem(slice []string, index int) (newSlice []string, item string) {
	if index > len(slice)-1 {
		panic(fmt.Sprintf("Index %d is out of bound", index))
	}

	item = slice[index]
	newSlice = append(slice[:index], slice[index+1:]...)

	return newSlice, item
}
