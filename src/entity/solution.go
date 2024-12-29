package entity

type Solution []rune

func (s Solution) Get(i, j int) rune {
	return s[i*BoardSize+j]
}

func (s Solution) GetRow(i int) []rune {
	return s[i*BoardSize : (i+1)*BoardSize]
}

func (s Solution) Set(i, j int, letter rune) {
	s[i*BoardSize+j] = letter
}

func (s Solution) SetRow(i int, word []rune) {
	for j, letter := range word {
		s.Set(i, j, letter)
	}
}

func (s Solution) Rows() [][]rune {
	rows := make([][]rune, BoardSize)
	for i := 0; i < BoardSize; i++ {
		rows[i] = s[i*BoardSize : (i+1)*BoardSize]
	}
	return rows
}

func (s Solution) String() string {
	ret := ""
	for _, row := range s.Rows() {
		ret += string(row) + "|"
	}
	return ret
}

func EmptySolution() Solution {
	empty := make([]rune, BoardSize*BoardSize)
	for i := 0; i < BoardSize*BoardSize; i++ {
		empty[i] = ' '
	}
	return empty
}
