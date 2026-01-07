package models

import (
	"errors"
	"fmt"
	"strings"
)

type ProblemDifficulty int

const (
	ProblemEasy ProblemDifficulty = iota + 1
	ProblemMedium
	ProblemHard
)

var (
	ErrInvalidProblemDifficulty = errors.New("invalid problem level")
)

func (p ProblemDifficulty) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", p.String()), nil
}

func NewProblemDifficulty(level string) (ProblemDifficulty, error) {
	switch strings.ToLower(level) {
	case "easy":
		return ProblemEasy, nil
	case "medium":
		return ProblemMedium, nil
	case "hard":
		return ProblemHard, nil
	default:
		return 0, ErrInvalidProblemDifficulty
	}
}

func (p ProblemDifficulty) String() string {
	switch p {
	case ProblemEasy:
		return "Easy"
	case ProblemMedium:
		return "Medium"
	case ProblemHard:
		return "Hard"
	default:
		return "Unknown"
	}
}
