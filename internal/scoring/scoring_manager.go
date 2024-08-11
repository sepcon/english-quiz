package scoring

import (
	"github.com/sepcon/quizprob/internal/scoring/ranker"
	"github.com/sepcon/quizprob/internal/scoring/user_score_db"
	"github.com/sepcon/quizprob/pkg/model/quiz"
	"github.com/sepcon/quizprob/pkg/model/scoring_service"
)

type LeaderboardUpdateCallback = func(leaderboard quiz.Leaderboard)
type UserStatusUpdateCallback = func(quizID quiz.QuizIDType, status *quiz.UserStatus)

type ScoringManager interface {
	CreateQuiz(request scoring_service.CreateQuizRequest) error
	UpdateUserScore(request scoring_service.UpdateScoreRequest) error
	GetLeaderboard(request scoring_service.LeaderBoardRequest) (*quiz.Leaderboard, error)
	GetOrInitUserStatus(request scoring_service.UserStatusRequest) (*quiz.UserStatus, error)
}

func NewScoringManager(ranker ranker.Ranker, options ...ScoringManagerOption) ScoringManager {
	impl := &scoringManagerImpl{ranker: ranker}
	for _, option := range options {
		option(impl)
	}
	return impl
}

type ScoringManagerOption = func(*scoringManagerImpl)

func OptPersistentUserScoreDB(persistentUserScoreDB user_score_db.UserScoreDB) ScoringManagerOption {
	return func(impl *scoringManagerImpl) {
		impl.persistentUserScoreDB = persistentUserScoreDB
	}
}
func OptLeaderboardUpdateCallback(cb LeaderboardUpdateCallback) ScoringManagerOption {
	return func(impl *scoringManagerImpl) {
		impl.leaderboardCallbacks = append(impl.leaderboardCallbacks, cb)
	}
}

func OptUserScoreUpdateCallback(cb UserStatusUpdateCallback) ScoringManagerOption {
	return func(impl *scoringManagerImpl) {
		impl.userStatusUpdateCallbacks = append(impl.userStatusUpdateCallbacks, cb)
	}
}

type scoringManagerImpl struct {
	ranker                    ranker.Ranker
	persistentUserScoreDB     user_score_db.UserScoreDB
	leaderboardCallbacks      []LeaderboardUpdateCallback
	userStatusUpdateCallbacks []UserStatusUpdateCallback
}

func (sm *scoringManagerImpl) CreateQuiz(request scoring_service.CreateQuizRequest) error {
	if err := sm.ranker.CreateQuiz(request.QuizID, request.LeaderboardSize); err != nil {
		return err
	}
	if sm.persistentUserScoreDB != nil {
		// NOTE: Must consider the case of undo the ranker as well
		return sm.persistentUserScoreDB.CreateQuiz(request.QuizID)
	}
	return nil
}

func (sm *scoringManagerImpl) UpdateUserScore(request scoring_service.UpdateScoreRequest) error {
	userScore, err := sm.ranker.UpdateUserScore(request.QuizUser, request.ScoreDelta)
	if err != nil {
		return err
	}

	if sm.persistentUserScoreDB != nil {
		err = sm.persistentUserScoreDB.UpdateUserScore(request.QuizUser, userScore)
		if err != nil {
			return err
		}
	}

	userRank, err := sm.ranker.GetUserRank(request.QuizUser)
	if err != nil {
		return err
	}

	sm.notifyUserScoreUpdated(request.QuizUser.QuizID, &quiz.UserStatus{
		User:  quiz.User{UserID: request.QuizUser.UserID},
		Score: userScore,
		Rank:  userRank,
	})

	leaderboard, err := sm.ranker.GetLeaderboard(request.QuizID)
	if err != nil {
		return err
	}
	sm.notifyLeaderboardUpdated(*leaderboard)
	return nil
}

func (sm *scoringManagerImpl) GetLeaderboard(request scoring_service.LeaderBoardRequest) (*quiz.Leaderboard, error) {
	return sm.ranker.GetLeaderboard(request.QuizID)
}

func (sm *scoringManagerImpl) GetOrInitUserStatus(request scoring_service.UserStatusRequest) (*quiz.UserStatus, error) {
	status, err := sm.ranker.GetUserStatus(request.QuizUser)
	if err == nil {
		return status, nil
	}
	var score quiz.ScoreType = 0
	if sm.persistentUserScoreDB != nil {
		score, err = sm.persistentUserScoreDB.GetOrInitUserStatus(request.QuizUser)
		if err != nil {
			return nil, err
		}
	}

	err = sm.ranker.SetUserScore(request.QuizUser, score)
	if err != nil {
		return nil, err
	}

	return sm.ranker.GetUserStatus(request.QuizUser)
}

func (sm *scoringManagerImpl) notifyUserScoreUpdated(quizID quiz.QuizIDType, status *quiz.UserStatus) {
	for _, callback := range sm.userStatusUpdateCallbacks {
		callback(quizID, status)
	}
}

func (sm *scoringManagerImpl) notifyLeaderboardUpdated(leaderboard quiz.Leaderboard) {
	for _, callback := range sm.leaderboardCallbacks {
		callback(leaderboard)
	}
}
