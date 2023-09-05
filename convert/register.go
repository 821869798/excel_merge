package convert

import (
	"github.com/821869798/excel_merge/define"
)

var convertModes map[string]define.IConvert

func init() {
	convertModes = make(map[string]define.IConvert)
	convertModes["csv"] = NewConvertCSV()
	convertModes["txt"] = NewConvertCSV()
}
