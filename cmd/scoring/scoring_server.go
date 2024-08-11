package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	scoringconf "github.com/sepcon/quizprob/config/scoring"
	"github.com/sepcon/quizprob/internal/common_errors"
	"github.com/sepcon/quizprob/internal/concurrent"
	"github.com/sepcon/quizprob/internal/scoring"
	"github.com/sepcon/quizprob/internal/scoring/ranker"
	"github.com/sepcon/quizprob/internal/scoring/user_score_db"
	"github.com/sepcon/quizprob/pkg/event_sdk"
	"github.com/sepcon/quizprob/pkg/model/event_service"
	"github.com/sepcon/quizprob/pkg/model/quiz"
	"github.com/sepcon/quizprob/pkg/model/scoring_service"
	"net/http"
)

type ScoringServer struct {
	scoringManager scoring.ScoringManager
	workerPool     concurrent.WorkerPool
	router         *gin.Engine
}

func NewScoringServer(conf scoringconf.Config) (ss *ScoringServer, err error) {
	var userScoreDB user_score_db.UserScoreDB
	if conf.Persistence != nil {
		userScoreDB, err = user_score_db.NewPostGreSQL(conf.Persistence)
		if err != nil {
			return
		}
	}

	ranker, err := ranker.NewRedisRanker(conf)
	if err != nil {
		return
	}

	publisher := event_sdk.NewPublisher(conf.EventServiceUrl)
	return &ScoringServer{
		scoringManager: scoring.NewScoringManager(
			ranker,
			scoring.OptPersistentUserScoreDB(userScoreDB),
			scoring.OptUserScoreUpdateCallback(func(quizID quiz.QuizIDType, status *quiz.UserStatus) {
				event, _ := json.Marshal(status)
				fmt.Printf("\nUser status updated: %s", event)
				publisher.Publish(context.Background(), event_service.MakeChannelID(quizID, status.UserID), event)
			}),
			scoring.OptLeaderboardUpdateCallback(func(leaderboard quiz.Leaderboard) {
				event, _ := json.Marshal(leaderboard)
				fmt.Printf("\nLeaderboard update: %s", event)
				publisher.Publish(context.Background(), event_service.MakeChannelID(leaderboard.QuizID, "leaderboard"),
					event)
			}),
		),
		workerPool: concurrent.NewWorkerPool(conf.Concurrent.NumberOfWorker, conf.Concurrent.MaxQueueLength),
		router:     gin.New(),
	}, nil
}

func (ss *ScoringServer) Serve(addr string) {
	ss.registerHandlers()
	ss.workerPool.Start()
	ss.router.Run(addr)
}

func (ss *ScoringServer) registerHandlers() {
	ss.router.POST("/create_quiz", ss.onCreateQuizRequest)
	ss.router.POST("/update_score", ss.onUserScoreUpdate)
	ss.router.GET("/leaderboard", ss.getLeaderBoard)
	ss.router.POST("/user_status", ss.getOrInitUserStatus)
}

func abortedByInvalidRequestFormat(ctx *gin.Context, req any) bool {
	if ctx.BindJSON(req) != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON format"})
		return true
	}
	return false
}

func makeErrorResponseBody(msg string) gin.H {
	return gin.H{"message": msg}
}

func (ss *ScoringServer) onCreateQuizRequest(ctx *gin.Context) {
	var req scoring_service.CreateQuizRequest
	if abortedByInvalidRequestFormat(ctx, &req) {
		return
	} else if err := ss.scoringManager.CreateQuiz(req); err != nil {
		if common_errors.Match(err, common_errors.EC_QUIZ_ALREADY_EXIST) {
			ctx.AbortWithStatusJSON(http.StatusConflict, makeErrorResponseBody(err.Error()))
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	} else {
		ctx.Status(http.StatusOK)
	}
}

func (ss *ScoringServer) onUserScoreUpdate(ctx *gin.Context) {
	var req scoring_service.UpdateScoreRequest
	if abortedByInvalidRequestFormat(ctx, &req) {
		return
	} else if err := ss.workerPool.Submit(func() {
		ss.scoringManager.UpdateUserScore(req)
	}); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusAccepted, &scoring_service.UpdateScoreResponse{
			Accepted: true,
		})
	}
}

func (ss *ScoringServer) getLeaderBoard(ctx *gin.Context) {
	var leaderboardRequest scoring_service.LeaderBoardRequest
	leaderboardRequest.QuizID = ctx.Query("quizid")
	if leaderboardRequest.QuizID == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
	} else if leaderboard, err := ss.scoringManager.GetLeaderboard(leaderboardRequest); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else if leaderboard == nil {
		ctx.AbortWithError(http.StatusNotFound,
			fmt.Errorf("not found leaderboard for quizid [%s]", leaderboardRequest.QuizID))
	} else {
		ctx.JSON(http.StatusOK, leaderboard)
	}
}

func (ss *ScoringServer) getOrInitUserStatus(ctx *gin.Context) {
	var req scoring_service.UserStatusRequest
	if ctx.BindJSON(&req) != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid JSON format"))
	} else if status, err := ss.scoringManager.GetOrInitUserStatus(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, status)
	}
}
