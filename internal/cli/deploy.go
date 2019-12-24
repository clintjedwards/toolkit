package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/clintjedwards/toolkit/osutil"
	"github.com/spf13/cobra"
)

// Deploy should initiate a connection via ssh to the deployment server
// download the binary from github and then run that binary with appropriate settings

var cmdDeploy = &cobra.Command{
	Use:   "deploy <semver> <user@host>",
	Short: "Controls the deployment process for the application",
	Long: `Downloads version of project specified using github release.
Then initiates a ssh connection using user host combination and runs commands
under "deploy" in configuration file`,
	Args: cobra.MinimumNArgs(2),
	Run:  runDeployCmd,
}

const downloadURLFmt string = "https://github.com/%s/%s/releases/download/v%s/%s"

func runDeployCmd(cmd *cobra.Command, args []string) {
	configFilePath, _ := cmd.Flags().GetString("config")
	project, err := getProjectInfo(configFilePath)
	if err != nil {
		log.Fatalf("could not get project info: %v", err)
	}

	hostParts := strings.Split(args[1], "@")

	// insert version into release struct
	version, err := semver.NewVersion(args[0])
	if err != nil {
		log.Fatalf("could not parse semver string: %v", err)
	}
	project.Version = version.String()

	downloadURL := fmt.Sprintf(downloadURLFmt,
		project.Username, project.Name, project.Version, project.Name)

	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		log.Fatalf("could not create tmp file: %v", err)
	}

	filename := file.Name()
	defer os.Remove(filename)

	log.Println("downloading binary")
	err = downloadFile(file, downloadURL)
	if err != nil {
		log.Fatalf("could not download binary: %v", err)
	}

	uploadFilePath := fmt.Sprintf("/tmp/%s_%s", project.Name, project.Version)
	uploadCmd := fmt.Sprintf("scp %s %s@%s:%s",
		filename, hostParts[0], hostParts[1], uploadFilePath)

	_, err = osutil.ExecuteBashCmd(uploadCmd, os.Environ(), "")
	if err != nil {
		log.Fatalf("could not run command '%s'; %v", uploadCmd, err)
	}

	var commandList []string
	for _, rawCommand := range project.Commands["deploy"] {
		command, err := populateCommandTemplate(rawCommand, project)
		if err != nil {
			log.Fatalf("could not populate command template for command %s; %v", rawCommand, err)
		}

		commandList = append(commandList, command)
	}

	runCommandsOverSSH(hostParts[0], hostParts[1], commandList)
}

func downloadFile(file *os.File, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func init() {

	rootCmd.AddCommand(cmdDeploy)
}
