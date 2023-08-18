package merge

import (
	"bufio"
	"errors"
	"excel_merge/config"
	"excel_merge/convert"
	"excel_merge/define"
	"excel_merge/source"
	"excel_merge/util"
	"fmt"
	"github.com/gookit/slog"
	"github.com/xuri/excelize/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
		slog.Infof("The merge file is not an Excel file, start comparison tools directly.")
		diffArg := util.FormatFieldName(config.GlobalConfig.MergeArgs, "base", originBaseFile, "remote", originRemoteFile, "local", originLocalFile, "output", outputFilePath)
		cmd := exec.Command(util.AbsOrRelExecutePath(config.GlobalConfig.CompareTools), diffArg...)
		output, err := cmd.CombinedOutput()
		if nil != err {
			slog.Panicf("[diff]execute compare tool output:%s\nerror:%v", output, err)
		}
		return
	}

	err := util.CreateDirIfNoExist(util.RelExecuteDir(define.WorkMergeTempDir))
	if err != nil {
		slog.Panicf("Back local file error: %v", err)
	}

	if util.ExistFile(outputFilePath) {
		backupFile := util.RelExecuteDir(define.WorkMergeTempDir, filepath.Base(outputFilePath))
		err = util.CopyFile(outputFilePath, backupFile)
		if err != nil {
			slog.Panicf("Back local file copy error: %v", err)
		}
		slog.Infof("Backup local excel file to %s", backupFile)
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
	tmpOutputFile := util.RelExecuteDir(define.WorkMergeTempDir, tmpOutputFileName)

	diffArg := util.FormatFieldName(config.GlobalConfig.MergeArgs, "base", baseFile, "remote", remoteFile, "local", localFile, "output", tmpOutputFile)
	cmd := exec.Command(util.AbsOrRelExecutePath(config.GlobalConfig.CompareTools), diffArg...)
	output, err := cmd.CombinedOutput()
	if nil != err {
		slog.Panicf("[diff]execute compare tool output:%s\nerror:%v", output, err)
	}
	slog.Infof(string(output))

	selectNumber := selectBaseFile()
	var mergeExcelFiles []string
	switch selectNumber {
	case 1:
		mergeExcelFiles = []string{originRemoteFile, originLocalFile}
	case 2:
		mergeExcelFiles = []string{originLocalFile, originRemoteFile}
	}

	err = mergeToExcel(tmpOutputFile, mergeExcelFiles, outputFilePath)
	if err != nil {
		slog.Panicf("merge excel mode[%v] error: %v", config.GlobalConfig.MergeOutputType, err)
		return
	}
	slog.Infof("Merge excel file complete:%s", outputFilePath)
	util.AnyKeyToQuit()
}

func convertFile(file string) string {
	excelData, err := source.ReadExcel(file, true)
	if err != nil {
		slog.Panicf("Read excel error: %v", err)
		return ""
	}

	fileName := filepath.Base(file)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	timestamp := time.Now().Unix()
	outputFileName := fmt.Sprintf("%s-%d.%s", fileNameWithoutExt, timestamp, config.GlobalConfig.MergeOutputType)
	outputFile := filepath.Join(os.TempDir(), define.WorkGenCSVDir, outputFileName)

	err = util.CreateDirIfNoExist(filepath.Dir(outputFile))
	if err != nil {
		slog.Panicf("%v", err)
		return ""
	}

	err = convert.RunConvert(config.GlobalConfig.MergeOutputType, excelData, outputFile)
	if err != nil {
		slog.Panicf("Convert excel mode[%v] error: %v", config.GlobalConfig.MergeOutputType, err)
		return ""
	}

	return outputFile
}

func selectBaseFile() int {
	fmt.Printf(`Select base excel file to merge.
The data is merged, but the formatting in the excel file of your choice is preferred.
数据是合并后的，但是优先保留你选择的excel文件中的格式。
1. remote (基于远程分支excel结构)
2. local (基于本地分支excel结构)`)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\nPlease enter your selection number(请输入你的选择数字):")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("input error:%v", err)
		}
		number, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
		if err != nil {
			fmt.Printf("input error:%v", err)
		}
		if number > 2 || number < 1 {
			fmt.Printf("input error number error!")
		} else {
			return int(number)
		}
	}
}

// mergeToExcel 将csvFile写入到excel中
func mergeToExcel(csvFilePath string, mergeExcelFiles []string, outputExcelFilePath string) error {

	// 读取csvFile,然后读取excel文件sourceExcelFile，然后写入csv数据到excel中，最后输出到outputExcelFile
	excelData, err := convert.RunReadToExcel(config.GlobalConfig.MergeOutputType, csvFilePath)
	if err != nil {
		return err
	}

	var excelFiles []*excelize.File
	// sheet name 对应哪个excel文件的索引，使用最优先的情况
	sheetMapping := make(map[string]int)

	for _, excelFile := range mergeExcelFiles {
		excel, err := excelize.OpenFile(excelFile)
		if err != nil {
			return err
		}
		excelFiles = append(excelFiles, excel)
	}

	// 从后往前遍历，优先使用前面的excel文件的sheet
	for i := len(excelFiles) - 1; i >= 0; i-- {
		excel := excelFiles[i]
		for _, sheetName := range excel.GetSheetMap() {
			if _, ok := sheetMapping[sheetName]; !ok {
				sheetMapping[sheetName] = i
			}
		}
	}

	const baseExcelFileIndex = 0

	// 基础excel文件的sheet
	baseExcelFile := excelFiles[baseExcelFileIndex]

	for _, sheet := range excelData.Sheets {
		excelIndex, ok := sheetMapping[sheet.SheetName]
		if !ok {
			return errors.New(fmt.Sprintf("Not found sheet name: %s", sheet.SheetName))
		}
		currentExcelSheet := excelFiles[excelIndex]
		if excelIndex != baseExcelFileIndex {
			// 拷贝当前sheet到baseExcelFile中
			newSheetIndex, err := baseExcelFile.NewSheet(sheet.SheetName)
			if err != nil {
				return err
			}
			_ = newSheetIndex
			// TODO 等待excelize支持从别的excel文件中复制sheet，以及支持调整sheet顺序，目前只能新增sheet到默认
		}

		excelRows, err := currentExcelSheet.GetRows(sheet.SheetName)
		if err != nil {
			return err
		}
		// 清除多余的行数据
		if len(sheet.RawData) <= 0 {
			if len(excelRows) > 0 {
				// clear sheet
				for rowIndex, cols := range excelRows {
					for colIndex, _ := range cols {
						axisStr, _ := excelize.CoordinatesToCellName(colIndex, rowIndex+1)
						err = baseExcelFile.SetCellValue(sheet.SheetName, axisStr, "")
						if err != nil {
							return err
						}
					}
				}
			}
			continue
		}

		// csv数据设置到excel中
		for rowIndex, rowRecord := range sheet.RawData {
			axisStr, _ := excelize.CoordinatesToCellName(1, rowIndex+1)
			err = baseExcelFile.SetSheetRow(sheet.SheetName, axisStr, &rowRecord)
			if err != nil {
				return err
			}
			// 清除多余的列数据
			for colIndex := len(rowRecord); colIndex < len(excelRows[rowIndex]); colIndex++ {
				axisStr, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
				err = baseExcelFile.SetCellValue(sheet.SheetName, axisStr, "")
				if err != nil {
					return err
				}
			}
		}

		// 清除多余行数据
		for rowIndex := len(sheet.RawData); rowIndex < len(excelRows); rowIndex++ {
			for colIndex, _ := range excelRows[rowIndex] {
				axisStr, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
				err = currentExcelSheet.SetCellValue(sheet.SheetName, axisStr, "")
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
