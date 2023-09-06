package main

import (
	"flag"
	"fmt"
	"github.com/821869798/excel_merge/config"
	"github.com/821869798/excel_merge/diff"
	"github.com/821869798/excel_merge/merge"
	"github.com/821869798/excel_merge/register_tools"
	"github.com/821869798/fankit/console"
	"github.com/821869798/fankit/fanpath"
	"github.com/gookit/slog"
	"os"
)

var (
	conf = flag.String("config", "config.toml", "runtime config path")
	mode = flag.String("mode", "", "execute mode:[diff,merge]")
	help = flag.Bool("help", false, "help with the usage")
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: excel_merge -mode <mode> file1 file2 ...
	diff example: excel_merge -mode diff left_file right_file
	merge example: excel_merge -mode merge base_file remote_file local_file output_file1
`)
	flag.PrintDefaults()
}

func main() {
	slog.SetLogLevel(slog.InfoLevel)

	flag.Parse()

	if *help {
		usage()
	}

	//slog.Infof(strings.Join(flag.Args(), "|"))

	defer func() {
		if err := recover(); err != nil {
			slog.Errorf("[main] catch exception: %v", err)
			console.AnyKeyToQuit()
			os.Exit(1)
		}
	}()

	// 初始化执行目录
	err := fanpath.InitExecutePath()
	if err != nil {
		slog.Panicf("init execute path error:%v", err)
	}

	fileList := flag.Args()
	if len(fileList) == 0 {
		err = config.ParseConfig(*conf)
		if err != nil {
			// 写入一份新的
			err = config.WriteNewConfig(*conf)
			if err != nil {
				slog.Panicf("Write new config toml file failed: %v", err)
			}
			parseConfig()
		}

		register_tools.Run()
		return
	}

	// 加载配置文件
	parseConfig()

	modeString := *mode

	if modeString != "" {
		switch modeString {
		case "diff":
			diff.Run(fileList)
		case "merge":
			merge.Run(fileList)
		default:
			usage()
		}
		return
	}

	if len(fileList) == 2 {
		diff.Run(fileList)
	} else if len(fileList) == 4 {
		merge.Run(fileList)
	} else {
		usage()
	}
}

func parseConfig() {
	err := config.ParseConfig(*conf)
	if err != nil {
		slog.Panicf("Load config toml file failed: %v", err)
	} else {
		slog.Infof("Load config toml success")
	}
}
