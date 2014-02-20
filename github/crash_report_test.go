package github

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/github/hub/fixtures"
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

func TestSaveAlwaysReportOption(t *testing.T) {
	checkSavedReportCrashOption(t, true, "a", "always")
	checkSavedReportCrashOption(t, true, "always", "always")
}

func TestSaveNeverReportOption(t *testing.T) {
	checkSavedReportCrashOption(t, false, "e", "never")
	checkSavedReportCrashOption(t, false, "never", "never")
}

func TestDoesntSaveYesReportOption(t *testing.T) {
	checkSavedReportCrashOption(t, false, "y", "")
	checkSavedReportCrashOption(t, false, "yes", "")
}

func TestDoesntSaveNoReportOption(t *testing.T) {
	checkSavedReportCrashOption(t, false, "n", "")
	checkSavedReportCrashOption(t, false, "no", "")
}

func checkSavedReportCrashOption(t *testing.T, always bool, confirm, expected string) {
	repo := fixtures.SetupTestRepo()
	defer repo.TearDown()

	saveReportConfiguration(confirm, always)
	assert.Equal(t, expected, reportCrashConfig())
}
