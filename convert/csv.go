package convert

import (
	"excel_merge/define"
)

type ConvertCSV struct {
}

func NewConvertCSV() define.IConvert {
	c := &ConvertCSV{}
	return c
}

func (c *ConvertCSV) Output(excelData *define.ExcelData, filePath string) error {
	// 生成csv
	return writeExcelToCSV(excelData, filePath)
}

func (c *ConvertCSV) ReadToExcel(csvFilePath string) (*define.ExcelData, error) {
	excelData, err := readCSVToExcel(csvFilePath)

	return excelData, err
}
