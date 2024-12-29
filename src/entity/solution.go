package entity

type Solution [][]rune

func EmptySolution() Solution {
	empty := make([][]rune, BoardSize)
	for i := 0; i < BoardSize; i++ {
		row := make([]rune, BoardSize)
		for j := 0; j < BoardSize; j++ {
			row[j] = ' '
		}
		empty[i] = row
	}
	return empty
}

func (s Solution) String() string {
	ret := ""
	for _, row := range s {
		ret += string(row) + "|"
	}
	return ret
}
