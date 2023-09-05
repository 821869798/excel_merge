package scm

import (
	"encoding/json"
	"github.com/821869798/excel_merge/define"
	"os"
	"os/exec"
	"path/filepath"
)

type RegisterFork struct {
}

func NewRegisterFork() *RegisterFork {
	return &RegisterFork{}
}

func (r *RegisterFork) Name() string {
	return "Fork"
}

func (r *RegisterFork) SupportSystem() define.SystemType {
	return define.SystemTypeWindows
}

func (r *RegisterFork) Register(toolPath string) bool {
	localPath, err := os.UserCacheDir()
	if err != nil {
		return false
	}

	err = killProcess("Fork.exe")
	if err != nil {
		return false
	}

	forkSettingPath := filepath.Join(localPath, "Fork", "settings.json")

	jsonData, err := os.ReadFile(forkSettingPath)
	if err != nil {
		return false
	}

	// 解码 JSON 数据为 map[string]interface{}
	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return false
	}

	// 修改配置
	data["MergeTool"] = map[string]interface{}{
		"Type":            "BeyondCompare",
		"ApplicationPath": toolPath,
	}

	data["ExternalDiffTool"] = map[string]interface{}{
		"Type":            "BeyondCompare",
		"ApplicationPath": toolPath,
	}
	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return false
	}

	// 将 JSON 字节切片写入回原始文件
	err = os.WriteFile(forkSettingPath, updatedJSON, 0755)
	if err != nil {
		return false
	}

	return true
}

func killProcess(name string) error {
	cmd := exec.Command("taskkill", "/IM", name, "/F")
	return cmd.Run()
}
