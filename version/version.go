package version

import (
	"errors"
	"strings"
)

// Info represents application information gleaned from the version string
type Info struct {
	Semver string
	Epoch  string
	Hash   string
}

// Parse takes a very specific version string format:
// <semver>_<epoch_time>_<git_hash>
// and returns its individual parts
func Parse(version string) (Info, error) {
	versionTuple := strings.Split(version, "_")

	if len(versionTuple) != 3 {
		return Info{},
			errors.New("version not in correct format: <semver>_<epoch_time>_<git_hash>")
	}

	return Info{
		Semver: versionTuple[0],
		Epoch:  versionTuple[1],
		Hash:   versionTuple[2],
	}, nil
}
