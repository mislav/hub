package commands

import (
	"fmt"
)

type Args struct {
	args []string
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

func NewArgs(args []string) *Args {
	return &Args{args}
}

func removeItem(slice []string, index int) (newSlice []string, item string) {
	if index > len(slice)-1 {
		panic(fmt.Sprintf("Index %d is out of bound", index))
	}

	item = slice[index]
	newSlice = append(slice[:index], slice[index+1:]...)

	return newSlice, item
}
