package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

const defaultEditor string = "vi"
const editorEnvVar string = "EDITOR"
const dateFmt string = "%s %d, %d"
const binaryPathFmt string = "/tmp/%s_%s"

var cmdRelease = &cobra.Command{
	Use:   "release <semver> --config <toolkit path>",
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

	configFilePath, _ := cmd.Flags().GetString("config")
	project, err := getProjectInfo(configFilePath)
	if err != nil {
		log.Fatalf("could not get project info: %v", err)
	}

	tokenFile, _ := cmd.Flags().GetString("tokenFile")
	token, err := getGithubToken(tokenFile)
	if err != nil {
		log.Fatalf("could not get github token: %v", err)
	}

	// insert version into release struct
	version, err := semver.NewVersion(args[0])
	if err != nil {
		log.Fatalf("could not parse semver string: %v", err)
	}
	project.Version = version.String()

	// insert date into release struct
	year, month, day := time.Now().Date()
	project.Date = fmt.Sprintf(dateFmt, month, day, year)

	// allow user to write changelog and store it
	skipChangelog, _ := cmd.Flags().GetBool("skipChangelog")
	changelog, err := getChangelog(skipChangelog, project)
	if err != nil {
		log.Fatalf("could not get changelog: %v", err)
	}
	project.Changelog = changelog

	// create github release and upload binary
	skipBinary, _ := cmd.Flags().GetBool("skipBinary")
	if !skipBinary {
		// set project build path so we have a predictable location
		project.BuildPath = fmt.Sprintf(binaryPathFmt, project.Name, project.Version)
		runBuildCmd(cmd, []string{project.Version, project.BuildPath})
	}

	err = createGithubRelease(project, token)
	if err != nil {
		log.Fatalf("could not create github release: %v", err)
	}
}

func getChangelog(skip bool, info projectInfo) ([]byte, error) {
	if skip {
		return []byte{}, nil
	}

	changelog, err := writeChangelog(info)
	if err != nil {
		return []byte{}, err
	}

	return changelog, nil
}

func init() {
	cmdRelease.Flags().Bool("skipChangelog", false, "don't add a changelog for this release")
	cmdRelease.Flags().Bool("skipBuild", false, "don't add a build asset for this release")
	cmdRelease.Flags().StringP("tokenFile", "t", "", "github api key file (default is $HOME/.github_token)")

	rootCmd.AddCommand(cmdRelease)
}
