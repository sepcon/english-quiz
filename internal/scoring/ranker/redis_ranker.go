package ranker

import (
	scoringconf "github.com/sepcon/quizprob/config/scoring"
)

func NewRedisRanker(conf scoringconf.Config) (Ranker, error) {
	// for moking purpose, we will use inMemoryRanker
	return NewInMemoryRanker(), nil
}
