package util

type Set[Elem comparable] struct {
	data map[Elem]struct{}
}

// Add adds an element to a Set.
func (s *Set[Elem]) Add(v Elem) {
	if s.data == nil {
		s.data = make(map[Elem]struct{})
	}
	s.data[v] = struct{}{}
}

// Contains reports whether v is in the Set.
func (s *Set[Elem]) Contains(v Elem) bool {
	if s.data == nil {
		return false
	}
	_, ok := s.data[v]
	return ok
}
