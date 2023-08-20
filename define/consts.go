package define

import (
	"path/filepath"
	"strings"
)

var (
	SheetHeadFormat  = "===[%s]==="
	SheetHeadStart   = "===["
	SheetHeadEnd     = "]==="
	WorkGenCSVDir    = "excel_merge_csv"
	WorkMergeTempDir = "excel_merge_temp"

	excelFileExt = []string{".xlsx", ".xlsm", ".xls"}

	// SafeExitCodeMapping 对比工具安全的退出码
	SafeExitCodeMapping = map[string]map[int]bool{
		"BCompare.exe": {
			0:  true,
			1:  true,
			2:  true,
			11: true,
			12: true,
			13: true,
		},
	}
)

// IsExcelFile 判断是否是excel文件
func IsExcelFile(fileName string) bool {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	// 判断是否是excel文件
	for _, ext := range excelFileExt {
		if fileExt == ext {
			return true
		}
	}
	return false
}

// IsCompareExitCodeSafe 判断对比工具的退出码是否安全
func IsCompareExitCodeSafe(filePath string, exitCode int) bool {
	compareTool := filepath.Base(filePath)
	if _, ok := SafeExitCodeMapping[compareTool]; ok {
		if _, ok := SafeExitCodeMapping[compareTool][exitCode]; ok {
			return true
		}
	}
	return false
}
