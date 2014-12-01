package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	flagSet     flag.FlagSet
	buildFlag   bool
	versionFlag bool
)

func init() {
	flagSet.BoolVar(&buildFlag, "b", false, "build Hub")
	flagSet.BoolVar(&versionFlag, "v", false, "show Hub version")
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: REPO_PATH [-b] [-v]\n\n")
}

func main() {
	flag.Usage = showUsage
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("make: ")

	args := flag.Args()
	if len(args) < 2 {
		showUsage()
		os.Exit(1)
	}

	path := args[0]
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	flagSet.Usage = showUsage
	flagSet.Parse(args[1:])

	if buildFlag {
		err = build(path)
	} else if versionFlag {
		err = version(path)
	} else {
		showUsage()
	}

	if err != nil {
		log.Fatal(err)
	}
}
