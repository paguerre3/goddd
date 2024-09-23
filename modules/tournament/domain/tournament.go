package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	minGames            = 0
	maxGames            = 19
	minBreakpoints      = 0
	maxBreakpoints      = 21
	minMatchDays        = -30
	minRoundNumber      = 0
	noSecondsFormat     = "2006-01-02T15:04"
	minTournamentDigits = 5
)

type Tournament struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Timestamp time.Time `json:"timestamp"`
	// Pre-requisite: Players creation is done in player-couple module.
	// Tournament registration is done here:
	PlayerCouples []PlayerCouple `json:"player_couples,omitempty"`
	Rounds        []Round        `json:"rounds,omitempty"`
}

// Custom JSON marshalling to format time without seconds:
func (t Tournament) MarshalJSON() ([]byte, error) {
	type Alias Tournament
	return json.Marshal(&struct {
		Timestamp string `json:"timestamp"`
		Alias
	}{
		Timestamp: t.Timestamp.Format(noSecondsFormat),
		Alias:     (Alias)(t),
	})
}

type Round struct {
	Number  int     `json:"number"`
	Matches []Match `json:"matches"`
}

type Match struct {
	ID        string       `json:"id"`
	Timestamp time.Time    `json:"timestamp"`
	Couple1   PlayerCouple `json:"couple1"`
	Couple2   PlayerCouple `json:"couple2"`
	Score     *Score       `json:"score,omitempty"`
}

// Custom JSON marshalling to format time without seconds:
func (m Match) MarshalJSON() ([]byte, error) {
	type Alias Match
	return json.Marshal(&struct {
		Timestamp string `json:"timestamp"`
		Alias
	}{
		Timestamp: m.Timestamp.Format(noSecondsFormat),
		Alias:     (Alias)(m),
	})
}

type Score struct {
	Set1 GameSet  `json:"set1"`
	Set2 GameSet  `json:"set2"`
	Set3 *GameSet `json:"set3,omitempty"`
}

type GameSet struct {
	GamesCouple1 int       `json:"gamesCouple1"`
	GamesCouple2 int       `json:"gamesCouple2"`
	Tiebreak     *Tiebreak `json:"tiebreak,omitempty"`
}

type Tiebreak struct {
	PointsCouple1 int `json:"pointsCouple1"`
	PointsCouple2 int `json:"pointsCouple2"`
}

func NewTournament(idGen IDGenerator, title string, timestamp time.Time,
	playerCouples []PlayerCouple, rounds []Round) (*Tournament, error) {
	if idGen == nil {
		return nil, fmt.Errorf("idGen cannot be nil")
	}
	if len(title) < minTournamentDigits {
		return nil, fmt.Errorf("invalid title: %s", title)
	}
	if timestamp.Before(time.Now().AddDate(0, 0, minMatchDays)) {
		return nil, fmt.Errorf("timestamp cannot older than %d days", minMatchDays)
	}
	return &Tournament{
		ID:            idGen.GenerateID(),
		Title:         title,
		Timestamp:     timestamp,
		PlayerCouples: playerCouples,
		Rounds:        rounds,
	}, nil
}

func NewRound(roundNumber int, matches []Match) (*Round, error) {
	if roundNumber < minRoundNumber {
		return nil, fmt.Errorf("invalid roundNumber: %d", roundNumber)
	}
	return &Round{
		Number:  roundNumber,
		Matches: matches,
	}, nil
}

func NewMatch(idGen IDGenerator, timestamp time.Time, couple1, couple2 PlayerCouple, score *Score) (*Match, error) {
	if idGen == nil {
		return nil, fmt.Errorf("idGen cannot be nil")
	}
	if timestamp.Before(time.Now().AddDate(0, 0, minMatchDays)) {
		return nil, fmt.Errorf("timestamp cannot older than %d days", minMatchDays)
	}
	if couple1.ID == couple2.ID {
		return nil, fmt.Errorf("couple1 and couple2 cannot be the same")
	}
	return &Match{
		ID:        idGen.GenerateID(),
		Timestamp: timestamp,
		Couple1:   couple1,
		Couple2:   couple2,
		Score:     score,
	}, nil
}

func NewScore(set1, set2, set3 *GameSet) (*Score, error) {
	if set1 == nil {
		return nil, fmt.Errorf("set1 cannot be nil")
	}
	if set2 == nil {
		return nil, fmt.Errorf("set2 cannot be nil")
	}
	return &Score{Set1: *set1, Set2: *set2, Set3: set3}, nil
}

func NewGameSet(gamesCouple1, gamesCouple2 int, tiebreak *Tiebreak) (*GameSet, error) {
	if gamesCouple1 < minGames || gamesCouple1 > maxGames {
		return nil, fmt.Errorf("invalid gamesCouple1: %d", gamesCouple1)
	}
	if gamesCouple2 < minGames || gamesCouple2 > maxGames {
		return nil, fmt.Errorf("invalid gamesCouple2: %d", gamesCouple2)
	}
	return &GameSet{
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
