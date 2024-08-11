package ranker

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/sepcon/quizprob/internal/common_errors"
	"github.com/sepcon/quizprob/pkg/model/quiz"
	"sort"
	"sync"
)

type userScoreMap = map[quiz.UserIDType]quiz.ScoreType
type quizd struct {
	userScores      userScoreMap
	leaderboardSize uint32
}
type quizScoresMap = map[quiz.QuizIDType]*quizd
type leaderboardMap = map[quiz.QuizIDType]treemap.Map
type inMemoryRanker struct {
	quiz         quizScoresMap
	leaderboards leaderboardMap
	mu           sync.Mutex
}

func NewInMemoryRanker() Ranker {
	return newInMemoryRanker()
}

func newQuizd(leaderboardSize uint32) *quizd {
	return &quizd{
		userScores:      make(userScoreMap),
		leaderboardSize: leaderboardSize,
	}
}
func newInMemoryRanker() *inMemoryRanker {
	return &inMemoryRanker{
		quiz:         make(quizScoresMap),
		leaderboards: make(leaderboardMap),
	}
}

func (r *inMemoryRanker) UpdateUserScore(quizUser quiz.QuizUser,
	scoreDelta quiz.ScoreType) (quiz.ScoreType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	latestScore, err := r.getUserScoreUnsafe(quizUser)
	if err != nil {
		return quiz.ScoreType(0), err
	}

	latestScore += scoreDelta
	return latestScore, r.setUserScoreUnsafe(quizUser, latestScore)
}

func (r *inMemoryRanker) SetUserScore(quizUser quiz.QuizUser, score quiz.ScoreType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.setUserScoreUnsafe(quizUser, score)
}

func (r *inMemoryRanker) GetUserScore(quizUser quiz.QuizUser) (quiz.ScoreType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.getUserScoreUnsafe(quizUser)
}

func (r *inMemoryRanker) GetUserRank(quizUser quiz.QuizUser) (quiz.RankType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.getUserRankUnsafe(quizUser)
}

func (r *inMemoryRanker) GetUserStatus(quizUser quiz.QuizUser) (*quiz.UserStatus, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	score, err := r.getUserScoreUnsafe(quizUser)
	if err != nil {
		return nil, err
	}
	rank, err := r.getUserRankUnsafe(quizUser)
	if err != nil {
		return nil, err
	}

	return &quiz.UserStatus{
		User:  quiz.User{UserID: quizUser.UserID},
		Score: score,
		Rank:  rank,
	}, nil
}

func (r *inMemoryRanker) GetLeaderboard(quizid quiz.QuizIDType) (*quiz.Leaderboard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	quizd, exists := r.quiz[quizid]
	if !exists {
		return nil, common_errors.NewQuizNotFoundError(quizid)
	}

	leaders := make([]quiz.UserStatus, 0, len(r.quiz))
	for userID, score := range quizd.userScores {
		leaders = append(leaders, quiz.UserStatus{
			User:  quiz.User{userID},
			Score: score,
		})
	}

	sort.Slice(leaders, func(i, j int) bool {
		return leaders[i].Score > leaders[j].Score
	})

	topN := quizd.leaderboardSize
	if topN > uint32(len(leaders)) {
		topN = uint32(len(leaders))
	}
	leaders = leaders[:topN]
	// specify rank for leaders
	for i := range leaders {
		leaders[i].Rank = quiz.RankType(i + 1)
	}
	return &quiz.Leaderboard{
		Quiz:    quiz.Quiz{quizid},
		Leaders: leaders[:topN],
	}, nil
}

func (r *inMemoryRanker) CreateQuiz(quizID quiz.QuizIDType, leaderboardSize uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exist := r.quiz[quizID]; exist {
		return common_errors.NewQuizAlreadyExistsError(quizID)
	}

	r.quiz[quizID] = newQuizd(leaderboardSize)
	return nil
}

func (r *inMemoryRanker) setUserScoreUnsafe(quizUser quiz.QuizUser, score quiz.ScoreType) error {
	if _, ok := r.quiz[quizUser.QuizID]; !ok {
		return common_errors.NewQuizNotFoundError(quizUser.QuizID)
	}

	r.quiz[quizUser.QuizID].userScores[quizUser.UserID] = score
	return nil
}

func (r *inMemoryRanker) getUserScoreUnsafe(quizUser quiz.QuizUser) (quiz.ScoreType, error) {
	if quizd, ok := r.quiz[quizUser.QuizID]; ok {
		if score, ok := quizd.userScores[quizUser.UserID]; !ok {
			return score, common_errors.NewUserNotFoundError(quizUser.UserID)
		} else {
			return score, nil
		}
	}
	return 0, common_errors.NewQuizNotFoundError(quizUser.QuizID)
}

func (r *inMemoryRanker) getUserRankUnsafe(quizUser quiz.QuizUser) (quiz.RankType, error) {
	quizd, exists := r.quiz[quizUser.QuizID]
	if !exists {
		return -1, common_errors.NewQuizNotFoundError(quizUser.QuizID)
	}

	score, exists := quizd.userScores[quizUser.UserID]
	if !exists {
		return -1, common_errors.NewUserNotFoundError(quizUser.UserID)
	}

	rank := 1
	for _, s := range quizd.userScores {
		if s > score {
			rank++
		}
	}
	return quiz.RankType(rank), nil
}
