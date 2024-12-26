package entity

type WordListDAG struct {
	// Fragment is the letters so far up to this node
	Fragment []rune
	Children map[rune]*WordListDAG
	// MaxValue is the maximum value among children not counting multipliers
	MaxValue int
}
