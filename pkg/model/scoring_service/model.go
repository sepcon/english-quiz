package scoring_service

import "github.com/sepcon/quizprob/pkg/model/quiz"

type UserStatusRequest struct {
	quiz.QuizUser
}

type UserStatusResponse struct {
	quiz.Quiz
	quiz.UserStatus
}

type UpdateScoreRequest struct {
	quiz.QuizUser
	ScoreDelta quiz.ScoreType `json:"delta"`
}

type UpdateScoreResponse struct {
	Accepted bool `json:"accepted"`
}

type CreateQuizRequest struct {
	quiz.Quiz
	LeaderboardSize uint32 `json:"leaderboard_size"`
}

type LeaderBoardRequest struct {
	quiz.Quiz
}

type LeaderboardResponse = quiz.Leaderboard
