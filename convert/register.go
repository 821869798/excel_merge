package convert

import (
	"excel_merge/define"
)

var convertModes map[string]define.IConvert

func init() {
	convertModes = make(map[string]define.IConvert)
	convertModes["csv"] = NewConvertCSV()
	convertModes["txt"] = NewConvertCSV()
}
