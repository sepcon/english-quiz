package quiz_service

import "github.com/sepcon/quizprob/pkg/model/quiz"

type QuestionIDType = string
type ChoiceType = string

type Question struct {
	ID      QuestionIDType        `json:"id,omitempty"`
	Ask     string                `json:"ask,omitempty"`
	Options map[ChoiceType]string `json:"options,omitempty"`
}

type JoinQuizRequest struct {
	quiz.QuizUser
}

type JoinQuizResponse struct {
	quiz.UserStatus
	Leaderboard []quiz.UserStatus `json:"leaderboard"`
	Questions   []Question        `json:"questions"`
}

type SubmitAnswerRequest struct {
	quiz.QuizUser
	QuestionID QuestionIDType `json:"question_id,omitempty"`
	Choice     ChoiceType     `json:"choice,omitempty"`
}

type SubmitAnswerResponse struct {
	Correct         bool           `json:"correct"`
	Message         string         `json:"message,omitempty"`
	AdditionalScore quiz.ScoreType `json:"additional_score"`
}
