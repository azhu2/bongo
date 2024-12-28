package entity

type WordListDAG struct {
	// Fragment is the letters so far up to this node
	Fragment []rune
	Children map[rune]*WordListDAG
	// IsWord marks if current node makes a valid word (still can have children)
	IsWord bool
}
