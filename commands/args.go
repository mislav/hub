package commands

type Args struct {
	args []string
}

func (a *Args) Get(i int) string {
	return a.args[i]
}

func (a *Args) First() string {
	return a.args[0]
}

func (a *Args) Rest() []string {
	return a.args[1:]
}

func (a *Args) Remove(i int) string {
	newArgs, item := removeItem(a.args, 0)
	a.args = newArgs

	return item
}

func (a *Args) Size() int {
	return len(a.args)
}

func NewArgs(args []string) *Args {
	return &Args{args}
}
