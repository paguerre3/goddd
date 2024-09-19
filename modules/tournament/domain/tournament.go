package domain

type Tournament struct {
	ID            string           `json:"id"`
	Title         string           `json:"title"`
	PlayerCouples []PlayerCouple   `json:"player_couples"`
	Rounds        []Round          `json:"rounds"`
	Scoreboard    map[string]Score `json:"scoreboard"`
}

type Round struct {
	ID      string  `json:"id"`
	Matches []Match `json:"matches"`
}

type Match struct {
	ID      string       `json:"id"`
	Couple1 PlayerCouple `json:"couple1"`
	Couple2 PlayerCouple `json:"couple2"`
	Score   Score        `json:"score"`
}

type Score struct {
	Set1 int  `json:"set1"`
	Set2 int  `json:"set2"`
	Set3 *int `json:"set3,omitempty"`
}
