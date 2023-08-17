package define

type IConvert interface {
	Output(excelData *ExcelData, filePath string) error
	ReadToExcel(csvFilePath string) (*ExcelData, error)
}
