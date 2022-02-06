package set

// Write a set with map... anyway...
type IntSet map[int]interface{}

func NewIntSet() IntSet {
	return IntSet{}
}

func (s IntSet) Add(i int) {
	s[i] = nil
}

func (s IntSet) Has(i int) bool {
	_, ok := s[i]
	return ok
}
