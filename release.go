package main

import (
	"fmt"
	"log"

	"github.com/clintjedwards/toolkit/github"
	"github.com/spf13/cobra"
)

const binaryPathFmt string = "/tmp/%s_%s"

var cmdRelease = &cobra.Command{
	Use:   "release <semver>",
	Short: "Controls the release process for an application",
	Long: `The release command uses semantic versioning to build a new version
of the provided application and create a new github release.

tokenFile should contain nothing but the github token with access to repo
`,
	Args: cobra.MinimumNArgs(1),
	Run:  runReleaseCmd,
}

// First we need to open a file where user can set the semver, changelog contents,
// then we can insert that changelog contents into the source files before we call make to build
func runReleaseCmd(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("config")

	newRelease, err := github.NewRelease(configFile, args)
	if err != nil {
		log.Fatalf("could not create new release instance: %v", err)
	}

	var binaryPath string
	skipBinary, _ := cmd.Flags().GetBool("skipBinary")
	if !skipBinary {
		// set project build path so we have a predictable location
		binaryPath = fmt.Sprintf(binaryPathFmt, newRelease.ProjectName, newRelease.Version)
		runBuildCmd(cmd, []string{newRelease.Version, binaryPath})
	}

	tokenFile, _ := cmd.Flags().GetString("tokenFile")
	newRelease.CreateGithubRelease(tokenFile, binaryPath)
}

func init() {
	cmdRelease.Flags().Bool("skipBuild", false, "don't add a build asset for this release")
	cmdRelease.Flags().StringP("tokenFile", "t", "", "github api key file (default is $HOME/.github_token)")

	rootCmd.AddCommand(cmdRelease)
}
