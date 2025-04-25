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

type graphqlTodayScreenResponse struct {
	TodayPage struct {
		Daily struct {
			Puzzles []struct {
				UrlPath string
				Puzzle  struct {
					Puzzle string
					Game   struct {
						Slug string
					}
				}
			}
		}
	}
}
