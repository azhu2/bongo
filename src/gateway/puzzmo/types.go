package puzzmo

type bongoResponse struct {
	StartOrFindGameplay struct {
		GamePlayed struct {
			Puzzle struct {
				Puzzle string
			}
		}
	}
}
