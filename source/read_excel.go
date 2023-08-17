package source

import (
	"excel_merge/define"
	"github.com/xuri/excelize/v2"
)

// ReadExcel 读取excel文件 consistentCol表示是否需要保持列数一致
func ReadExcel(file string, consistentCol bool) (outData *define.ExcelData, resultError error) {
	f, resultError := excelize.OpenFile(file)
	if resultError != nil {
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			resultError = err
		}
	}()

	outData = define.NewExcelData()

	for _, name := range f.GetSheetList() {

		rawData, err := f.GetRows(name)
		if err != nil {
			resultError = err
			return
		}

		if consistentCol {
			maxCells := 0
			for _, row := range rawData {
				if len(row) > maxCells {
					maxCells = len(row)
				}
			}
			for ri, row := range rawData {
				if len(row) < maxCells {
					for i := len(row); i < maxCells; i++ {
						row = append(row, "")
					}
				}
				rawData[ri] = row
			}
		}

		sheetData := new(define.ExcelSheet)
		sheetData.SheetName = name
		sheetData.RawData = rawData
		outData.Sheets = append(outData.Sheets, sheetData)
		outData.SheetMapping[sheetData.SheetName] = sheetData
	}
	return
}
