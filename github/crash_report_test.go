package github

import (
	"testing"

	"github.com/github/hub/v2/internal/assert"
)

func TestStackRemoveSelfAndPanic(t *testing.T) {
	actual := `goroutine 1 [running]:
runtime.panic(0x2bca00, 0x665b8a)
	/usr/local/go/src/pkg/runtime/panic.c:266 +0xb6
github.com/jingweno/gh/github.ReportCrash(0xc2000b5000, 0xc2000b49c0)
	/Users/calavera/github/go/src/github.com/jingweno/gh/github/crash_report.go:16 +0x97
github.com/jingweno/gh/commands.create(0x47f8a0, 0xc2000cf770)
	/Users/calavera/github/go/src/github.com/jingweno/gh/commands/create.go:54 +0x63
github.com/jingweno/gh/commands.(*Runner).Execute(0xc200094640, 0xc200094640, 0x21, 0xc2000b0a40)
	/Users/calavera/github/go/src/github.com/jingweno/gh/commands/runner.go:72 +0x3b7
main.main()
	/Users/calavera/github/go/src/github.com/jingweno/gh/main.go:10 +0xad`

	expected := `goroutine 1 [running]:
github.com/jingweno/gh/commands.create(0x47f8a0, 0xc2000cf770)
	/Users/calavera/github/go/src/github.com/jingweno/gh/commands/create.go:54 +0x63
github.com/jingweno/gh/commands.(*Runner).Execute(0xc200094640, 0xc200094640, 0x21, 0xc2000b0a40)
	/Users/calavera/github/go/src/github.com/jingweno/gh/commands/runner.go:72 +0x3b7
main.main()
	/Users/calavera/github/go/src/github.com/jingweno/gh/main.go:10 +0xad`

	s := formatStack([]byte(actual))
	assert.Equal(t, expected, s)
}
