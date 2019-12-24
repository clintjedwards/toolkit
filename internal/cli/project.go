package main

import (
	"fmt"
	"io/ioutil"
)

type projectInfo struct {
	Changelog   []byte
	BuildPath   string
	Date        string // human readable date ex: June 19, 2019
	Name        string
	Username    string
	Version     string
	VersionFull string
	Commands    map[string][]string
}

func getProjectInfo(configFilePath string) (projectInfo, error) {
	p := projectInfo{}

	configText, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return projectInfo{}, fmt.Errorf("could not get config file: %v", err)
	}

	// fill release struct with config values
	config, err := parseConfigFile(configText)
	if err != nil {
		return projectInfo{}, fmt.Errorf("could not parse config file: %v", err)
	}

	err = unmarshalConfig(config, &p)
	if err != nil {
		return projectInfo{}, fmt.Errorf("could not unmarshal config info: %v", err)
	}

	return p, nil
}
