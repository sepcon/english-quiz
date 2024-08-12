package main

import (
	service_shared "github.com/sepcon/quizprob/cmd/shared"
	"github.com/sepcon/quizprob/pkg/scoring_sdk"
)

func main() {
	NewQuizServer(":"+service_shared.CONST_QUIZ_SERVICE_PORT,
		service_shared.CONST_QUIZ_GAME_ID, 10000,
		scoring_sdk.NewScoringClient("http://localhost:"+service_shared.CONST_SCORING_SERVICE_PORT),
	).Start()

}
