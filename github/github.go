package github

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/clintjedwards/toolkit/changelog"
	"github.com/clintjedwards/toolkit/config"
	"github.com/clintjedwards/toolkit/utils"
	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
)

const tokenEnv string = "GITHUB_TOKEN"
const tokenFileName string = ".github_token"
const dateFmt string = "%s %d, %d"

// Release contains information pertaining to a specific github release
type Release struct {
	User        string
	Date        string // date in format: month day, year
	Changelog   []byte
	Repository  string // full repository name from config
	ProjectName string // the project name grabbed from the repository
	Version     string // semver without the v; ex: 1.0.0
	VersionFull string // ex: <semver>_<epoch>_<commit>
	Commands    map[string][]string
}

// NewRelease creates a prepopulated release struct using the config file and other sources
func NewRelease(configFile string, args []string) (*Release, error) {
	config := &config.Config{}
	err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config file: %w", err)
	}

	// insert version into build struct
	version, err := semver.NewVersion(args[0])
	if err != nil {
		return nil, fmt.Errorf("could not parse semver string: %w", err)
	}

	versionFull, err := getVersionFull(version.String())
	if err != nil {
		return nil, fmt.Errorf("could not get full version string: %w", err)
	}

	user, projectName, err := ParseGithubURL(config.Repository)
	if err != nil {
		return nil, fmt.Errorf("could not parse github URL: %w", err)
	}

	// insert date into release struct
	year, month, day := time.Now().Date()
	date := fmt.Sprintf(dateFmt, month, day, year)

	cl, err := changelog.HandleChangelog(projectName, version.String(), date)
	if err != nil {
		return nil, fmt.Errorf("could not get changelog: %w", err)
	}

	return &Release{
		Changelog:   cl,
		Commands:    config.Commands,
		Date:        date,
		ProjectName: projectName,
		Repository:  config.Repository,
		User:        user,
		Version:     version.String(),
		VersionFull: versionFull,
	}, nil
}

// CreateGithubRelease cuts a new release, tags the current commit with semver, and uploads the changelog as a description
func (r *Release) CreateGithubRelease(tokenFile, binaryPath string) error {
	ctx := context.Background()

	token, err := getGithubToken(tokenFile)
	if err != nil {
		return fmt.Errorf("could not get github token: %w", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	release := &github.RepositoryRelease{
		TagName: github.String("v" + r.Version),
		Name:    github.String("v" + r.Version),
		Body:    github.String(string(r.Changelog)),
	}

	createdRelease, _, err := client.Repositories.CreateRelease(ctx, r.User, r.ProjectName, release)
	if err != nil {
		return err
	}

	if binaryPath == "" {
		return nil
	}

	_, err = os.Stat(binaryPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("could not find binary file: %s; %w", binaryPath, err)
	}

	f, err := os.Open(binaryPath)
	if err != nil {
		return err
	}
	defer f.Close()

	client.Repositories.UploadReleaseAsset(ctx, r.User, r.ProjectName, createdRelease.GetID(),
		&github.UploadOptions{Name: r.ProjectName}, f)

	return nil
}

// getGithubToken attempts to load a github token and returns an error if none exists
func getGithubToken(tokenFile string) (token string, err error) {

	token = os.Getenv(tokenEnv)

	if token != "" {
		return token, nil
	}

	if tokenFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return "", fmt.Errorf("could not get user home dir: %w", err)
		}

		tokenFile = fmt.Sprintf("%s/%s", home, tokenFileName)
	}

	rawToken, err := setGithubTokenFromFile(tokenFile)
	if err != nil {
		return "", err
	}

	return string(rawToken), nil
}

func setGithubTokenFromFile(filename string) ([]byte, error) {

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not find github token: %s; %w", filename, err)
	}
	if len(contents) == 0 {
		return nil, fmt.Errorf("could not load github token contents empty: %s", filename)
	}

	token := bytes.TrimSpace(contents)
	return token, nil
}

// ParseGithubURL parses the githubURL and return a username and repo name
func ParseGithubURL(githubURL string) (username, projectName string, err error) {
	splitURL := strings.Split(githubURL, "/")

	if len(splitURL) != 2 {
		return "", "", fmt.Errorf("github URL not in correct format: username/repo")
	}

	return splitURL[0], splitURL[1], nil
}

// getVersionFull generates a long version string in format <semver>_<epoch>_<githash>
func getVersionFull(semver string) (string, error) {
	versionFmt := "%s_%s_%s"
	gitCmd := "git rev-parse --short HEAD"
	dateCmd := "date +%s"

	env := os.Environ()
	commit, err := utils.ExecuteBashCmd(gitCmd, env, "")
	epoch, err := utils.ExecuteBashCmd(dateCmd, env, "")
	if err != nil {
		return "", fmt.Errorf("could not determine version: %w", err)
	}

	commit = bytes.TrimSpace(commit)
	epoch = bytes.TrimSpace(epoch)

	return fmt.Sprintf(versionFmt, semver, epoch, commit), nil
}
