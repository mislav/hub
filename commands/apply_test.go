package commands

import (
	"fmt"
	"github.com/bmizerany/assert"
	"regexp"
	"strings"
	"testing"
)

func TestTransformApplyArgs(t *testing.T) {
	args := NewArgs([]string{"apply", "https://github.com/jingweno/gh/pull/55"})
	transformApplyArgs(args)

	cmds := args.Commands()
	assert.Equal(t, 2, len(cmds))
	curlString := fmt.Sprintf("curl -#LA %s https://github.com/jingweno/gh/pull/55.patch -o .+/55.patch", fmt.Sprintf("gh %s", Version))
	curlRegexp := regexp.MustCompile(curlString)
	applyString := "git apply"
	assert.T(t, curlRegexp.MatchString(cmds[0].String()))
	assert.T(t, strings.Contains(cmds[1].String(), applyString))

	args = NewArgs([]string{"apply", "--ignore-whitespace", "https://github.com/jingweno/gh/commit/fdb9921"})
	transformApplyArgs(args)

	cmds = args.Commands()
	assert.Equal(t, 2, len(cmds))
	curlString = fmt.Sprintf("curl -#LA %s https://github.com/jingweno/gh/commit/fdb9921.patch -o .+/fdb9921.patch", fmt.Sprintf("gh %s", Version))
	curlRegexp = regexp.MustCompile(curlString)
	applyString = "git apply --ignore-whitespace"
	assert.T(t, curlRegexp.MatchString(cmds[0].String()))
	assert.T(t, strings.Contains(cmds[1].String(), applyString))

	args = NewArgs([]string{"apply", "https://gist.github.com/8da7fb575debd88c54cf"})
	transformApplyArgs(args)

	cmds = args.Commands()
	assert.Equal(t, 2, len(cmds))
	curlString = fmt.Sprintf("curl -#LA %s https://gist.github.com/8da7fb575debd88c54cf.txt -o .+8da7fb575debd88c54cf.txt", fmt.Sprintf("gh %s", Version))
	curlRegexp = regexp.MustCompile(curlString)
	applyString = "git apply"
	assert.T(t, curlRegexp.MatchString(cmds[0].String()))
	assert.T(t, strings.Contains(cmds[1].String(), applyString))
}
