package main

import (
	service_shared "github.com/sepcon/quizprob/cmd/shared"
	scoringconf "github.com/sepcon/quizprob/config/scoring"
	"github.com/sirupsen/logrus"
)

func main() {
	// Should load from environment variables
	server, err := NewScoringServer(scoringconf.Config{
		ServerPort: service_shared.CONST_SCORING_SERVICE_PORT,
		Concurrent: scoringconf.Concurrent{
			NumberOfWorker: 5,
			MaxQueueLength: 1000000,
		},
		Redis: scoringconf.Redis{
			Host: "",
			Port: 0,
		},
		Persistence:     nil,
		EventServiceUrl: "localhost:" + service_shared.CONST_EVENT_SERVICE_PORT,
	})
	if err != nil {
		logrus.Fatalf("Failed to setup server: %s", err.Error())
	}
	server.Serve(":" + service_shared.CONST_SCORING_SERVICE_PORT)
}
