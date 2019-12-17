package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseConfigFile(t *testing.T) {
	file := []byte(`
repoName: clintjedwards/toolkit
commands:
  build:
    - echo "Test"
`)

	expectedConfig := Config{
		RepoName: "clintjedwards/toolkit",
		Commands: map[string][]string{
			"build": []string{
				"echo \"Test\"",
			},
		},
	}

	config, err := parseConfigFile(file)
	if err != nil {
		t.Errorf("could not parse config: %v", err)
	}

	if !cmp.Equal(expectedConfig, config) {
		t.Errorf("config does not contain expected output; Diff below: \n%v", cmp.Diff(expectedConfig, config))
	}
}

func TestUnmarshalConfig(t *testing.T) {
	testConfig := Config{
		RepoName: "clintjedwards/toolkit",
	}

	expectedProject := projectInfo{
		Username: "clintjedwards",
		Name:     "toolkit",
	}

	var project projectInfo

	err := unmarshalConfig(testConfig, &project)
	if err != nil {
		t.Errorf("could not unmarshal config: %v", err)
	}

	if !cmp.Equal(expectedProject, project) {
		t.Errorf("project does not contain expected output; Diff below: \n%v", cmp.Diff(expectedProject, project))
	}
}
