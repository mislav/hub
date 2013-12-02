package git

import (
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestGitDir(t *testing.T) {
	gitDir, _ := Dir()
	assert.T(t, strings.Contains(gitDir, ".git"))
}

func TestGitPullReqMsgFile(t *testing.T) {
	gitPullReqMsgFile, _ := PullReqMsgFile()
	assert.T(t, strings.Contains(gitPullReqMsgFile, "PULLREQ_EDITMSG"))
}

func TestGitEditor(t *testing.T) {
	gitEditor, err := Editor()
	if err == nil {
		assert.NotEqual(t, "", gitEditor)
	}
}

func TestGitRemote(t *testing.T) {
	gitRemote, _ := OriginRemote()
	assert.Equal(t, "origin", gitRemote.Name)
	assert.T(t, strings.Contains(gitRemote.URL, "gh"))
}

func TestGitHead(t *testing.T) {
	gitHead, _ := Head()
	assert.NotEqual(t, "", gitHead)
}

func TestGitLog(t *testing.T) {
	logs, _ := Log("master", "HEAD")
	assert.T(t, len(logs) >= 0)
}

func TestGitRef(t *testing.T) {
	gitRef, err := Ref("master")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", gitRef)
}

func TestGitRefList(t *testing.T) {
	refList, err := RefList("e357a98a1a580b09d4f1d9bf613a6a51e131ef6e", "49e984e2fe86f68c386aeb133b390d39e4264ec1")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(refList))

	assert.Equal(t, "49e984e2fe86f68c386aeb133b390d39e4264ec1", refList[0])
}

func TestGitShow(t *testing.T) {
	output, err := Show("ce20e63ad00751bfed5d08072b11cf1b43af1995")
	assert.Equal(t, nil, err)
	assert.Equal(t, "Add Git.RefList", output)
}
