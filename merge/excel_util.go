package merge

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

const (
	ExcelMergeSheetTemp = "excel_merge_sheet_temp"
)

func ExcelGetSheetIndexMap(excel *excelize.File) map[string]int {
	SheetsIndex := make(map[string]int, excel.SheetCount)
	for index, sheetName := range excel.GetSheetList() {
		SheetsIndex[sheetName] = index
	}
	return SheetsIndex
}

func ExcelInsertSheet(excel *excelize.File, index int, sheetName string) error {
	if index > excel.SheetCount || index < 0 {
		return errors.New(fmt.Sprintf("index is out of range:%v", index))
	}
	sheetNameList := excel.GetSheetList()
	for _, name := range sheetNameList {
		if name == sheetName {
			return errors.New(fmt.Sprintf("sheet name is exist:%v", sheetName))
		}
	}
	if index == excel.SheetCount {
		_, err := excel.NewSheet(sheetName)
		return err
	}

	lastIndex, err := excel.NewSheet(sheetName)
	if err != nil {
		return err
	}
	for i := lastIndex - 1; i >= index; i-- {
		err = excel.CopySheet(i, i+1)
		if err != nil {
			return err
		}
		sourceName := excel.GetSheetName(i)
		nextName := excel.GetSheetName(i + 1)
		err = excel.SetSheetName(sourceName, ExcelMergeSheetTemp)
		if err != nil {
			return err
		}
		err = excel.SetSheetName(nextName, sourceName)
		if err != nil {
			return err
		}
		err = excel.SetSheetName(ExcelMergeSheetTemp, nextName)
		if err != nil {
			return err
		}
	}
	newSheetIndex, err := excel.NewSheet(ExcelMergeSheetTemp)
	if err != nil {
		return err
	}
	err = excel.CopySheet(newSheetIndex, index)
	if err != nil {
		return nil
	}
	err = excel.DeleteSheet(ExcelMergeSheetTemp)
	if err != nil {
		return err
	}
	return nil
}

func ExcelSwapSheetByName(excel *excelize.File, srcName, destName string) error {
	if srcName == destName {
		return errors.New(fmt.Sprintf("srcName is equal destName: %v", srcName))
	}

	nameCount := 0
	sheetNameList := excel.GetSheetList()
	for _, name := range sheetNameList {
		if srcName == name {
			nameCount++
		} else if destName == name {
			nameCount++
		}
	}
	if nameCount != 2 {
		return errors.New(fmt.Sprintf("srcName or destName is not exist: %v %v", srcName, destName))
	}
	// 交换两个sheet的位置
	srcIndex, err := excel.GetSheetIndex(srcName)
	if err != nil {
		return err
	}
	destIndex, err := excel.GetSheetIndex(destName)
	if err != nil {
		return err
	}

	newSheetIndex, err := excel.NewSheet(ExcelMergeSheetTemp)
	if err != nil {
		return err
	}
	err = excel.CopySheet(srcIndex, newSheetIndex)
	if err != nil {
		return err
	}

	err = excel.CopySheet(destIndex, srcIndex)
	if err != nil {
		return err
	}

	err = excel.CopySheet(newSheetIndex, destIndex)
	if err != nil {
		return err
	}

	err = excel.DeleteSheet(ExcelMergeSheetTemp)
	if err != nil {
		return err
	}

	err = excel.SetSheetName(srcName, ExcelMergeSheetTemp)
	if err != nil {
		return err
	}
	err = excel.SetSheetName(destName, srcName)
	if err != nil {
		return err
	}
	err = excel.SetSheetName(ExcelMergeSheetTemp, destName)
	if err != nil {
		return err
	}

	return nil
}
