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
	// rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "", "aws profile to use")
	// rootCmd.PersistentFlags().StringVar(&awsRegion, "region", "us-west-2", "aws region to use")
}

func run(cmd *cobra.Command, args []string) {
	logInit()

	logrus.Info("Hi! Starting yadr server...")

	r, err := registry.New("/hacking/data")
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

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

func check(e error) {
	if e != nil {
		logrus.Panic(e)
	}
}