package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sepcon/quizprob/pkg/model/quiz"
	"github.com/sepcon/quizprob/pkg/model/quiz_service"
	"github.com/sepcon/quizprob/pkg/model/scoring_service"
	"github.com/sepcon/quizprob/pkg/scoring_sdk"
	"net/http"
	"strconv"
	"strings"
)

type QuizServer struct {
	quiz            quiz.Quiz
	leaderboardSize uint32
	questions       Questions
	correctAnswers  CorrectAnswers
	scoringClient   scoring_sdk.Client
	router          *gin.Engine
	addr            string
}

func NewQuizServer(addr string, quizID quiz.QuizIDType, leaderboardSize uint32, client scoring_sdk.Client) *QuizServer {
	q := &QuizServer{
		quiz:            quiz.Quiz{QuizID: quizID},
		leaderboardSize: leaderboardSize,
		scoringClient:   client,
		router:          gin.New(),
		addr:            addr,
	}
	q.questions, q.correctAnswers = MakeQuestionBank()
	return q
}

func (q *QuizServer) Start() error {
	if err := q.scoringClient.CreateQuiz(context.Background(), &scoring_service.CreateQuizRequest{
		Quiz:            q.quiz,
		LeaderboardSize: q.leaderboardSize,
	}); err != nil {
		fmt.Println("Error creating quiz: ", err.Error())
		if !strings.Contains(err.Error(), strconv.Itoa(http.StatusConflict)) {
			return err
		}
	}
	q.registerHandlers()
	return q.router.Run(q.addr)
}

func (q *QuizServer) registerHandlers() {
	q.router.POST("/join_quiz", q.onUserJoinQuiz)
	q.router.POST("/submit_answer", q.onUserAnswerSubmitted)
}

func (q *QuizServer) onUserJoinQuiz(ctx *gin.Context) {
	var req quiz_service.JoinQuizRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid Json format"})
		return
	}
	userStatus, err := q.scoringClient.GetOrInitUserStatus(ctx.Request.Context(), &scoring_service.UserStatusRequest{
		req.QuizUser,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Cannot initialize or get the latest user stauts"})
		return
	}
	if userStatus == nil {
		fmt.Println("Error, null user status: ", req.UserID)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	leaderboardResponse, err := q.scoringClient.GetLeaderboard(ctx.Request.Context(), req.QuizID)
	if err != nil {
		fmt.Println("Error getting leaderboard: ", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, &quiz_service.JoinQuizResponse{
		UserStatus:  *userStatus,
		Leaderboard: leaderboardResponse.Leaders,
		Questions:   q.questions,
	})
}

func (q *QuizServer) onUserAnswerSubmitted(ctx *gin.Context) {
	var req quiz_service.SubmitAnswerRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid Json format"})
		return
	}
	correctChoice, validQuestion := q.correctAnswers[req.QuestionID]
	if !validQuestion {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unkown question id: %v", req.QuestionID)})
		return
	}
	correct := req.Choice == correctChoice.Choice
	message := ""
	additionalScore := quiz.ScoreType(0)
	if correct {
		message = "congrat!"
		additionalScore = correctChoice.Score
	} else {
		message = "sorry, its not a correct answer"
		additionalScore = -correctChoice.Score
	}
	err = q.scoringClient.UpdateUserScore(ctx.Request.Context(), &scoring_service.UpdateScoreRequest{
		QuizUser:   req.QuizUser,
		ScoreDelta: additionalScore,
	})

	if err != nil {
		fmt.Println("Error updating user score: ", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := &quiz_service.SubmitAnswerResponse{
		Correct:         correct,
		Message:         message,
		AdditionalScore: additionalScore,
	}

	ctx.JSON(http.StatusOK, response)
}
