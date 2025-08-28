package internal

import (
	"errors"
	"net/url"
	"strings"
)

type GitHubTarget struct {
	Owner      string
	Repository string
	Directory  string
	Ref        string
}

// InferTargetFromUrl will take a valid GitHub URL and create a GitHubTarget object from it.
// A GitHub URL will take the form "https://github.com/{owner}/{repo}/(tree/<ref>)?/(<path>/<to>/<root>)?", where:
// - "tree/<ref>" is optional, if not present it defaults to the default branch of the repo
// - "<ref>" can either be a commit hash or branch
// - "<path>/<to>/<root>" is optional
// - If "<path>/<to>/<root>" is present, then "tree/<ref>" MUST be present
func InferTargetFromUrl(githubUrl string) (GitHubTarget, error) {
	target := GitHubTarget{}

	_url, ok := strings.CutPrefix(githubUrl, "https://github.com/")
	if !ok {
		return GitHubTarget{}, errors.New(`invalid GitHub URL`)
	}

	_url, err := url.QueryUnescape(_url)
	if err != nil {
		return GitHubTarget{}, errors.New(`invalid GitHub URL`)
	}

	parts := strings.Split(_url, "/")

	if len(parts) < 2 {
		return GitHubTarget{}, errors.New(`invalid GitHub URL`)
	}
	target.Owner = parts[0]
	target.Repository = parts[1]

	// Nothing more to process, we are at the root of the repo in the default branch
	if len(parts) == 2 {
		return target, nil
	}

	// Otherwise, the URL will include at least the "/tree/<ref>" part
	if len(parts) < 4 || parts[2] != "tree" {
		return GitHubTarget{}, errors.New(`invalid GitHub URL`)
	}

	target.Ref = parts[3]

	// If the target includes the <path>/<to>/<root> part.
	// It will be the rest of the string.
	if len(parts) > 4 {
		target.Directory = strings.Join(parts[4:], "/")
	}

	return target, nil
}
