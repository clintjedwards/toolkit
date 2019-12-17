package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
)

const tokenEnv string = "GITHUB_TOKEN"
const defaultFileName string = ".github_token"

func createGithubRelease(project projectInfo, token string) error {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	release := &github.RepositoryRelease{
		TagName: github.String("v" + project.Version),
		Name:    github.String("v" + project.Version),
		Body:    github.String(string(project.Changelog)),
	}

	createdRelease, _, err := client.Repositories.CreateRelease(ctx, project.Username, project.Name, release)
	if err != nil {
		return err
	}

	if project.BuildPath == "" {
		return nil
	}

	// check if binary file exists first
	exists, err := fileExists(project.BuildPath)
	if !exists {
		return fmt.Errorf("could not find binary file: %s; %v", project.BuildPath, err)
	}
	if err != nil {
		return err
	}

	f, err := os.Open(project.BuildPath)
	if err != nil {
		return err
	}
	defer f.Close()

	client.Repositories.UploadReleaseAsset(ctx, project.Username, project.Name, createdRelease.GetID(),
		&github.UploadOptions{Name: project.Name}, f)

	return nil
}

// getGithubToken attempts to load a github token and returns an error if none exists
func getGithubToken(tokenFile string) (token string, err error) {

	// if user didn't set a tokenFile set the default
	if tokenFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return "", fmt.Errorf("could not get user home dir: %v", err)
		}

		tokenFile = fmt.Sprintf("%s/%s", home, defaultFileName)
	}

	setGithubTokenFromFile(tokenFile)

	token = os.Getenv(tokenEnv)

	if token == "" {
		return "", fmt.Errorf("env var GITHUB_TOKEN not found")
	}

	return token, nil
}

func setGithubTokenFromFile(filename string) {

	contents, _ := ioutil.ReadFile(filename)
	if len(contents) == 0 {
		log.Printf("could not find github token file: %v", filename)
		return
	}

	contents = bytes.TrimSpace(contents)

	os.Setenv(tokenEnv, string(contents))
}

// parseGithubURL parses the githubURL and return a username and repo name
func parseGithubURL(githubURL string) (username, projectName string, err error) {
	splitURL := strings.Split(githubURL, "/")

	if len(splitURL) != 2 {
		return "", "", fmt.Errorf("github URL not in correct format: username/repo")
	}

	return splitURL[0], splitURL[1], nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
