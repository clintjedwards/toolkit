package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents per project configuration loaded from the toolkit.yml file
type Config struct {
	Repository string `yaml:"repository"` // In form: username/project_name
	Commands   map[string][]string
}

// Load reads in a config file and unmarshals it into config struct
func (c *Config) Load(filename string) error {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(f, &c)
	return err
}
