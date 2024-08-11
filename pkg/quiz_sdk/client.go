package quiz_sdk

import (
	"context"
	"github.com/sepcon/quizprob/pkg/model/quiz_service"
)

type Client interface {
	JoinQuiz(ctx context.Context, req *quiz_service.JoinQuizRequest) (*quiz_service.JoinQuizResponse, error)
	SubmitAnswer(ctx context.Context, req *quiz_service.SubmitAnswerRequest) (*quiz_service.SubmitAnswerResponse, error)
}
