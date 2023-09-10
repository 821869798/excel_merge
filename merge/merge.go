package merge

import (
	"errors"
	"fmt"
	"github.com/821869798/excel_merge/config"
	"github.com/821869798/excel_merge/convert"
	"github.com/821869798/excel_merge/define"
	"github.com/821869798/excel_merge/source"
	"github.com/821869798/fankit/console"
	"github.com/821869798/fankit/fanopen"
	"github.com/821869798/fankit/fanpath"
	"github.com/821869798/fankit/fanstr"
	"github.com/AlecAivazis/survey/v2"
	_ "github.com/AlecAivazis/survey/v2"
	"github.com/gookit/slog"
	"github.com/xuri/excelize/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func Run(fileList []string) {
	outputFilePath := fileList[3]
	originBaseFile := fileList[0]
	originRemoteFile := fileList[1]
	originLocalFile := fileList[2]

	if !define.IsExcelFile(outputFilePath) {
		// 非Excel文件，直接调用对比工具
		slog.Infof("[merge]The merge file is not an Excel file, start comparison tools directly.")
		diffArg := fanstr.FormatFieldName(config.GlobalConfig.MergeArgs, "base", originBaseFile, "remote", originRemoteFile, "local", originLocalFile, "output", outputFilePath)
		cmd := exec.Command(fanpath.AbsOrRelExecutePath(config.GlobalConfig.CompareTools), diffArg...)
		output, err := cmd.CombinedOutput()
		if nil != err {
			slog.Panicf("[merge]execute compare tool output:%s\nerror:%v", output, err)
		}
		return
	}

	err := fanpath.CreateDirIfNoExist(fanpath.RelExecuteDir(define.WorkMergeTempDir))
	if err != nil {
		slog.Panicf("[merge]Back local file error: %v", err)
	}

	if fanpath.ExistFile(outputFilePath) {
		backupFile := fanpath.RelExecuteDir(define.WorkMergeTempDir, filepath.Base(outputFilePath))
		err = fanpath.CopyFile(outputFilePath, backupFile)
		if err != nil {
			slog.Panicf("[merge]Back local file copy error: %v", err)
		}
		slog.Infof("[merge]Backup local excel file to %s", backupFile)
	}

	// 转换csv
	baseFile := convertFile(originBaseFile)
	remoteFile := convertFile(originRemoteFile)
	localFile := convertFile(originLocalFile)

	defer os.Remove(baseFile)
	defer os.Remove(remoteFile)
	defer os.Remove(localFile)

	fileName := filepath.Base(outputFilePath)
	tmpOutputFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + config.GlobalConfig.MergeOutputType
	tmpOutputFile := fanpath.RelExecuteDir(define.WorkMergeTempDir, tmpOutputFileName)

	if fanpath.ExistFile(tmpOutputFile) {
		err = os.Remove(tmpOutputFile)
		if err != nil {
			slog.Panicf("[merge] remove tmp output file error: %v", err)
		}
	}

	diffArg := fanstr.FormatFieldName(config.GlobalConfig.MergeArgs, "base", baseFile, "remote", remoteFile, "local", localFile, "output", tmpOutputFile)
	cmd := exec.Command(fanpath.AbsOrRelExecutePath(config.GlobalConfig.CompareTools), diffArg...)
	output, err := cmd.CombinedOutput()
	if nil != err {
		slog.Panicf("[merge]execute compare tool output:%s\nerror:%v", output, err)
	}
	slog.Infof(string(output))

	selectNumber := selectMergeBase()
	var mergeExcelFiles []string
	switch selectNumber {
	case 1:
		mergeExcelFiles = []string{originRemoteFile, originLocalFile}
	case 2:
		mergeExcelFiles = []string{originLocalFile, originRemoteFile}
	}

	err = mergeToExcel(tmpOutputFile, mergeExcelFiles, outputFilePath)
	if err != nil {
		slog.Panicf("[merge] excel mode[%v] error: %v", config.GlobalConfig.MergeOutputType, err)
		return
	}
	slog.Infof("[merge] excel file complete:%s", outputFilePath)
	console.AnyKeyToQuit()
}

func convertFile(file string) string {
	excelData, err := source.ReadExcel(file, true)
	if err != nil {
		slog.Panicf("[merge] Read excel error: %v", err)
		return ""
	}

	fileName := filepath.Base(file)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	timestamp := time.Now().Unix()
	outputFileName := fmt.Sprintf("%s-%d.%s", fileNameWithoutExt, timestamp, config.GlobalConfig.MergeOutputType)
	outputFile := filepath.Join(os.TempDir(), define.WorkGenCSVDir, outputFileName)

	err = fanpath.CreateDirIfNoExist(filepath.Dir(outputFile))
	if err != nil {
		slog.Panicf("%v", err)
		return ""
	}

	err = convert.RunConvert(config.GlobalConfig.MergeOutputType, excelData, outputFile)
	if err != nil {
		slog.Panicf("[merge] Convert excel mode[%v] error: %v", config.GlobalConfig.MergeOutputType, err)
		return ""
	}

	return outputFile
}

func selectMergeBase() int {
	var answer survey.OptionAnswer
	prompt := &survey.Select{
		Message: `Select base excel file to merge.
The data is merged, but the formatting in the excel file of your choice is preferred.
数据是合并后的，但是优先保留你选择的excel文件中的格式`,
		Options: []string{"1. remote (基于远程分支excel结构)", "2. local (基于本地分支excel结构)"},
	}
	err := survey.AskOne(prompt, &answer)
	if err != nil {
		slog.Panicf("[merge] select merge base excel file error: %v", err)
	}
	return answer.Index + 1
}

// mergeToExcel 将csvFile写入到excel中
func mergeToExcel(csvFilePath string, mergeExcelFiles []string, outputExcelFilePath string) error {

	// 读取csvFile,然后读取excel文件sourceExcelFile，然后写入csv数据到excel中，最后输出到outputExcelFile
	excelData, err := convert.RunReadToExcel(config.GlobalConfig.MergeOutputType, csvFilePath)
	if err != nil {
		return err
	}

	var excelFiles []*excelize.File
	for _, excelFile := range mergeExcelFiles {
		excel, err := excelize.OpenFile(excelFile)
		if err != nil {
			return err
		}
		excelFiles = append(excelFiles, excel)
	}

	// sheet name 对应哪个excel文件的索引，使用最优先的情况
	sheetMapping := make(map[string]int)

	// 从后往前遍历，优先使用前面的excel文件的sheet
	for i := len(excelFiles) - 1; i >= 0; i-- {
		excel := excelFiles[i]
		for _, sheetName := range excel.GetSheetMap() {
			sheetMapping[sheetName] = i
		}
	}

	// 基础excel文件的sheet
	const baseExcelFileIndex = 0
	baseExcelFile := excelFiles[baseExcelFileIndex]

	// 删除base excel 多余的Sheet
	for _, sheetName := range baseExcelFile.GetSheetMap() {
		if _, ok := excelData.SheetMapping[sheetName]; !ok {
			err = baseExcelFile.DeleteSheet(sheetName)
			if err != nil {
				return err
			}
		}
	}

	// 对应基础的excel文件中的sheet索引顺序，可能有修改顺序
	baseExcelSheetsIndex := ExcelGetSheetIndexMap(baseExcelFile)

	for sheetIndex, sheet := range excelData.Sheets {
		excelIndex, ok := sheetMapping[sheet.SheetName]
		if !ok {
			return errors.New(fmt.Sprintf("Not found sheet name: %s", sheet.SheetName))
		}
		if excelIndex != baseExcelFileIndex {
			// 在当前位置拆入一张新sheet，但是没有样式  // TODO 等待excelize支持从别的excel文件中复制sheet
			err := ExcelInsertSheet(baseExcelFile, sheetIndex, sheet.SheetName)
			if err != nil {
				return err
			}

			// 重新获取Sheet顺序
			baseExcelSheetsIndex = ExcelGetSheetIndexMap(baseExcelFile)

		} else {
			// 验证Sheet顺序和csv文本中的是否一致,如果不一致则调整
			sourceSheetIndex := baseExcelSheetsIndex[sheet.SheetName]
			if sourceSheetIndex != sheetIndex {
				adjustSheetName := baseExcelFile.GetSheetName(sheetIndex)
				err = ExcelSwapSheetByName(baseExcelFile, sheet.SheetName, adjustSheetName)
				if err != nil {
					return err
				}

				// 重新获取Sheet顺序
				baseExcelSheetsIndex = ExcelGetSheetIndexMap(baseExcelFile)
			}
		}

		excelRows, err := baseExcelFile.GetRows(sheet.SheetName)
		if err != nil {
			return err
		}

		// csv数据设置到excel中
		for rowIndex, rowRecord := range sheet.RawData {
			axisStr, _ := excelize.CoordinatesToCellName(1, rowIndex+1)
			err = baseExcelFile.SetSheetRow(sheet.SheetName, axisStr, &rowRecord)
			if err != nil {
				return err
			}
			// 清除多余的列数据
			if rowIndex < len(excelRows) {
				for colIndex := len(rowRecord); colIndex < len(excelRows[rowIndex]); colIndex++ {
					axisStr, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
					err = baseExcelFile.SetCellValue(sheet.SheetName, axisStr, "")
					if err != nil {
						return err
					}
				}
			}
		}

		// 清除多余行数据
		for rowIndex := len(sheet.RawData); rowIndex < len(excelRows); rowIndex++ {
			for colIndex, _ := range excelRows[rowIndex] {
				axisStr, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
				err = baseExcelFile.SetCellValue(sheet.SheetName, axisStr, "")
				if err != nil {
					return err
				}
			}
		}

		if err != nil {
			return err
		}
	}

	err = baseExcelFile.SaveAs(outputExcelFilePath)
	return err
}

func ViewHistoryPath() {
	mergePath := fanpath.RelExecuteDir(define.WorkMergeTempDir)
	if !fanpath.ExistPath(mergePath) {
		slog.Panicf("[merge] Not found merge history path: %s", mergePath)
		return
	}
	err := fanopen.Start(mergePath)
	if err != nil {
		slog.Panicf("[merge] Open merge history path error: %v", err)
	}
}
