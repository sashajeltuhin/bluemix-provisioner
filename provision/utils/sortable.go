package utils

import (
	"sort"
)

type SortablePair struct {
	Key   string
	Value string
}

type SortedPairArr []SortablePair

func (s SortedPairArr) Len() int {
	return len(s)
}

func (s SortedPairArr) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedPairArr) Less(i, j int) bool {
	return s[i].Value < s[j].Value
}

func SortMapbyVal(m map[string]string) SortedPairArr {
	count := len(m)
	i := 0
	var list SortedPairArr = make(SortedPairArr, count)
	for key := range m {
		list[i] = SortablePair{key, m[key]}
		i++
	}

	// Sort by val
	sort.Sort(SortedPairArr(list))
	return list
}
