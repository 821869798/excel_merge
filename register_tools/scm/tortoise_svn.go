package scm

import (
	"github.com/821869798/excel_merge/define"
	"github.com/821869798/excel_merge/register_tools/reg"
	"golang.org/x/sys/windows/registry"
)

type RegisterTortoiseSVN struct {
}

func NewRegisterTortoiseSVN() *RegisterTortoiseSVN {
	return &RegisterTortoiseSVN{}
}

func (r *RegisterTortoiseSVN) Name() string {
	return "TortoiseSVN"
}

func (r *RegisterTortoiseSVN) SupportSystem() define.SystemType {
	return define.SystemTypeWindows
}

func (r *RegisterTortoiseSVN) Register(toolPath string) bool {
	err := reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseSVN`, "Diff", toolPath)
	if err != nil {
		return false
	}
	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseSVN`, "Merge", toolPath)
	if err != nil {
		return false
	}
	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseSVN\DiffTools`, ".xlsx", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseSVN\DiffTools`, ".xlsm", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseSVN\DiffTools`, ".xls", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseSVN\MergeTools`, ".xlsx", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseSVN\MergeTools`, ".xlsm", toolPath)
	if err != nil {
		return false
	}

	err = reg.WriteRegistry(registry.CURRENT_USER, `SOFTWARE\TortoiseSVN\MergeTools`, ".xls", toolPath)
	if err != nil {
		return false
	}

	return true
}
