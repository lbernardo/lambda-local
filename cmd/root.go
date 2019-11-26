package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lam",
	Short: "Execute functions",
	Long:  `Execute functions lambda aws`,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
