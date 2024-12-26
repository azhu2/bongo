package gameimporter

type graphqlBoardResponse struct {
	StartOrFindGameplay struct {
		GamePlayed struct {
			Puzzle struct {
				Puzzle string
			}
		}
	}
}
