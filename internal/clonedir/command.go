package clonedir

import (
	"github.com/google/go-github/v74/github"
	"github.com/ludanortmun/ghu/internal"
)

type DownloadCommand struct {
	target     internal.GitHubTarget
	outputPath string
	client     *github.Client
}

func NewDownloadCommand(target internal.GitHubTarget, outputPath string, client *github.Client) *DownloadCommand {
	return &DownloadCommand{
		target:     target,
		outputPath: outputPath,
		client:     client,
	}
}

func (cmd *DownloadCommand) Execute() error {
	files, err := cmd.startDownload(cmd.target.Directory)
	if err != nil {
		return err
	}

	return saveToDisk(files, cmd.outputPath)
}
