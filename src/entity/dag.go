package entity

type WordListDAG struct {
	Children map[rune]WordListDAG
	MaxValue int
	Depth    int
}
