package scm

import (
	"github.com/821869798/excel_merge/define"
	"github.com/821869798/excel_merge/register_tools/reg"
	"golang.org/x/sys/windows/registry"
)

type RegisterTortoiseGit struct {
}

func NewRegisterTortoiseGit() *RegisterTortoiseGit {
	return &RegisterTortoiseGit{}
}

func (r *RegisterTortoiseGit) Name() string {
	return "TortoiseGit"
}

func (r *RegisterTortoiseGit) SupportSystem() define.SystemType {
	return define.SystemTypeWindows
}

func (r *RegisterTortoiseGit) Register(toolPath string) bool {
	err := reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseGit`, "Diff", toolPath)
	if err != nil {
		return false
	}
	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseGit`, "Merge", toolPath)
	if err != nil {
		return false
	}
	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseGit\DiffTools`, ".xlsx", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseGit\DiffTools`, ".xlsm", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseGit\DiffTools`, ".xls", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseGit\MergeTools`, ".xlsx", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseGit\MergeTools`, ".xlsm", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseGit\MergeTools`, ".xls", toolPath)
	if err != nil {
		return false
	}

	return true
}
