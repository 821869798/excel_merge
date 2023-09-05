package compare

import (
	"github.com/821869798/excel_merge/define"
)

type ICompareTools interface {
	Name() string
	SupportSystem() define.SystemType
	GetExecuteFilePath() (string, bool)
}

var (
	compareToolsRegister []ICompareTools
)

func init() {
	// TODO 支持更多工具
	compareToolsRegister = []ICompareTools{
		&winBeyondCompare{},
	}
}

func SupportCompareTools(systemType define.SystemType) []ICompareTools {
	var compares []ICompareTools
	for _, c := range compareToolsRegister {
		if c.SupportSystem() == systemType {
			compares = append(compares, c)
		}
	}
	return compares
}
