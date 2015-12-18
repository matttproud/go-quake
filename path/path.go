// Package path provides asset access mechanisms.
package path

import (
	"io"
	"os"
	"path/filepath"

	"github.com/matttproud/go-quake/pak"
)

type dirPath struct {
	base    string
	handles map[string]*os.File
}

func fAsSectionReader(f *os.File) *io.SectionReader {
	stat, _ := f.Stat()
	return io.NewSectionReader(f, 0, stat.Size())
}

func (p *dirPath) Load(n string) (*io.SectionReader, error) {
	fp := filepath.Join(p.base, n)
	if f, ok := p.handles[fp]; ok {
		return fAsSectionReader(f), nil
	}
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	p.handles[fp] = f
	return fAsSectionReader(f), nil
}

func (p *dirPath) Close() {
	for _, f := range p.handles {
		f.Close()
	}
	p.handles = nil
}

type pakPath struct {
	f       *os.File
	p       *pak.Pak
	handles map[string]*pak.File
}

func pAsSectionReader(f *pak.File) *io.SectionReader {
	return io.NewSectionReader(f, 0, int64(f.Size))
}

func (p *pakPath) Load(n string) (*io.SectionReader, error) {
	if f, ok := p.handles[n]; ok {
		return pAsSectionReader(f), nil
	}
	for _, f := range p.p.Files {
		if f.Name != n {
			continue
		}
		p.handles[n] = f
		return pAsSectionReader(f), nil
	}
	return nil, os.ErrNotExist
}

func (p *pakPath) Close() {
	p.f.Close()
}

type Path struct {
	sources []source
}

type source interface {
	Load(string) (*io.SectionReader, error)
	Close()
}

func New(searchPaths ...string) (*Path, error) {
	p := &Path{sources: make([]source, len(searchPaths))}
	for i, sp := range searchPaths {
		fi, err := os.Stat(sp)
		if err != nil {
			return nil, err
		}
		if fi.IsDir() {
			dp := &dirPath{base: sp, handles: make(map[string]*os.File)}
			p.sources[i] = dp
			continue
		}
		fh, err := os.Open(sp)
		if err != nil {
			return nil, err
		}
		pp, err := pak.Open(fh)
		if err != nil {
			return nil, err
		}
		p.sources[i] = &pakPath{f: fh, p: pp, handles: make(map[string]*pak.File)}
	}
	return p, nil
}

func (p *Path) Load(name string) (*io.SectionReader, error) {
	for _, s := range p.sources {
		r, err := s.Load(name)
		switch {
		case err == nil:
			return r, nil
		case !os.IsNotExist(err):
			return nil, err
		}
	}
	return nil, os.ErrNotExist
}

func (p *Path) Close() {
	for _, s := range p.sources {
		s.Close()
	}
}
