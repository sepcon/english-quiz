package scoring

import (
	"errors"
	"github.com/sepcon/quizprob/internal/common_errors"
	rankermocks "github.com/sepcon/quizprob/internal/scoring/ranker/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupScoringManager(t *testing.T, scoringManagerOpts ...ScoringManagerOption) (ScoringManager, *rankermocks.Ranker) {
	ranker := rankermocks.NewRanker(t)
	sm := NewScoringManager(ranker, scoringManagerOpts...)
	return sm, ranker
}
func Test_scoringManagerImpl(t *testing.T) {
	sm, ranker := setupScoringManager(t)
	user := scoring_service.QuizUser{
		Quiz: scoring_service.Quiz{"quiz1"},
		User: scoring_service.User{"user1"},
	}
	t.Run("UpdateOrInitUserStatus", func(t *testing.T) {
		t.Run("Status exists on ranker then return status", func(t *testing.T) {
			statusRequest := scoring_service.UserStatusRequest{}
			statusRequest.UserID = "userid"
			statusRequest.QuizID = "quizid"
			expectedUserStatus := scoring_service.UserStatus{
				User:  scoring_service.User{statusRequest.UserID},
				Score: 1000,
				Rank:  1000,
			}
			ranker.On("GetUserStatus", statusRequest.QuizUser).Return(&expectedUserStatus, nil)

			status, err := sm.GetOrInitUserStatus(statusRequest)
			assert.NoError(t, err)
			assert.Equal(t, &expectedUserStatus, status)
		})
		t.Run("User Not exist", func(t *testing.T) {
			expectedErr := errors.New("Failed to GetUserStatus")

			ranker.On("GetUserStatus", user).Return(nil, expectedErr).Once()

			t.Run("SetUserScoreSuccessfully", func(t *testing.T) {
				expectedUserStatus := scoring_service.UserStatus{
					User:  user.User,
					Score: 0,
					Rank:  -1,
				}

				ranker.On("GetUserStatus", user).Return(&expectedUserStatus, nil).Once()
				ranker.On("SetUserScore", user, scoring_service.ScoreType(0)).Return(nil)
				status, err := sm.GetOrInitUserStatus(scoring_service.UserStatusRequest{user})
				assert.Equal(t, expectedUserStatus, *status)
				assert.NoError(t, err)
			})
			t.Run("SetUserScoreFailed", func(t *testing.T) {
				expectedErr = common_errors.NewQuizNotFoundError(user.QuizID)
				ranker.On("SetUserScore", user, scoring_service.ScoreType(0)).Return(expectedErr)
				status, err := sm.GetOrInitUserStatus(scoring_service.UserStatusRequest{user})
				assert.Nil(t, status)
				assert.True(t, common_errors.Match(err, common_errors.EC_QUIZ_NOTFOUND))
			})
		})

	})

	t.Run("UpdateUserScore", func(t *testing.T) {

		const scoreDelta = scoring_service.ScoreType(10)
		const expectedScore = scoring_service.ScoreType(110)
		const zeroScore = scoring_service.ScoreType(0)
		const expectedRank = scoring_service.RankType(200)
		expectedLeaderboard := new(scoring_service.Leaderboard)
		userScoreUpdateCallbackCalledCount := 0
		leaderboardUpdatedCallbackCalledCount := 0
		sm, ranker := setupScoringManager(t,
			OptLeaderboardUpdateCallback(func(leaderboard scoring_service.Leaderboard) {
				assert.Equal(t, *expectedLeaderboard, leaderboard)
				leaderboardUpdatedCallbackCalledCount++
			}),
			OptUserScoreUpdateCallback(func(quizUser scoring_service.QuizUser, score scoring_service.ScoreType, rank scoring_service.RankType) {
				assert.Equal(t, quizUser, quizUser)
				assert.Equal(t, expectedScore, score)
				assert.Equal(t, expectedRank, rank)
				userScoreUpdateCallbackCalledCount++
			}),
		)

		t.Run("Ranker.UpdateUserScoreFailed", func(t *testing.T) {
			expectedErr := errors.New("UpdateScoreFailed")
			ranker.On("UpdateUserScore", user, scoreDelta).Return(zeroScore, expectedErr)
			sm.UpdateUserScore(scoring_service.UpdateScoreRequest{
				QuizUser:   user,
				ScoreDelta: scoreDelta,
			})
			assert.Zero(t, userScoreUpdateCallbackCalledCount)
			assert.Zero(t, leaderboardUpdatedCallbackCalledCount)
		})
		t.Run("Ranker.UpdateUserScoreSucceeded", func(t *testing.T) {
			ranker.On("UpdateUserScore", user, scoreDelta).Return(expectedScore, nil)
			ranker.On("GetUserRank", user).Return(expectedRank, nil)
			ranker.On("GetLeaderboard", user.QuizID).Return(expectedLeaderboard, nil)
			sm.UpdateUserScore(scoring_service.UpdateScoreRequest{
				QuizUser:   user,
				ScoreDelta: scoreDelta,
			})
			assert.Equal(t, 1, userScoreUpdateCallbackCalledCount)
			assert.Equal(t, 1, leaderboardUpdatedCallbackCalledCount)
		})
	})
}
