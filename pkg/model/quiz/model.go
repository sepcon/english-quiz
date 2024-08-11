package quiz

type QuizIDType = string
type UserIDType = string
type ScoreType = int32
type RankType int32

type Quiz struct {
	QuizID QuizIDType `json:"quiz_id"`
}

type User struct {
	UserID UserIDType `json:"userid"`
}

type QuizUser struct {
	Quiz
	User
}

func NewQuizUser(quizID QuizIDType, userID UserIDType) QuizUser {
	return QuizUser{
		Quiz: Quiz{quizID},
		User: User{userID},
	}
}

type UserStatus struct {
	User
	Score ScoreType `json:"score"`
	Rank  RankType  `json:"rank"`
}

type Leaderboard struct {
	Quiz
	Leaders []UserStatus `json:"leaders"`
}
