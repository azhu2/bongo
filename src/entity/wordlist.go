package entity

type WordList struct {
	Root *DAGNode
	// map of column index -> rune -> nodes
	NodeMap map[int]map[rune][]*DAGNode
}

type DAGNode struct {
	// Fragment is the letters so far up to this node
	Fragment []rune
	Children map[rune]*DAGNode
	Prev     *DAGNode
	// IsWord marks if current node makes a valid word (still can have children)
	IsWord bool
}
