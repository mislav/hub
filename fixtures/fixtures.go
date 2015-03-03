package fixtures

import (
	"os"
	"path/filepath"
)

func Path(segment ...string) string {
	pwd, _ := os.Getwd()
	p := []string{pwd, "..", "fixtures"}
	p = append(p, segment...)

	return filepath.Join(p...)
}
