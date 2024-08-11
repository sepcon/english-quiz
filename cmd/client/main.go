package main

import (
	"flag"
	"fmt"
	service_shared "github.com/sepcon/quizprob/cmd/shared"
	"github.com/sepcon/quizprob/pkg/event_sdk"
	"os"
)

func main() {
	userID := ""
	flag.StringVar(&userID, "userid", "", "unique userid to join the quiz")
	flag.Parse()
	if userID != "" {
		subscriber := event_sdk.NewSubscriber("localhost:"+service_shared.CONST_EVENT_SERVICE_PORT, userID)
		NewWorkflow(service_shared.CONST_QUIZ_GAME_ID, userID, subscriber).Play()
	} else {
		fmt.Println("Error: userid must not be empty and must be unique")
		os.Exit(1)
	}

}
