package scoring_sdk

import (
	"context"
	"github.com/sepcon/quizprob/pkg/model/quiz"
	"github.com/sepcon/quizprob/pkg/model/scoring_service"
)

type Client interface {
	CreateQuiz(ctx context.Context, req *scoring_service.CreateQuizRequest) error
	UpdateUserScore(ctx context.Context, req *scoring_service.UpdateScoreRequest) error
	GetLeaderboard(ctx context.Context, quizID quiz.QuizIDType) (*quiz.Leaderboard, error)
	GetOrInitUserStatus(ctx context.Context, req *scoring_service.UserStatusRequest) (*quiz.UserStatus, error)
}
