package register_tools

import (
	"errors"
	"fmt"
	"github.com/821869798/excel_merge/config"
	"github.com/821869798/excel_merge/define"
	"github.com/821869798/excel_merge/register_tools/compare"
	"github.com/821869798/excel_merge/register_tools/scm"
	"github.com/821869798/fankit/admin"
	"github.com/821869798/fankit/console"
	"github.com/821869798/fankit/fanpath"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/slog"
	"github.com/ncruces/zenity"
	"os"
	"runtime"
)

var (
	registerSCMs []scm.IRegisterSCM = make([]scm.IRegisterSCM, 0)
)

func init() {
	registerSCMs = append(registerSCMs, scm.NewRegisterTortoiseGit())
	registerSCMs = append(registerSCMs, scm.NewRegisterTortoiseSVN())
	registerSCMs = append(registerSCMs, scm.NewRegisterFork())
}

func Run() {
	var currentSystem define.SystemType = define.SystemTypeNone
	switch runtime.GOOS {
	case "windows":
		currentSystem = define.SystemTypeWindows
		if !admin.IsAdministrator() {
			err := admin.StartRunAdministrator(fanpath.ExecuteFilePath(), []string{})
			if err != nil {
				slog.Panicf("run as administrator error:%v", err)
			}
			os.Exit(200)
		}
	case "linux":
		currentSystem = define.SystemTypeLinux
	case "darwin":
		currentSystem = define.SystemTypeMac
	}
	// 判断当前支持注册的工具
	supportRegisters := make([]scm.IRegisterSCM, 0)
	for _, scmTools := range registerSCMs {
		if scmTools.SupportSystem() == currentSystem {
			supportRegisters = append(supportRegisters, scmTools)
		}
	}

	if len(supportRegisters) <= 0 {
		slog.Warnf("[register] The current system does not support registration tools. Please manually configure config.toml and manually configure version tools\n(当前系统不支持注册工具，请手动配置config.toml,以及手动配置版本工具)")
		console.AnyKeyToQuit()
		return
	}

	// 通过对话框选择对比工具
	filePath, err := SelectCompareTools(currentSystem)
	if err != nil {
		slog.Panicf("[register] select compare_tools error: %v", err)
		return
	}

	// 更新配置文件
	config.GlobalConfig.CompareTools = filePath
	err = config.UpdateCompareTool(filePath)
	if err != nil {
		slog.Panicf("[register] Unable to update config file(无法更新配置文件):%v", err)
		return
	}

	// 选择需要执行的注册工具
	var options = make([]string, len(supportRegisters))
	var checked = make(map[int]bool, len(supportRegisters))
	for index, scmTools := range supportRegisters {
		options[index] = scmTools.Name()
		checked[index] = true
	}

	// 选择注册的scm工具
	var answers = make([]survey.OptionAnswer, 0, len(supportRegisters))
	prompt := &survey.MultiSelect{
		Message: "Please select the SCM tool that needs to be registered(请选择需要注册的SCM工具,选择好后回车):",
		Options: options,
		Default: options,
	}
	err = survey.AskOne(prompt, &answers, survey.WithIcons(func(icons *survey.IconSet) {
		// you can set any icons
		icons.MarkedOption.Text = "[√]"
	}))
	if err != nil {
		slog.Panicf("[register] select register_tools scm error: %v", err)
		return
	}

	if len(answers) == 0 {
		slog.Warnf("[register] No SCM tool is registered(没有选择注册SCM工具)")
		console.AnyKeyToQuit()
		return
	}

	// 开始注册
	for _, ans := range answers {
		scmTools := supportRegisters[ans.Index]
		ok := scmTools.Register(fanpath.ExecuteFilePath())
		if ok {
			slog.Infof("[register] register_tools success(注册SCM工具成功): %s", scmTools.Name())
		} else {
			slog.Warnf("[register] Failed to register SCM tool, possibly not installed or open\n(注册SCM工具失败,可能是没有安装,或者是打开状态): %s", scmTools.Name())
		}
	}

	slog.Infof("Registration completed, please reopen the version tool(注册完毕，请重新打开版本工具)")

	console.AnyKeyToQuit()
}

// SelectCompareTools 选择对比工具
func SelectCompareTools(currentSystem define.SystemType) (string, error) {
	// 选择对比工具
	compareTools := compare.SupportCompareTools(currentSystem)
	var supportToolsPath []string
	var supportToolsName []string
	for _, ct := range compareTools {
		toolsPath, ok := ct.GetExecuteFilePath()
		if ok {
			supportToolsPath = append(supportToolsPath, toolsPath)
			supportToolsName = append(supportToolsName, ct.Name())
		}
	}
	if len(supportToolsName) <= 0 {
		// 直接选择自定义工具
		return SelectCustomCompareTools(currentSystem)
	}

	supportToolsName = append(supportToolsName, "Select Custom(选择自定义)")

	var answer survey.OptionAnswer
	prompt := &survey.Select{
		Message: `Please select a basic comparison tool(Please select a basic comparison tool):`,
		Options: supportToolsName,
	}
	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return "", errors.New(fmt.Sprintf("[register] select compare tools failed error: %v", err))
	}

	if answer.Index == len(supportToolsName)-1 {
		return SelectCustomCompareTools(currentSystem)
	}

	return supportToolsPath[answer.Index], nil
}

// SelectCustomCompareTools 通过对话框选择对比工具
func SelectCustomCompareTools(currentSystem define.SystemType) (string, error) {
	extensions := "*.*"
	if currentSystem == define.SystemTypeWindows {
		extensions = "*.exe"
	}

	filePath, err := zenity.SelectFile(
		zenity.Title("Please select a comparison tool(请选择对比工具，例如Beyond Compare)"),
		zenity.FileFilters{
			{fmt.Sprintf("execute file(%s)", extensions), []string{extensions}, true},
		})
	return filePath, err
}
