package compare

import (
	"github.com/821869798/excel_merge/define"
	"github.com/821869798/excel_merge/register_tools/reg"
	"strings"
)

type winBeyondCompare struct {
}

func (w *winBeyondCompare) Name() string {
	return "Beyond Compare"
}

func (w *winBeyondCompare) SupportSystem() define.SystemType {
	return define.SystemTypeWindows
}

func (w *winBeyondCompare) GetExecuteFilePath() (string, bool) {
	value, err := reg.ReadRegistry(reg.LOCAL_MACHINE, `SOFTWARE\Classes\BeyondCompare.SettingsPackage\shell\open\command`, "")
	if err == nil {
		path, ok := parseBCPath(value)
		if ok {
			return path, true
		}
	}
	return "", false
}

func parseBCPath(str string) (string, bool) {
	// 查找第一个双引号的位置
	firstQuoteIndex := strings.Index(str, "\"")
	if firstQuoteIndex == -1 {
		return "", false
	}

	// 从第一个双引号之后查找第二个双引号的位置
	secondQuoteIndex := strings.Index(str[firstQuoteIndex+1:], "\"")
	if secondQuoteIndex == -1 {
		return "", false
	}

	// 提取双引号中的内容
	content := str[firstQuoteIndex+1 : firstQuoteIndex+1+secondQuoteIndex]
	return content, true
}
