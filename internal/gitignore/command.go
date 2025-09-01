package gitignore

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v74/github"
)

const (
	repoOwner = "github"
	repoName  = "gitignore"
)

type GetGitignoreCommand struct {
	language   string
	outputPath string
	client     *github.Client
}

func NewGetGitignoreCommand(language, outputPath string, client *github.Client) *GetGitignoreCommand {
	return &GetGitignoreCommand{
		language:   language,
		outputPath: outputPath,
		client:     client,
	}
}

func (cmd *GetGitignoreCommand) Execute() error {
	log.Printf("Downloading .gitignore for %s in %s\n", cmd.language, cmd.outputPath)

	log.Println("Retrieving supported languages")
	languages, _ := cmd.fetchSupportedLanguages()

	url, exists := findMatch(cmd.language, languages)
	if !exists {
		return errors.New("language not supported")
	}

	log.Printf("Found match: %s -> %s\n", cmd.language, url)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("failed to download .gitignore file")
	}

	contents, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	parent := filepath.Dir(cmd.outputPath)
	err = os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(cmd.outputPath, contents, os.ModePerm)
}

func findMatch(language string, languages map[string]string) (string, bool) {
	lowerLang := strings.ToLower(language)
	url, exists := languages[lowerLang]
	return url, exists
}

func (cmd *GetGitignoreCommand) fetchSupportedLanguages() (map[string]string, error) {
	_, dir, _, err := cmd.client.Repositories.GetContents(
		context.Background(),
		repoOwner,
		repoName,
		"",
		&github.RepositoryContentGetOptions{})

	if err != nil {
		return nil, err
	}

	languages := make(map[string]string, len(dir))

	for _, item := range dir {
		if item.GetType() != "file" {
			continue
		}

		if !strings.HasSuffix(item.GetName(), ".gitignore") {
			continue
		}

		name := strings.TrimSuffix(item.GetName(), ".gitignore")
		languages[strings.ToLower(name)] = item.GetDownloadURL()
	}

	return languages, nil
}
