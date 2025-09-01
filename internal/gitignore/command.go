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

	url, ok := findMatch(cmd.language, languages)
	if !ok {
		return errors.New("language not supported")
	}

	log.Printf("Found match: %s -> %s\n", cmd.language, url)

	content, err := downloadFile(url)
	if err != nil {
		return err
	}

	return saveToFile(cmd.outputPath, content)
}

func findMatch(language string, languages map[string]string) (string, bool) {
	l := strings.ToLower(language)
	url, ok := languages[l]

	if !ok {
		url, ok = languages[aliasesMap[l]]
		if ok {
			log.Printf("Using alias: %s -> %s\n", language, aliasesMap[l])
		}
	}

	return url, ok
}

func saveToFile(path string, content []byte) error {
	parent := filepath.Dir(path)
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, os.ModePerm)
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download file: " + resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
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
