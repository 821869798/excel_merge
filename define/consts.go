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
