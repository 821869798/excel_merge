package convert

import (
	"errors"
	"fmt"
	"github.com/821869798/excel_merge/define"
)

func RunConvert(mode string, excelData *define.ExcelData, filePath string) error {
	c, ok := convertModes[mode]
	if !ok {
		return errors.New(fmt.Sprintf("Not support convert mode: %s", mode))
	}
	return c.Output(excelData, filePath)
}

func RunReadToExcel(mode string, csvFilePath string) (*define.ExcelData, error) {
	c, ok := convertModes[mode]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Not support merge mode: %s", mode))
	}
	return c.ReadToExcel(csvFilePath)
}
