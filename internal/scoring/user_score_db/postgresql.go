package user_score_db

import (
	scoringconf "github.com/sepcon/quizprob/config/scoring"
	"github.com/sepcon/quizprob/pkg/model/quiz"
)

func NewPostGreSQL(conf *scoringconf.PersistentStorage) (UserScoreDB, error) {
	return &postGreSQL{}, nil
}

type postGreSQL struct {
}

func (p *postGreSQL) CreateQuiz(quizID quiz.QuizIDType) error {
	//TODO implement me
	panic("implement me")
}

func (p *postGreSQL) UpdateUserScore(quizUser quiz.QuizUser, score quiz.ScoreType) error {
	//TODO implement me
	panic("implement me")
}

func (p *postGreSQL) GetOrInitUserStatus(quizUser quiz.QuizUser) (quiz.ScoreType, error) {
	//TODO implement me
	panic("implement me")
}
