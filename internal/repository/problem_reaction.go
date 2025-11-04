package repository

import (
	"context"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/google/uuid"
)

type ProblemReactionRepository interface {
	Create(ctx context.Context, userId, problemId uuid.UUID, reaction models.ProblemReactionType) error

	GetByUserAndProblem(ctx context.Context, userId, problemId uuid.UUID) (models.ProblemReactionType, error)
	Update(ctx context.Context, userId, problemId uuid.UUID, reaction models.ProblemReactionType) error

	GetLikesCountForProblem(ctx context.Context, problemId uuid.UUID) (int, error)
	GetDislikesCountForProblem(ctx context.Context, problemId uuid.UUID) (int, error)
	HasUserReacted(ctx context.Context, userId, problemId uuid.UUID) (reaction models.ProblemReactionType, err error)
}
