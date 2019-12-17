package main

import (
	"bytes"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Config represents per project information pulled from the config file
type Config struct {
	RepoName string `yaml:"repoName"`
	Commands map[string][]string
}

// parseConfigFile unmarshals yml config file into release struct
func parseConfigFile(file []byte) (Config, error) {
	var c Config

	err := yaml.Unmarshal(file, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func unmarshalConfig(config Config, r *projectInfo) error {

	username, projectName, err := parseGithubURL(config.RepoName)
	if err != nil {
		return err
	}

	r.Username = username
	r.Name = projectName
	r.Commands = config.Commands

	return nil
}

// takes in a command and returns the command with the variables filled in
func populateCommandTemplate(command string, info projectInfo) (string, error) {

	b := bytes.NewBuffer([]byte{})

	tmpl := template.Must(template.New("").Parse(command))
	err := tmpl.Execute(b, info)
	if err != nil {
		return "", err
	}

	return string(b.String()), nil
}
