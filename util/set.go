package util

type Set map[int]bool

func NewSet() Set {
	return make(Set)
}

func (s Set) Add(item int) {
	s[item] = true
}

func (s Set) Contains(item int) bool {
	return s[item]
}

func (s Set) Remove(item int) {
	delete(s, item)
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Slice() []int {
	var intSlice []int
	for item := range s {
		intSlice = append(intSlice, item)
	}
	return intSlice
}
