package cmd

import (
	"fmt"
	"os"
	"runtime"
	"github.com/spf13/cobra"
	"github.com/pthomison/yadr/registry"
    "github.com/sirupsen/logrus"
    "path"
)

var(
	dataDirectory string
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: "",
	Long:  ``,
	Run: run,
}

func Execute() {
	check(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dataDirectory, "data-directory", "./data", "where to store registry data")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "INFO", "INFO/DEBUG/ERROR")
}

func run(cmd *cobra.Command, args []string) {
	logInit()

	logrus.Info("Hi! Starting yadr server...")

	r, err := registry.New(dataDirectory)
	check(err)

	r.Serve()
}

func logInit() {
    logrus.SetFormatter(&logrus.TextFormatter{
	    CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
	        function = ""
	        file = fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
	        return 
	    },
    })

    var level logrus.Level

	switch logLevel {
	case "DEBUG":
		level = logrus.DebugLevel
	case "INFO":
		level = logrus.InfoLevel
	case "ERROR":
		level = logrus.ErrorLevel
	default:
		level = logrus.DebugLevel
	}

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)
}

func check(e error) {
	if e != nil {
		logrus.Panic(e)
	}
}