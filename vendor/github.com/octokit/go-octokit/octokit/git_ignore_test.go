package octokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitIgnoreService_All(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/gitignore/templates", "git_ignore_templates", nil)

	templates, result := client.GitIgnore().All(nil)
	assert.False(t, result.HasError())
	assert.Equal(t, "AppceleratorTitanium", templates[2])
	assert.Equal(t, "Autotools", templates[3])
	assert.Len(t, templates, 7)
}

func TestGitIgnoreService_One(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/gitignore/templates/C", "git_ignore_c_template", nil)

	template, result := client.GitIgnore().One(&GitIgnoreURL, M{"name": "C"})
	assert.False(t, result.HasError())
	assert.Equal(t, "C", template.Name)
	assert.Equal(t, "# Object files\n*.o\n\n# Libraries\n*.lib\n*.a\n\n# Shared objects (inc. Windows DLLs)\n*.dll\n*.so\n*.so.*\n*.dylib\n\n# Executables\n*.exe\n*.out\n*.app\n", template.Source)
}
