package convert

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"excel_merge/define"
	"fmt"
	"io"
	"os"
	"strings"
)

// WriteExcelToCSV 将excelData写入到csv文件中
func writeExcelToCSV(excelData *define.ExcelData, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, sheet := range excelData.Sheets {
		err = writer.Write([]string{fmt.Sprintf(define.SheetHeadFormat, sheet.SheetName)})
		if err != nil {
			return err
		}
		for _, record := range sheet.RawData {
			if err = writer.Write(record); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReadCSVToExcel 读取csvFile到ExcelData中
func readCSVToExcel(csvFilePath string) (*define.ExcelData, error) {
	// 读取csvFile,然后读取excel文件sourceExcelFile，然后写入csv数据到excel中，最后输出到outputExcelFile
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	excelData := define.NewExcelData()
	var currentSheet *define.ExcelSheet
	var buffers []*bytes.Buffer
	var buf *bytes.Buffer

	br := bufio.NewReader(csvFile)
	for {
		byteDatas, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		content := string(byteDatas)
		if strings.HasPrefix(content, define.SheetHeadStart) && strings.HasSuffix(content, define.SheetHeadEnd) {
			sheetName := strings.TrimPrefix(content, define.SheetHeadStart)
			sheetName = strings.TrimSuffix(sheetName, define.SheetHeadEnd)
			sheetData := &define.ExcelSheet{
				SheetName: sheetName,
			}
			excelData.Sheets = append(excelData.Sheets, sheetData)
			excelData.SheetMapping[sheetData.SheetName] = sheetData
			currentSheet = sheetData

			buf = new(bytes.Buffer)
			buffers = append(buffers, buf)
			continue
		}
		if currentSheet == nil {
			return nil, fmt.Errorf("csv file format error,not a sheet name in begin")
		}
		buf.Write(byteDatas)
		buf.WriteByte('\n')
	}

	for index, sheet := range excelData.Sheets {
		sheet.RawData, err = csv.NewReader(buffers[index]).ReadAll()
		if err != nil {
			return nil, err
		}
	}

	return excelData, nil
}
