package diff

import (
	"excel_merge/config"
	"excel_merge/convert"
	"excel_merge/define"
	"excel_merge/source"
	"excel_merge/util"
	"fmt"
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

	diffArg := util.FormatFieldName(config.GlobalConfig.DiffArgs, "left", file1, "right", file2)
	cmd := exec.Command(util.AbsOrRelExecutePath(config.GlobalConfig.CompareTools), diffArg...)
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
		slog.Panicf("Read Excel error: %v", err)
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
		slog.Panicf("Read Excel error: %v", err)
		return ""
	}

	err = util.CreateDirIfNoExist(filepath.Dir(outputFile))
	if err != nil {
		slog.Panicf("%v", err)
		return ""
	}

	err = convert.RunConvert(config.GlobalConfig.DiffOutputType, excelData, outputFile)
	if err != nil {
		slog.Panicf("Convert Excel to mode[%v] error: %v", config.GlobalConfig.DiffOutputType, err)
		return ""
	}

	return outputFile
}
