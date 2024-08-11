package scoring_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sepcon/quizprob/pkg/model/quiz"
	"github.com/sepcon/quizprob/pkg/model/scoring_service"
	"io"
	"net/http"
)

type scoringClientImpl struct {
	baseURL string
	client  *http.Client
}

func NewScoringClient(baseURL string) *scoringClientImpl {
	return &scoringClientImpl{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (sc *scoringClientImpl) CreateQuiz(ctx context.Context, req *scoring_service.CreateQuizRequest) error {
	url := fmt.Sprintf("%s/create_quiz", sc.baseURL)
	return sc.sendRequest(ctx, http.MethodPost, url, req, nil)
}

func (sc *scoringClientImpl) UpdateUserScore(ctx context.Context, req *scoring_service.UpdateScoreRequest) error {
	url := fmt.Sprintf("%s/update_score", sc.baseURL)
	return sc.sendRequest(ctx, http.MethodPost, url, req, nil)
}

func (sc *scoringClientImpl) GetLeaderboard(ctx context.Context, quizID quiz.QuizIDType) (*quiz.Leaderboard, error) {
	url := fmt.Sprintf("%s/leaderboard?quizid=%s", sc.baseURL, quizID)
	var leaderboard quiz.Leaderboard
	err := sc.sendRequest(ctx, http.MethodGet, url, nil, &leaderboard)
	if err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

func (sc *scoringClientImpl) GetOrInitUserStatus(ctx context.Context, req *scoring_service.UserStatusRequest) (*quiz.UserStatus, error) {
	url := fmt.Sprintf("%s/user_status", sc.baseURL)
	var status quiz.UserStatus
	err := sc.sendRequest(ctx, http.MethodPost, url, req, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (sc *scoringClientImpl) sendRequest(ctx context.Context, method, url string, requestBody, responseBody interface{}) error {
	var reqBody io.Reader
	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := sc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	if responseBody != nil {
		err = json.NewDecoder(resp.Body).Decode(responseBody)
		if err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}
