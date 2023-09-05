package scm

import (
	"github.com/821869798/excel_merge/define"
)

type IRegisterSCM interface {
	Name() string
	SupportSystem() define.SystemType
	Register(toolPath string) bool
}
