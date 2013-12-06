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
	assert.T(t, strings.Contains(gitRemote.URL.String(), "gh"))
}

func TestGitLog(t *testing.T) {
	log, err := Log("e357a98a1a580b09d4f1d9bf613a6a51e131ef6e", "49e984e2fe86f68c386aeb133b390d39e4264ec1")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", log)
}

func TestGitRef(t *testing.T) {
	gitRef, err := Ref("1c1077c052d32a83aa13a8afaa4a9630d2f28ef6")
	assert.Equal(t, nil, err)
	assert.Equal(t, "1c1077c052d32a83aa13a8afaa4a9630d2f28ef6", gitRef)
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
