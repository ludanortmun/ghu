package internal

import (
	"log"

	"github.com/google/go-github/v74/github"
)

func CreateGithubClient() *github.Client {
	client := github.NewClient(nil)
	ghToken, ok := GetAuthToken()

	if ok {
		log.Println("Using authenticated GitHub API client.")
		client = client.WithAuthToken(ghToken)
	} else {
		log.Println("No PAT was found, using unauthenticated GitHub API client. " +
			"If you want to access private repositories, please set a PAT using the `ghws auth set-token` command.")
	}

	return client
}
