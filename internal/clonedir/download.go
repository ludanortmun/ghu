package clonedir

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v74/github"
)

// startDownload will start the download process for the given path.
// It will check if the path is a file or a directory and handle it accordingly.
func (cmd *DownloadCommand) startDownload(path string) (*fileTree, error) {
	opts := &github.RepositoryContentGetOptions{}
	if cmd.target.Ref != "" {
		opts.Ref = cmd.target.Ref
	}

	file, dir, _, err := cmd.client.Repositories.GetContents(
		context.Background(),
		cmd.target.Owner,
		cmd.target.Repository,
		path,
		opts,
	)

	if err != nil {
		return nil, err
	}

	// If the target is a file, download it directly
	if file != nil {
		return downloadFile(file)
	}

	// If the target is not a file, it will be a directory.
	pathParts := strings.Split(path, "/")
	dirname := pathParts[len(pathParts)-1]
	return cmd.downloadDirectory(dirname, dir)
}

// downloadFile will download the file from the GitHub repository using its download URL.
func downloadFile(item *github.RepositoryContent) (*fileTree, error) {
	res, err := http.Get(item.GetDownloadURL())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &fileTree{
		name:     item.GetName(),
		content:  bytes,
		children: nil,
	}, nil
}

// downloadDirectory will download the contents of a directory from the GitHub repository.
// It will recursively download all files and directories within the given directory.
func (cmd *DownloadCommand) downloadDirectory(dirname string, items []*github.RepositoryContent) (*fileTree, error) {

	result := &fileTree{
		name:    dirname,
		content: nil,
	}

	children := make([]*fileTree, len(items))

	for i, item := range items {
		if item.GetType() == "file" {
			f, err := downloadFile(item)
			if err != nil {
				return nil, err
			}
			children[i] = f
		} else if item.GetType() == "dir" {
			dir, err := cmd.startDownload(item.GetPath())
			if err != nil {
				return nil, err
			}
			children[i] = dir
		}
	}

	result.children = children
	return result, nil
}
