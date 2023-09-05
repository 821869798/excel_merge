package config

import (
	"fmt"
	"github.com/821869798/fankit/fanpath"
	"github.com/BurntSushi/toml"
	"os"
	"regexp"
	"runtime"
	"strings"
)

type RawGlobalConfig struct {
	DiffOutputType  string `toml:"diff_output"`
	MergeOutputType string `toml:"merge_output"`
	CompareTools    string `toml:"compare_exe"`
	DiffArgs        string `toml:"diff_arg"`
	MergeArgs       string `toml:"merge_arg"`
}

var (
	GlobalConfig   *RawGlobalConfig
	ConfigFilePath string
)

func ParseConfig(configFile string) error {
	ConfigFilePath = fanpath.AbsOrRelExecutePath(configFile)
	GlobalConfig = new(RawGlobalConfig)
	if _, err := toml.DecodeFile(ConfigFilePath, GlobalConfig); err != nil {
		return err
	}
	return nil
}

func WriteConfig(configFile string) error {
	if configFile == "" {
		configFile = ConfigFilePath
	}
	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	if err := toml.NewEncoder(f).Encode(GlobalConfig); err != nil {
		return err
	}
	return nil
}

func UpdateCompareTool(newFilePath string) error {
	content, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		return err
	}

	// 定义正则表达式
	pattern := `compare_exe\s*=\s*"(.*?)"`
	regex := regexp.MustCompile(pattern)

	if runtime.GOOS == "windows" {
		// windows系统下，将路径中的反斜杠替换为双反斜杠
		newFilePath = strings.ReplaceAll(newFilePath, "\\", "\\\\")
	}
	// 使用正则表达式替换字符串
	newContent := regex.ReplaceAllString(string(content), fmt.Sprintf("compare_exe = \"%s\"", newFilePath))

	// 将替换后的内容写入文件
	err = os.WriteFile(ConfigFilePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}
	return nil

	// 选择需要注册的版本软件工具

}
