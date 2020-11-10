package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/Masterminds/semver"
	"github.com/clintjedwards/toolkit/config"
	"github.com/clintjedwards/toolkit/github"
	"github.com/clintjedwards/toolkit/utils"
	"github.com/spf13/cobra"
)

var cmdBuild = &cobra.Command{
	Use:   "build <semver> <path>",
	Short: "Controls the build process for an application",
	Long: `Runs the commands under 'build' in config file to build the application
Injects variables in template format: {{.ExampleVar}}

Variables injected: ProjectName, Path, Version, VersionFull
`,
	Args: cobra.MinimumNArgs(2),
	Run:  runBuildCmd,
}

type build struct {
	ProjectName string // the project name grabbed from the repository
	Path        string // path where binary will be build
	Version     string // semver without the v; ex: 1.0.0
	VersionFull string // ex: <semver>_<epoch>_<commit>
	Commands    map[string][]string
}

func newBuild(configFile string, args []string) (*build, error) {
	config := &config.Config{}
	err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config file: %w", err)
	}

	// insert version into build struct
	version, err := semver.NewVersion(args[0])
	if err != nil {
		return nil, fmt.Errorf("could not parse semver string: %w", err)
	}

	versionFull, err := getVersionFull(version.String())
	if err != nil {
		return nil, fmt.Errorf("could not get full version string: %w", err)
	}

	_, projectName, err := github.ParseGithubURL(config.Repository)
	if err != nil {
		return nil, fmt.Errorf("could not parse github URL: %w", err)
	}

	return &build{
		ProjectName: projectName,
		Path:        args[1],
		Version:     version.String(),
		VersionFull: versionFull,
		Commands:    config.Commands,
	}, nil
}

func runBuildCmd(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("config")

	newBuild, err := newBuild(configFile, args)
	if err != nil {
		log.Fatalf("could not create build instance: %v", err)
	}

	echoCommands, _ := cmd.Flags().GetBool("echoCommands")
	env := os.Environ()

	var commandList []string
	for _, rawCommand := range newBuild.Commands["build"] {
		command, err := newBuild.substituteTemplate(rawCommand)
		if err != nil {
			log.Fatalf("could not populate command template for command %s; %v", rawCommand, err)
		}

		commandList = append(commandList, command)
	}

	for _, command := range commandList {
		if echoCommands {
			fmt.Println("> " + command)
		}

		output, err := utils.ExecuteBashCmd(command, env, "")
		if err != nil {
			log.Fatalf("could not run command '%s'; %v", command, err)
		}

		hideOutput, _ := cmd.Flags().GetBool("hideOutput")
		if !hideOutput && len(output) != 0 {
			fmt.Println(string(output))
		}
	}
}

// takes in a command and returns the command with the variables from build struct filled in
func (b *build) substituteTemplate(command string) (string, error) {

	cmdBuffer := bytes.NewBuffer([]byte{})

	tmpl := template.Must(template.New("").Parse(command))
	err := tmpl.Execute(cmdBuffer, b)
	if err != nil {
		return "", err
	}

	return string(cmdBuffer.String()), nil
}

// getVersionFull generates a long version string in format <semver>_<epoch>_<githash>
func getVersionFull(semver string) (string, error) {
	versionFmt := "%s_%s_%s"
	gitCmd := "git rev-parse --short HEAD"
	dateCmd := "date +%s"

	env := os.Environ()
	commit, err := utils.ExecuteBashCmd(gitCmd, env, "")
	epoch, err := utils.ExecuteBashCmd(dateCmd, env, "")
	if err != nil {
		return "", fmt.Errorf("could not determine version: %w", err)
	}

	commit = bytes.TrimSpace(commit)
	epoch = bytes.TrimSpace(epoch)

	return fmt.Sprintf(versionFmt, semver, epoch, commit), nil
}

func init() {
	rootCmd.AddCommand(cmdBuild)
}
