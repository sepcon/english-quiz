package ranker

import (
	"github.com/sepcon/quizprob/internal/common_errors"
	"github.com/sepcon/quizprob/pkg/model/quiz"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type quizHelper struct {
	quizID quiz.QuizIDType
}

func Test_inMemoryRanker(t *testing.T) {
	ranker := newInMemoryRanker()
	const quiz1 = "quiz1"
	t.Run("Quiz not exist", func(t *testing.T) {
		score, err := ranker.GetUserScore(quiz.QuizUser{
			Quiz: quiz.Quiz{quiz1},
			User: quiz.User{"usr1"},
		})
		assert.True(t, common_errors.Match(err, common_errors.EC_QUIZ_NOTFOUND))
		assert.Equal(t, quiz.ScoreType(0), score)
	})
	t.Run("Quiz exist", func(t *testing.T) {
		assert.NoError(t, ranker.CreateQuiz(quiz1, 2))
		quizUser := quiz.NewQuizUser(quiz1, "usr1")
		t.Run("Query  non exist user's score must be failed", func(t *testing.T) {
			score, err := ranker.GetUserScore(quizUser)
			assert.True(t, common_errors.Match(err, common_errors.EC_USER_NOTFOUND))
			assert.Equal(t, quiz.ScoreType(0), score)
		})
		t.Run("Update user score and get back to confirm", func(t *testing.T) {
			const (
				expectedScore quiz.ScoreType = 100
				scoreDelta    quiz.ScoreType = 20
			)
			ranker.SetUserScore(quizUser, expectedScore)
			score, err := ranker.GetUserScore(quizUser)
			assert.NoError(t, err)
			assert.Equal(t, expectedScore, score)

			score, err = ranker.UpdateUserScore(quizUser, scoreDelta)
			assert.NoError(t, err)
			assert.Equal(t, expectedScore+scoreDelta, score)
			score, err = ranker.GetUserScore(quizUser)
			assert.NoError(t, err)
			assert.Equal(t, expectedScore+scoreDelta, score)
		})

		t.Run("Set multi user and get leader board", func(t *testing.T) {
			userScores := []struct {
				user         quiz.QuizUser
				score        quiz.ScoreType
				scoreDelta   quiz.ScoreType
				expectedRank quiz.RankType
			}{
				{
					user:         quiz.NewQuizUser(quiz1, "usr1"),
					score:        100,
					scoreDelta:   10,
					expectedRank: 4,
				}, {
					user:         quiz.NewQuizUser(quiz1, "usr2"),
					score:        100,
					scoreDelta:   50,
					expectedRank: 2,
				}, {
					user:         quiz.NewQuizUser(quiz1, "usr3"),
					score:        100,
					scoreDelta:   20,
					expectedRank: 3,
				}, {
					user:         quiz.NewQuizUser(quiz1, "usr4"),
					score:        100,
					scoreDelta:   60,
					expectedRank: 1,
				}, {
					user:         quiz.NewQuizUser(quiz1, "usr5"),
					score:        100,
					scoreDelta:   5,
					expectedRank: 5,
				},
			}

			// set initial score for user first
			for _, us := range userScores {
				assert.NoError(t, ranker.SetUserScore(us.user, us.score))
			}
			// update new score
			for _, us := range userScores {
				score, err := ranker.UpdateUserScore(us.user, us.scoreDelta)
				assert.NoError(t, err)
				assert.Equal(t, us.score+us.scoreDelta, score)
			}

			leaders, err := ranker.GetLeaderboard(quiz1)
			assert.NoError(t, err)
			assert.NotNil(t, leaders)

			assert.True(t, reflect.DeepEqual(&quiz.Leaderboard{
				Quiz: quiz.Quiz{quiz1},
				Leaders: []quiz.UserStatus{
					{User: quiz.User{"usr4"}, Score: 160, Rank: 1},
					{User: quiz.User{"usr2"}, Score: 150, Rank: 2},
				},
			}, leaders))

			for _, us := range userScores {
				rank, err := ranker.GetUserRank(us.user)
				assert.NoError(t, err)
				assert.Equal(t, us.expectedRank, rank,
					"user %v must with score %v must have rank %v", us.user.UserID, us.scoreDelta+us.scoreDelta, us.expectedRank)
			}
		})

	})
}
