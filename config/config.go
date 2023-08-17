package config

import (
	"excel_merge/util"
	"github.com/BurntSushi/toml"
)

type RawGlobalConfig struct {
	DiffOutputType  string `toml:"diff_output"`
	MergeOutputType string `toml:"merge_output"`
	CompareTools    string `toml:"compare_exe"`
	DiffArgs        string `toml:"diff_arg"`
	MergeArgs       string `toml:"merge_arg"`
}

var GlobalConfig *RawGlobalConfig

func ParseConfig(configFile string) error {
	configFile = util.AbsOrRelExecutePath(configFile)
	GlobalConfig = new(RawGlobalConfig)
	if _, err := toml.DecodeFile(configFile, GlobalConfig); err != nil {
		return err
	}
	return nil
}
