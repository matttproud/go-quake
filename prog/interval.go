package prog

import (
	"fmt"
	"sort"
)

type interval []int

func newInterval(heads ...int) interval {
	hs := make([]int, len(heads))
	copy(hs, heads)
	sort.Ints(hs)
	return interval(hs)
}

func (i interval) End(start int) (e int) {
	f := func(j int) bool { return i[j] >= start }
	j := sort.Search(len(i), f)
	switch {
	case j+1 == len(i):
		return -1
	case j+1 < len(i):
		return i[j+1]
	default:
		panic(fmt.Sprintf("%#v: does not contain %v", i, start))
	}
}
