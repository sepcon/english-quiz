package main

import (
	"context"
	"encoding/json"
	"fmt"
	service_shared "github.com/sepcon/quizprob/cmd/shared"
	"github.com/sepcon/quizprob/pkg/event_sdk"
	"github.com/sepcon/quizprob/pkg/model/event_service"
	"github.com/sepcon/quizprob/pkg/model/quiz"
	"github.com/sepcon/quizprob/pkg/model/quiz_service"
	"github.com/sepcon/quizprob/pkg/quiz_sdk"
	"time"
)

type Workflow struct {
	quizUser            quiz.QuizUser
	userStatus          quiz.UserStatus
	quizClient          quiz_sdk.Client
	quizEventSubscriber event_sdk.Subscriber
}

type ignoreErrHandler struct {
}
type leaderboardUpdateHandler struct {
	ignoreErrHandler
}
type userStatusUpdateHandler struct {
	ignoreErrHandler
}

func (u ignoreErrHandler) OnConnectionError(err error) {
	fmt.Println("Got connection error ", err)
}

func (u *userStatusUpdateHandler) OnMessage(bytes []byte) {
	var userStatus quiz.UserStatus
	err := json.Unmarshal(bytes, &userStatus)
	if err != nil {
		fmt.Println("Error leaderboard update is not welformed!", string(bytes))
	} else {
		fmt.Printf("\nYou - Rank: %d, Score: %d", userStatus.Rank, userStatus.Score)
	}
}

func (q *leaderboardUpdateHandler) OnMessage(bytes []byte) {
	var leaderboard quiz.Leaderboard
	err := json.Unmarshal(bytes, &leaderboard)
	if err != nil {
		fmt.Println("Error leaderboard update is not welformed!", string(bytes))
	} else {
		if len(leaderboard.Leaders) > 0 {
			fmt.Println()
			fmt.Println("Leaderboard: ")
			for i, leader := range leaderboard.Leaders {
				fmt.Println(i+1, ". ", leader.UserID, ": ", leader.Score)
			}
		}
	}
}

func NewWorkflow(quizID quiz.QuizIDType, userID quiz.UserIDType, quizEventSubscriber event_sdk.Subscriber) *Workflow {
	return &Workflow{
		quizUser:            quiz.NewQuizUser(quizID, userID),
		quizEventSubscriber: quizEventSubscriber,
		quizClient:          quiz_sdk.NewClient("http://localhost:" + service_shared.CONST_QUIZ_SERVICE_PORT),
	}
}

func (w *Workflow) Play() {
	questions, err := w.joinQuiz()
	if err != nil {
		fmt.Println("Error joining quiz: ", err.Error())
		return
	}
	fmt.Printf("\nWelcome [%s] to the quiz [%s]", w.quizUser.UserID, w.quizUser.QuizID)
	for i, question := range questions {
		fmt.Printf("\n%d. %s:", i, question.Ask)
		for choice, option := range question.Options {
			fmt.Printf("\n%s. %s", choice, option)
		}
		for {
			fmt.Println()
			fmt.Println("Which one is your final choice?: ")
			answer := ""
			fmt.Scanf("%s", &answer)
			if _, valid := question.Options[answer]; valid {
				if err := w.submitAnswer(question, answer); err != nil {
					fmt.Println("submitting answer failed: ", err.Error())
					// must break all here.
				}
				//Wait for leaderboard to be updated to not getting mess
				time.Sleep(20 * time.Millisecond)
				break
			} else if answer == "i" {
				break
			} else {
				fmt.Println("Invalid choice, please try again or ignore by typing `i + enter`:")
			}

		}
	}
}

func (w *Workflow) displayUserStatus() {
	fmt.Println("YOUR SCORE: ", w.userStatus.Score)
	fmt.Println("YOUR RANK: ", w.userStatus.Rank)
}

func (w *Workflow) joinQuiz() ([]quiz_service.Question, error) {
	ctx := context.Background()
	response, err := w.quizClient.JoinQuiz(ctx, &quiz_service.JoinQuizRequest{w.quizUser})
	if err != nil {
		fmt.Println("Error joining quizz: ", err.Error())
		return nil, err
	}
	w.userStatus = response.UserStatus
	w.displayUserStatus()

	err = w.quizEventSubscriber.Connect(ctx)
	if err != nil {
		fmt.Println("Error: Cannot connect to event service", err.Error())
		return nil, err
	}

	w.quizEventSubscriber.Subscribe(ctx,
		event_service.MakeChannelID(w.quizUser.QuizID, "leaderboard"), &leaderboardUpdateHandler{})

	w.quizEventSubscriber.Subscribe(ctx,
		event_service.MakeChannelID(w.quizUser.QuizID, w.quizUser.UserID), &userStatusUpdateHandler{})

	return response.Questions, nil
}

func (w *Workflow) submitAnswer(question quiz_service.Question, userChoice quiz_service.ChoiceType) error {
	response, err := w.quizClient.SubmitAnswer(context.Background(), &quiz_service.SubmitAnswerRequest{
		QuizUser:   w.quizUser,
		QuestionID: question.ID,
		Choice:     userChoice,
	})
	if err != nil {
		fmt.Println("Error summiting answer: ", err.Error())
		return err
	}

	fmt.Println(response.Message)
	if response.Correct {
		fmt.Printf("\nYou've just gained +%d", response.AdditionalScore)
	} else {
		fmt.Print("\nYou've just lost -%d", response.AdditionalScore)
	}
	return nil
}
