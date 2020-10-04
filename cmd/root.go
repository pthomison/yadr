package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/pthomison/yadr/registry"
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
	fmt.Println("hi")

	r, err := registry.New("/hacking/data")
	check(err)

	r.Serve()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}