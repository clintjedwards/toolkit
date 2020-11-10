package main

import (
	"fmt"
	"os"
	"time"

	"github.com/clintjedwards/toolkit/changelog"
	"github.com/clintjedwards/toolkit/config"
	"github.com/clintjedwards/toolkit/github"
	"github.com/spf13/cobra"
	"github.com/theckman/yacspin"
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

func initSpinner(suffix string) (*yacspin.Spinner, error) {
	cfg := yacspin.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           yacspin.CharSets[14],
		Suffix:            " " + suffix,
		SuffixAutoColon:   true,
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		StopFailCharacter: "x",
		StopFailColors:    []string{"fgRed"},
	}

	spinner, err := yacspin.New(cfg)
	if err != nil {
		return nil, err
	}

	return spinner, nil
}

// First we need to open a file where user can set the semver, changelog contents,
// then we can insert that changelog contents into the source files before we call make to build
func runReleaseCmd(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("config")
	config := &config.Config{}
	err := config.Load(configFile)
	if err != nil {
		fmt.Printf("could not load config file: %v\n", err)
		os.Exit(1)
		return
	}

	spinner, err := initSpinner(fmt.Sprintf("Releasing v%s of %s", args[0], config.Repository))
	if err != nil {
		fmt.Println("could not init spinner")
		os.Exit(1)
		return
	}
	spinner.Start()

	newRelease, err := github.NewRelease(config, args, spinner)
	if err != nil {
		spinner.StopFailMessage(fmt.Sprintf("%v", err))
		spinner.StopFail()
		os.Exit(1)
		return
	}

	cl, err := changelog.HandleChangelog(newRelease.ProjectName, newRelease.Version, newRelease.Date, spinner)
	if err != nil {
		spinner.StopFailMessage(fmt.Sprintf("%v", err))
		spinner.StopFail()
		os.Exit(1)
		return
	}

	newRelease.Changelog = cl

	var binaryPath string
	skipBinary, _ := cmd.Flags().GetBool("skipBinary")
	if !skipBinary {
		// set project build path so we have a predictable location
		binaryPath = fmt.Sprintf(binaryPathFmt, newRelease.ProjectName, newRelease.Version)
		runBuildCmd(cmd, []string{newRelease.Version, binaryPath})
	}

	tokenFile, _ := cmd.Flags().GetString("tokenFile")
	err = newRelease.CreateGithubRelease(tokenFile, binaryPath, spinner)
	if err != nil {
		spinner.StopFailMessage(fmt.Sprintf("%v", err))
		spinner.StopFail()
		os.Exit(1)
		return
	}

	spinner.Suffix(" Finished release")
	spinner.Stop()
}

func init() {
	cmdRelease.Flags().Bool("skipBinary", false, "don't add a build asset for this release")
	cmdRelease.Flags().StringP("tokenFile", "t", "", "github api key file (default is $HOME/.github_token)")

	rootCmd.AddCommand(cmdRelease)
}
