package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Masterminds/semver"
	"github.com/clintjedwards/toolkit/config"
	"github.com/clintjedwards/toolkit/github"
	"github.com/clintjedwards/toolkit/sshutil"
	"github.com/clintjedwards/toolkit/utils"
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

type deploy struct {
	Host           string
	Name           string
	Version        string
	DownloadURL    string
	UploadFilePath string
	Commands       map[string][]string
}

func newDeploy(configFile string, args []string) (*deploy, error) {
	config := &config.Config{}
	err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config file: %w", err)
	}

	version, err := semver.NewVersion(args[0])
	if err != nil {
		log.Fatalf("could not parse semver string: %v", err)
	}

	projectUser, projectName, err := github.ParseGithubURL(config.Repository)
	if err != nil {
		return nil, fmt.Errorf("could not parse github URL: %w", err)
	}

	downloadURLFmt := "https://github.com/%s/%s/releases/download/v%s/%s"
	downloadURL := fmt.Sprintf(downloadURLFmt,
		projectUser, projectName, version.String(), projectName)

	uploadFilePath := fmt.Sprintf("/tmp/%s_%s", projectName, version.String())

	return &deploy{
		Host:           args[1],
		Name:           projectName,
		DownloadURL:    downloadURL,
		UploadFilePath: uploadFilePath,
		Version:        version.String(),
		Commands:       config.Commands,
	}, nil
}

func runDeployCmd(cmd *cobra.Command, args []string) {
	configFilePath, _ := cmd.Flags().GetString("config")

	newDeploy, err := newDeploy(configFilePath, args)
	if err != nil {
		log.Fatalf("could not create deploy instance: %v", err)
	}

	err = newDeploy.transferBinary()
	if err != nil {
		log.Fatalf("could not put binary on server: %v", err)
	}

	var commandList []string
	for _, rawCommand := range newDeploy.Commands["deploy"] {
		command, err := newDeploy.substituteTemplate(rawCommand)
		if err != nil {
			log.Fatalf("could not populate command template for command %s; %v", rawCommand, err)
		}

		commandList = append(commandList, command)
	}

	sshutil.RunCommandsOverSSH(newDeploy.Host, commandList)
}

func (d *deploy) transferBinary() error {

	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return fmt.Errorf("could not create tmp file: %w", err)
	}

	filename := file.Name()
	defer os.Remove(filename)

	log.Println("downloading binary")
	err = downloadFile(file, d.DownloadURL)
	if err != nil {
		return fmt.Errorf("could not download binary: %w", err)
	}

	// Upload the binary we just downloaded to server mentioned
	uploadCmdFmt := "scp %s %s:%s"
	uploadCmd := fmt.Sprintf(uploadCmdFmt, filename, d.Host, d.UploadFilePath)

	log.Println("uploading binary")
	_, err = utils.ExecuteBashCmd(uploadCmd, os.Environ(), "")
	if err != nil {
		return fmt.Errorf("could not run command '%s'; %w", uploadCmd, err)
	}

	return nil
}

// downloadFile downloads a file from url and writes it specified file
func downloadFile(file *os.File, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// takes in a command and returns the command with the variables from build struct filled in
func (d *deploy) substituteTemplate(command string) (string, error) {

	cmdBuffer := bytes.NewBuffer([]byte{})

	tmpl := template.Must(template.New("").Parse(command))
	err := tmpl.Execute(cmdBuffer, d)
	if err != nil {
		return "", err
	}

	return string(cmdBuffer.String()), nil
}

func init() {

	rootCmd.AddCommand(cmdDeploy)
}
