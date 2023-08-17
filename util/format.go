package util

import (
	"fmt"
	"strings"
)

// FormatByName FormatByName("File {file} had error {error}", "file", file, "error", err)
func FormatByName(format string, args ...interface{}) string {
	args2 := make([]string, len(args))
	for i, v := range args {
		if i%2 == 0 {
			args2[i] = fmt.Sprintf("{%v}", v)
		} else {
			args2[i] = fmt.Sprint(v)
		}
	}
	r := strings.NewReplacer(args2...)
	return r.Replace(format)
}

func FormatFieldName(format string, args ...string) []string {
	fields := strings.Fields(format)
	formatMapping := make(map[string]string)
	for i := 0; i < len(args); i += 2 {
		key := "{" + args[i] + "}"
		value := args[i+1]
		formatMapping[key] = value
	}
	for i, v := range fields {
		if value, ok := formatMapping[v]; ok {
			fields[i] = value
		}
	}

	return fields
}
