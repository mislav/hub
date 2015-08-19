package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/utils"
)

var cmdPr = &Command{
	Run:          pr,
	GitExtension: true,
	Usage:        "pr PULLREQ-NUMBER",
	Short:        "Treat pullrequest number",
	Long: `
`,
}

var lengthParams int
var prRegex *regexp.Regexp

const NotPr = "--not-pr"
const prRegexText = "^#?[1-9][0-9]*$"

func init() {
	CmdRunner.Use(cmdPr)
}

func pr(command *Command, args *Args) {
	if !args.IsParamsEmpty() {
		transformPrArgs(args)
	}
}

func transformPrArgs(args *Args) {
	//Remove firstParam
	firstParam := args.FirstParam()
	args.RemoveParam(0)
	//Attatch a protection to not pullrequest numbers
	notPr := false
	idx := 0
	prRegex = regexp.MustCompile(prRegexText)
	for idx < args.ParamsSize() {
		param := args.GetParam(idx)
		if param == NotPr {
			notPr = true
			args.RemoveParam(idx)
			continue
		}
		if notPr && prRegex.MatchString(param) {
			args.ReplaceParam(idx, NotPr+param)
		}
		idx++
		notPr = false
	}
	//Handle args with firstParam
	if strings.HasPrefix(firstParam, "--") {
		action := strings.TrimPrefix(firstParam, "--")
		switch action {
		case "apply", "am":
			convertPrToUrl(args)
			args.Executable = "hub"
			args.Command = action
		case "browse":
			convertPrToUrl(args)
			args.Executable = "open"
			args.Command = ""
		default:
			utils.Check(fmt.Errorf("Error: command not found"))
		}
	} else {
		utils.Check(fmt.Errorf("Error: command not found"))
	}
	//Detach protections
	notPrRegex := regexp.MustCompile("^--not-pr#?[1-9][0-9]*$")
	lengthParams = args.ParamsSize()
	for i := 0; i < lengthParams; i++ {
		param := args.GetParam(i)
		if notPrRegex.MatchString(param) {
			args.ReplaceParam(i, strings.TrimPrefix(param, "--not-pr"))
		}
	}
}

func convertPrToUrl(args *Args) {
	localRepo, err := github.LocalRepo()
	utils.Check(err)
	remote, err := localRepo.OriginRemote()
	utils.Check(err)
	project, err := remote.Project()
	utils.Check(err)
	lengthParams = args.ParamsSize()
	prRegex = regexp.MustCompile(prRegexText)
	for i := 0; i < lengthParams; i++ {
		param := args.GetParam(i)
		if prRegex.MatchString(param) {
			args.ReplaceParam(i, "https://github.com/"+project.Owner+"/"+project.Name+"/pull/"+strings.TrimPrefix(param, "#"))
		}
	}
}
