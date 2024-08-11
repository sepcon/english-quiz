package ranker

import "github.com/sepcon/quizprob/pkg/model/quiz"

type Ranker interface {
	UpdateUserScore(quizUser quiz.QuizUser, scoreDelta quiz.ScoreType) (quiz.ScoreType, error)
	SetUserScore(quizUser quiz.QuizUser, score quiz.ScoreType) error
	GetUserScore(quizUser quiz.QuizUser) (quiz.ScoreType, error)
	GetUserRank(quizUser quiz.QuizUser) (quiz.RankType, error)
	GetUserStatus(quizUser quiz.QuizUser) (*quiz.UserStatus, error)
	GetLeaderboard(quizid quiz.QuizIDType) (*quiz.Leaderboard, error)
	CreateQuiz(quizID quiz.QuizIDType, leaderboardSize uint32) error
}
