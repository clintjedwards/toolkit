package config

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {

	file, err := ioutil.TempFile("/tmp", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.Write([]byte(`repository: clintjedwards/toolkit
commands:
  build:
    - echo "Test"
`))

	expectedConfig := Config{
		Repository: "clintjedwards/toolkit",
		Commands: map[string][]string{
			"build": []string{
				"echo \"Test\"",
			},
		},
	}

	config := Config{}
	err = config.Load(file.Name())
	if err != nil {
		t.Errorf("could not parse config: %v", err)
	}

	if !cmp.Equal(expectedConfig, config) {
		t.Errorf("config does not contain expected output; Diff below: \n%v", cmp.Diff(expectedConfig, config))
	}
}
