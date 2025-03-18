package set

type Set[T comparable] struct {
	lookup    map[T]struct{}
	elements []T
}

func (s *Set[T]) Add(elem T) {
    s.lookup[elem] = struct{}{}
    s.elements = append(s.elements, elem)
}

func (s *Set[T]) Contains(elem T) bool {
    _, ok := s.lookup[elem]
    return ok
}
