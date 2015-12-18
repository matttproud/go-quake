package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/matttproud/go-quake/path"
)

var (
	flagBaseDir string
	flagGame    string
)

func GamePath() (*path.Path, error) {
	var paths []string
	if flagGame != "" {
		paths = append(paths, makePath(flagGame))
	}
	id1Path := makePath("id1")
	if _, err := os.Stat(id1Path); err != nil {
		return nil, fmt.Errorf("could not access id1 directory: %v", err)
	}
	paths = append(paths, id1Path)
	for _, p := range paths {
		paks, err := filepath.Glob(filepath.Join(p, "*.pak"))
		if err != nil {
			return nil, err
		}
		paths = append(paths, paks...)
	}
	return path.New(paths...)
}

func makePath(p string) string { return filepath.Join(flagBaseDir, p) }

func init() {
	flag.StringVar(&flagBaseDir, "basedir", ".", "the directory that contains id1 directory")
	flag.StringVar(&flagGame, "game", "", "an alternative game directory")
}
