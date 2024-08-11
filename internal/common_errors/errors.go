package common_errors

import (
	"fmt"
	"github.com/sepcon/quizprob/pkg/model/quiz"
)

type QuizErrorCode int

const (
	EC_QUIZ_NOTFOUND QuizErrorCode = iota
	EC_USER_NOTFOUND
	EC_QUIZ_ALREADY_EXIST
)

type QuizError struct {
	error
	Code QuizErrorCode
}

func Match(err error, code QuizErrorCode) bool {
	quizErr, ok := err.(*QuizError)
	if ok {
		return quizErr.Code == code
	}
	return false
}

func NewQuizNotFoundError(quizID quiz.QuizIDType) error {
	return &QuizError{
		error: fmt.Errorf("Quiz %v not found", quizID),
		Code:  EC_QUIZ_NOTFOUND,
	}
}
func NewQuizAlreadyExistsError(quizID quiz.QuizIDType) error {
	return &QuizError{
		error: fmt.Errorf("Quiz %v already exists", quizID),
		Code:  EC_QUIZ_ALREADY_EXIST,
	}
}
func NewUserNotFoundError(userID quiz.UserIDType) error {
	return &QuizError{
		error: fmt.Errorf("User %v not found", userID),
		Code:  EC_USER_NOTFOUND,
	}
}
