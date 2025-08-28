package clonedir

import (
	"errors"
	"os"
)

type fileTree struct {
	name     string
	content  []byte
	children []*fileTree
}

func (f *fileTree) isDirectory() bool {
	return f.children != nil
}

func saveToDisk(node *fileTree, path string) error {
	if node == nil {
		return errors.New("no files to save")
	}

	if !node.isDirectory() {
		return os.WriteFile(path+"/"+node.name, node.content, os.ModePerm)
	}

	newBasePath := path + "/" + node.name

	// Ensure the directory exists
	// We call this in the directory path to ensure it's only called once for each directory
	if err := os.MkdirAll(newBasePath, os.ModePerm); err != nil {
		return err
	}

	for _, child := range node.children {
		if err := saveToDisk(child, newBasePath); err != nil {
			return err
		}
	}

	return nil
}
