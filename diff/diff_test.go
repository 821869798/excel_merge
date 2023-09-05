package diff

import (
	"github.com/821869798/fankit/fanpath"
	"github.com/821869798/fankit/fanstr"
	"os"
	"strings"
	"testing"
)

func TestFunc1(t *testing.T) {
	//os.CreateTemp()
	t.Logf(os.TempDir())
	var files []string = []string{"a", "b", "c"}
	var actions []func()
	for _, file := range files {
		actions = append(actions, func() {
			t.Logf(file)
		})
	}

	for _, action := range actions {
		action()
	}

	file := "my file.txt"
	err := "file not found"

	result := fanstr.FormatFieldName("File {file} had error {error}", "file", file, "error", err)
	t.Logf(strings.Join(result, "|"))

	t.Logf(fanpath.GetFileNameWithoutExt("qwe.txt"))
}
