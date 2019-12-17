package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Masterminds/semver"
	"github.com/clintjedwards/toolkit/osutil"
	"github.com/spf13/cobra"
)

var cmdBuild = &cobra.Command{
	Use:   "build <semver> <path> --config <toolkit path>",
	Short: "Controls the build process for an application",
	Args:  cobra.MinimumNArgs(2),
	Run:   runBuildCmd,
}

func runBuildCmd(cmd *cobra.Command, args []string) {
	configFilePath, _ := cmd.Flags().GetString("config")
	project, err := getProjectInfo(configFilePath)
	if err != nil {
		log.Fatalf("could not get project info: %v", err)
	}

	// insert version into release struct
	version, err := semver.NewVersion(args[0])
	if err != nil {
		log.Fatalf("could not parse semver string: %v", err)
	}
	project.Version = version.String()

	// insert build_path into release struct
	buildPath := args[1]
	if err != nil {
		log.Fatalf("could not parse semver string: %v", err)
	}
	project.BuildPath = buildPath

	echoCommands, _ := cmd.Flags().GetBool("echoCommands")
	env := os.Environ()

	var commandList []string
	for _, rawCommand := range project.Commands["build"] {
		command, err := populateCommandTemplate(rawCommand, project)
		if err != nil {
			log.Fatalf("could not populate command template for command %s; %v", rawCommand, err)
		}

		commandList = append(commandList, command)
	}

	for _, command := range commandList {
		if echoCommands {
			fmt.Println("> " + command)
		}

		output, err := osutil.ExecuteBashCmd(command, env, "")
		if err != nil {
			log.Fatalf("could not run command '%s'; %v", command, err)
		}

		hideOutput, _ := cmd.Flags().GetBool("hideOutput")
		if !hideOutput {
			fmt.Println(string(output))
		}
	}

}

func init() {
	rootCmd.AddCommand(cmdBuild)
}
