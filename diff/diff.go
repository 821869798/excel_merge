package diff

import (
	"fmt"
	"github.com/821869798/excel_merge/config"
	"github.com/821869798/excel_merge/convert"
	"github.com/821869798/excel_merge/define"
	"github.com/821869798/excel_merge/source"
	"github.com/821869798/fankit/fanpath"
	"github.com/821869798/fankit/fanstr"
	"github.com/gookit/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func Run(fileList []string) {
	// Excel文件会转成csv对比，非Excel会直接调用对比工具
	var file1 string = fileList[0]
	if define.IsExcelFile(file1) {
		file1 = convertFile(file1)
		defer os.Remove(file1)
	}
	var file2 string = fileList[1]
	if define.IsExcelFile(file2) {
		file2 = convertFile(file2)
		defer os.Remove(file2)
	}

	diffArg := fanstr.FormatFieldName(config.GlobalConfig.DiffArgs, "left", file1, "right", file2)
	cmd := exec.Command(fanpath.AbsOrRelExecutePath(config.GlobalConfig.CompareTools), diffArg...)
	output, err := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()
	if nil != err && !define.IsCompareExitCodeSafe(config.GlobalConfig.CompareTools, exitCode) {
		slog.Panicf("[diff]execute compare tool output:%s\nerror:%v", output, err)
		return
	}
	slog.Infof(string(output))
}

func convertFile(file string) string {
	fileInfo, err := os.Stat(file) // replace with your file
	if err != nil {
		slog.Panicf("[diff]Read Excel error: %v", err)
		return ""
	}
	fileName := filepath.Base(file)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	timestamp := time.Now().Unix()
	outputFileName := fmt.Sprintf("%s-%d.%s", fileNameWithoutExt, timestamp, config.GlobalConfig.DiffOutputType)
	outputFile := filepath.Join(os.TempDir(), define.WorkGenCSVDir, outputFileName)

	if fileInfo.Size() == 0 {
		_ = os.WriteFile(outputFile, []byte{}, 0755)
		return outputFile
	}

	excelData, err := source.ReadExcel(file, false)
	if err != nil {
		slog.Panicf("[diff]Read Excel error: %v", err)
		return ""
	}

	err = fanpath.CreateDirIfNoExist(filepath.Dir(outputFile))
	if err != nil {
		slog.Panicf("[diff] %v", err)
		return ""
	}

	err = convert.RunConvert(config.GlobalConfig.DiffOutputType, excelData, outputFile)
	if err != nil {
		slog.Panicf("[diff] Convert Excel to mode[%v] error: %v", config.GlobalConfig.DiffOutputType, err)
		return ""
	}

	return outputFile
}
