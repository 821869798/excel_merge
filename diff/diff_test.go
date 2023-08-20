package diff

import (
	"excel_merge/util"
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

	result := util.FormatFieldName("File {file} had error {error}", "file", file, "error", err)
	t.Logf(strings.Join(result, "|"))

	t.Logf(util.GetFileNameWithoutExt("qwe.txt"))
}
