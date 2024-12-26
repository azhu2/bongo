package entity

// Stack is a basic generic stack
type Stack[T any] struct {
	root *stackNode[T]
}

type stackNode[T any] struct {
	val  T
	prev *stackNode[T]
}

func (s *Stack[T]) IsEmpty() bool {
	return s.root == nil
}

func (s *Stack[T]) Push(val T) {
	new := stackNode[T]{
		val:  val,
		prev: s.root,
	}
	s.root = &new
}

func (s *Stack[T]) Pop() T {
	cur := s.root.val
	s.root = s.root.prev
	return cur
}

func (s *Stack[T]) Peek() T {
	return s.root.val
}
