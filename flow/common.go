package flow

import (
	"github.com/github/hub/cmd"
	"github.com/github/hub/git"
	ini "gopkg.in/ini.v1"
)

func launchCmdGit(cmdGit [][]string) (err error) {
	for i := range cmdGit {
		err = git.Spawn(cmdGit[i]...)

		if err != nil {
			break
		}
	}

	return
}

func HubCmd(args ...string) (err error) {
	cmd := cmd.New("hub")

	for _, a := range args {
		cmd.WithArg(a)
	}

	return cmd.Spawn()
}

type GitConfig struct {
	cfg        *ini.File
	configFile string

	sectionBranch string
	sectionPrefix string
}

func NewConfig() (gitConfig *GitConfig, err error) {
	file := ".git/config"

	cfg, err := ini.Load(file)

	gitConfig = &GitConfig{
		cfg:           cfg,
		configFile:    file,
		sectionBranch: "gitflow \"branch\"",
		sectionPrefix: "gitflow \"prefix\"",
	}

	return
}

func (gitConfig *GitConfig) GetSection(name string) (section *ini.Section) {
	section = gitConfig.cfg.Section(name)
	return
}

func (gitConfig *GitConfig) CreateSections(args ...string) (err error) {
	for _, name := range args {
		_, err = gitConfig.cfg.NewSection(name)
	}

	return
}

func (gitConfig *GitConfig) CreateKey(sectionName string, key string, value string) {
	section := gitConfig.cfg.Section(sectionName)
	section.NewKey(key, value)
}

func (gitConfig *GitConfig) Save() {
	gitConfig.cfg.SaveTo(gitConfig.configFile)
}

func (gitConfig *GitConfig) GetPrefix(name string) (prefix string, err error) {
	section := gitConfig.cfg.Section(gitConfig.sectionPrefix)

	key, err := section.GetKey(name)

	if err != nil {
		return
	}

	prefix = key.Value()

	return
}

func (gitConfig *GitConfig) GetBranch(name string) (prefix string, err error) {
	section := gitConfig.cfg.Section(gitConfig.sectionBranch)

	key, err := section.GetKey(name)

	if err != nil {
		return
	}

	prefix = key.Value()

	return
}
