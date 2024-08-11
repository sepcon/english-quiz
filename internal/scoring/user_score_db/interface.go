package user_score_db

import "github.com/sepcon/quizprob/pkg/model/quiz"

type UserScoreDB interface {
	CreateQuiz(quizID quiz.QuizIDType) error
	UpdateUserScore(quizUser quiz.QuizUser, score quiz.ScoreType) error
	GetOrInitUserStatus(quizUser quiz.QuizUser) (quiz.ScoreType, error)
}
