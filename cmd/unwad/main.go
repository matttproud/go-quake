// unwad extracts files from wad archives.
package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/matttproud/go-quake/wad"
)

var (
	matchExp  string
	matchRe   *regexp.Regexp
	dest      string
	overwrite bool
)

func main() {
	flag.Parse()
	matchRe = regexp.MustCompile(matchExp)
	for _, fp := range flag.Args() {
		f, err := os.Open(fp)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		if err := extractPak(f); err != nil {
			log.Println(err)
			return
		}
	}
}

func extractPak(r io.ReaderAt) error {
	p, err := wad.Open(r)
	if err != nil {
		return err
	}
	for _, f := range p.Files {
		if !matchRe.MatchString(f.Name) {
			continue
		}
		if err := extract(f); err != nil {
			return err
		}
	}
	return nil
}

func extract(w *wad.File) error {
	destPath := filepath.Join(dest, w.Name)
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0744); err != nil {
		return err
	}
	if _, err := os.Stat(destDir); !os.IsNotExist(err) && !overwrite {
		return nil
	}
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, w)
	if err == nil {
		log.Println("Extracted", w.Name)
	}
	return err
}

func init() {
	flag.StringVar(&matchExp, "match_exp", ".*", "the regular expression for file matches")
	flag.StringVar(&dest, "dest", ".", "the path for the output files")
	flag.BoolVar(&overwrite, "overwrite", false, "whether to overwrite pre-existing files")
}
