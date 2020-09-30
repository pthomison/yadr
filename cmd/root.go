package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"github.com/pthomison/yadr/registry"
)

var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: "",
	Long:  ``,
	Run: run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "", "aws profile to use")
	// rootCmd.PersistentFlags().StringVar(&awsRegion, "region", "us-west-2", "aws region to use")
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("hi")
	registry.Serve()
}