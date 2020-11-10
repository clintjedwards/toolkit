package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for the basecoat cli
var rootCmd = &cobra.Command{
	Use:   "toolkit",
	Short: "Helper for simple releases",
}

const defaultConfigPath = "./.toolkit.yml"

func main() {
	rootCmd.PersistentFlags().Bool("hideOutput", false, "Hide output from commands")
	rootCmd.PersistentFlags().Bool("echoCommands", false, "Print commands before running")
	rootCmd.PersistentFlags().StringP("config", "c", ".toolkit.yml", "Path of toolkit config file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
