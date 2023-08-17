package define

type ExcelData struct {
	fileName     string
	Sheets       []*ExcelSheet
	SheetMapping map[string]*ExcelSheet
}

func NewExcelData() *ExcelData {
	return &ExcelData{
		SheetMapping: make(map[string]*ExcelSheet),
	}
}

type ExcelSheet struct {
	SheetName string
	RawData   [][]string
}
