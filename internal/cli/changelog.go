package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// getEditorPath attempts to find a suitible editor
// returns an editor binary and argument string
// ex. /usr/bin/vscode --wait
func getEditorPath() (string, error) {

	editorFromEnv := os.Getenv(editorEnvVar)
	if editorFromEnv != "" {
		return editorFromEnv, nil
	}

	path, err := exec.LookPath(defaultEditor)
	if err != nil {
		return "", err
	}

	return path, nil
}

// openFileInEditor attempts to find an editor and open a specific file
func openFileInEditor(filename string) error {
	editorPath, err := getEditorPath()
	if err != nil {
		return err
	}

	// split the path parsed into parts so we can manipulate and add into Command func
	editorPathParts := strings.Split(editorPath, " ")
	editorPathParts = append(editorPathParts, filename)

	cmd := exec.Command(editorPathParts[0], editorPathParts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// writeChangelog opens a temp file for editing so the user can insert values
func writeChangelog(info projectInfo) ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return nil, err
	}

	filename := file.Name()
	defer os.Remove(filename)

	tmpl := template.Must(template.New("").Parse(pretext))
	err = tmpl.Execute(file, info)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	err = openFileInEditor(filename)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	changelog := removeFileComments(bytes)

	return changelog, nil
}

func removeFileComments(data []byte) []byte {

	var newFile [][]byte
	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		if !bytes.HasPrefix(bytes.TrimSpace(line), []byte("//")) {
			newFile = append(newFile, line)
		}
	}

	return bytes.Join(newFile, []byte("\n"))
}
