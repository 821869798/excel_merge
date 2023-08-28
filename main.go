package main

import (
	"excel_merge/config"
	"excel_merge/diff"
	"excel_merge/merge"
	"excel_merge/util"
	"flag"
	"fmt"
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
			util.AnyKeyToQuit()
			os.Exit(1)
		}
	}()

	err := config.ParseConfig(*conf)
	if err != nil {
		slog.Panicf("Load config toml file failed: %v", err)
	} else {
		slog.Infof("Load config toml success")
	}

	fileList := flag.Args()
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
