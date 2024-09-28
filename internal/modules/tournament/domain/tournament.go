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
	minIdDigits         = 3
)

type Tournament struct {
	ID        string    `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	// Pre-requisite: Players creation is done in player-couple module.
	// Tournament registration is done here:
	PlayerCouples []PlayerCouple `bson:"player_couples,omitempty" json:"player_couples,omitempty"`
	Rounds        []Round        `bson:"rounds,omitempty" json:"rounds,omitempty"`
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
	Number  int     `bson:"number" json:"number"`
	Matches []Match `bson:"matches" json:"matches"`
}

type Match struct {
	ID        string       `bson:"_id" json:"id"`
	Timestamp time.Time    `bson:"timestamp" json:"timestamp"`
	Couple1   PlayerCouple `bson:"couple1" json:"couple1"`
	Couple2   PlayerCouple `bson:"couple2" json:"couple2"`
	Score     *Score       `bson:"score,omitempty" json:"score,omitempty"`
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
	Set1 GameSet  `bson:"set1" json:"set1"`
	Set2 GameSet  `bson:"set2" json:"set2"`
	Set3 *GameSet `bson:"set3,omitempty" json:"set3,omitempty"`
}

type GameSet struct {
	GamesCouple1 int       `bson:"gamesCouple1" json:"gamesCouple1"`
	GamesCouple2 int       `bson:"gamesCouple2" json:"gamesCouple2"`
	Tiebreak     *Tiebreak `bson:"tiebreak,omitempty" json:"tiebreak,omitempty"`
}

type Tiebreak struct {
	PointsCouple1 int `bson:"pointsCouple1" json:"pointsCouple1"`
	PointsCouple2 int `bson:"pointsCouple2" json:"pointsCouple2"`
}

func NewTournament(title string, timestamp time.Time,
	playerCouples []PlayerCouple, rounds []Round) (*Tournament, error) {
	if len(title) < minTournamentDigits {
		return nil, fmt.Errorf("invalid title: %s", title)
	}
	if timestamp.Before(time.Now().AddDate(0, 0, minMatchDays)) {
		return nil, fmt.Errorf("timestamp cannot older than %d days", minMatchDays)
	}
	return &Tournament{
		//ID:          auto generated ID set in the repository.
		Title:         title,
		Timestamp:     timestamp,
		PlayerCouples: playerCouples,
		Rounds:        rounds,
	}, nil
}

func NewRound(roundNumber int, matches []Match) (*Round, error) {
	// Round is an embedded struct inside a root aggregate that will be stored in the tournament repository.
	if roundNumber < minRoundNumber {
		return nil, fmt.Errorf("invalid roundNumber: %d", roundNumber)
	}
	return &Round{
		Number:  roundNumber,
		Matches: matches,
	}, nil
}

func NewMatch(id string, timestamp time.Time, couple1, couple2 PlayerCouple, score *Score) (*Match, error) {
	// Match is an embedded struct inside a root aggregate that will be stored in the tournament repository.
	if len(id) < minIdDigits {
		return nil, fmt.Errorf("invalid id: %s", id)
	}
	if timestamp.Before(time.Now().AddDate(0, 0, minMatchDays)) {
		return nil, fmt.Errorf("timestamp cannot older than %d days", minMatchDays)
	}
	if couple1.ID == couple2.ID {
		return nil, fmt.Errorf("couple1 and couple2 cannot be the same")
	}
	return &Match{
		ID:        id,
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
