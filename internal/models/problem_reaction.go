package models

type ProblemReactionType string

const (
	ProblemReactionLike    ProblemReactionType = "like"
	ProblemReactionDislike ProblemReactionType = "dislike"
	ProblemReactionNone    ProblemReactionType = "none"
)

func (r ProblemReactionType) IsValid() bool {
	switch r {
	case ProblemReactionLike, ProblemReactionDislike, ProblemReactionNone:
		return true
	default:
		return false
	}
}
