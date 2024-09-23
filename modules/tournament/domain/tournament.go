package domain

import "fmt"

const (
	minGames       = 0
	maxGames       = 19
	minBreakpoints = 0
	maxBreakpoints = 21
)

type Tournament struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	PlayerCouples []PlayerCouple `json:"player_couples"`
	Rounds        []Round        `json:"rounds"`
	// map key is ROUND_ID:
	Scoreboard map[string]Score `json:"scoreboard"`
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
	Set1 GamesTuple  `json:"set1"`
	Set2 GamesTuple  `json:"set2"`
	Set3 *GamesTuple `json:"set3,omitempty"`
}

type GamesTuple struct {
	GamesCouple1 int       `json:"gamesCouple1"`
	GamesCouple2 int       `json:"gamesCouple2"`
	Tiebreak     *Tiebreak `json:"tiebreak,omitempty"`
}

type Tiebreak struct {
	PointsCouple1 int `json:"pointsCouple1"`
	PointsCouple2 int `json:"pointsCouple2"`
}

func NewScore(set1, set2, set3 *GamesTuple) (*Score, error) {
	if set1 == nil {
		return nil, fmt.Errorf("set1 cannot be nil")
	}
	if set2 == nil {
		return nil, fmt.Errorf("set2 cannot be nil")
	}
	return &Score{Set1: *set1, Set2: *set2, Set3: set3}, nil
}

func NewGamesTuple(gamesCouple1, gamesCouple2 int, tiebreak *Tiebreak) (*GamesTuple, error) {
	if gamesCouple1 < minGames || gamesCouple1 > maxGames {
		return nil, fmt.Errorf("invalid gamesCouple1: %d", gamesCouple1)
	}
	if gamesCouple2 < minGames || gamesCouple2 > maxGames {
		return nil, fmt.Errorf("invalid gamesCouple2: %d", gamesCouple2)
	}
	return &GamesTuple{
		GamesCouple1: gamesCouple1,
		GamesCouple2: gamesCouple2,
		Tiebreak:     tiebreak,
	}, nil
}

func NewTiebreak(pointsCouple1, pointsCouple2 int) (*Tiebreak, error) {
	if pointsCouple1 < minBreakpoints || pointsCouple1 > maxBreakpoints {
		return nil, fmt.Errorf("invalid pointsCouple1: %d", pointsCouple1)
	}
	if pointsCouple2 < minBreakpoints || pointsCouple2 > maxBreakpoints {
		return nil, fmt.Errorf("invalid pointsCouple2: %d", pointsCouple2)
	}
	return &Tiebreak{
		PointsCouple1: pointsCouple1,
		PointsCouple2: pointsCouple2,
	}, nil
}
