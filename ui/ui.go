package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

type UI interface {
	Promptf(format string, a ...interface{}) (n int, err error)
	Promptln(a ...interface{}) (n int, err error)

	Printf(format string, a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)

	Errorf(format string, a ...interface{}) (n int, err error)
	Errorln(a ...interface{}) (n int, err error)
}

func init() {
	console := Console{
		Stdout: Stdout,
		Stderr: Stderr,
		TTY:    Stdout,
	}

	TTY = os.Stdin
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0644)
	if err == nil {
		TTY = tty
		console.TTY = tty
	}

	Default = console
}

var (
	TTY     *os.File
	Stdout  = colorable.NewColorableStdout()
	Stderr  = colorable.NewColorableStderr()
	Default UI
)

func Promptf(format string, a ...interface{}) (n int, err error) {
	return Default.Promptf(format, a...)
}

func Promptln(a ...interface{}) (n int, err error) {
	return Default.Promptln(a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return Default.Printf(format, a...)
}

func Println(a ...interface{}) (n int, err error) {
	return Default.Println(a...)
}

func Errorf(format string, a ...interface{}) (n int, err error) {
	return Default.Errorf(format, a...)
}

func Errorln(a ...interface{}) (n int, err error) {
	return Default.Errorln(a...)
}

func IsTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd())
}

type Console struct {
	TTY    io.Writer
	Stdout io.Writer
	Stderr io.Writer
}

func (c Console) Promptf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.TTY, format, a...)
}

func (c Console) Promptln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.TTY, a...)
}

func (c Console) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stdout, format, a...)
}

func (c Console) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.Stdout, a...)
}

func (c Console) Errorf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stderr, format, a...)
}

func (c Console) Errorln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.Stderr, a...)
}
