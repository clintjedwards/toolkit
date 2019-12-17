package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for the basecoat cli
var rootCmd = &cobra.Command{
	Use:   "tools",
	Short: "I got fed up of makefiles",
}

// executeCmd wraps a context around a given shell command and executes it.
// dir refers to the working directory of command to be run
func executeCmd(path string, args []string, env []string, workingDir string) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = env
	cmd.Dir = workingDir

	// Execute command
	return cmd.CombinedOutput()
}

// execute adds all child commands to the root command and sets flags appropriately.
func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
