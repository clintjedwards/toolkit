package utils

import (
	"context"
	"os/exec"
	"time"
)

// ExecuteBashCmd takes in a command string and executes it using bash
// returns both stdout and stderr to user as a combined blob
func ExecuteBashCmd(command string, env []string, workingDir string) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Env = env
	cmd.Dir = workingDir

	// Execute command and return combined output to user
	return cmd.CombinedOutput()
}
