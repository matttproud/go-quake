package prog

import "fmt"

type stringRepo struct {
	vals map[int]string
	data []byte
}

func newStringRepo(data []byte) *stringRepo {
	return &stringRepo{vals: make(map[int]string), data: data}
}

func scan(data []byte, at int) (string, error) {
	if at < 0 || at >= len(data) {
		return "", fmt.Errorf("requested %#v, which is outside of %#v and %#v", at, 0, len(data))
	}
	for i, b := range data[at:] {
		if b == 0 {
			return string(data[at : at+i]), nil
		}
	}
	return string(data[at:]), nil
}

func (s *stringRepo) Lookup(at int) string {
	if v, ok := s.vals[at]; ok {
		return v
	}
	scanned, err := scan(s.data, at)
	if err != nil {
		panic(err)
	}
	s.vals[at] = scanned
	return scanned
}
