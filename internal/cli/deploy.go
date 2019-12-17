package main

import (
	"github.com/spf13/cobra"
)

// Deploy should initiate a connection via ssh to the deployment server
// download the binary from github and then run that binary with appropriate settings

var cmdDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "Controls the release process for an application",
	Long: ` The release command uses semantic versioning to create a new github release.
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  runDeployCmd,
}

// First we need to open a file where user can set the semver, changelog contents,
// then we can insert that changelog contents into the source files before we call make to build
func runDeployCmd(cmd *cobra.Command, args []string) {

}

func init() {
	var skipChangelog bool
	cmdDeploy.Flags().BoolVar(&skipChangelog, "skipChangelog", false,
		"don't add a changelog for this release")

	rootCmd.AddCommand(cmdDeploy)
}
