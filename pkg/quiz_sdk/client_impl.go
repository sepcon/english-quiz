package quiz_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sepcon/quizprob/pkg/model/quiz_service"
	"net/http"
)

type clientImpl struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) *clientImpl {
	return &clientImpl{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (qc *clientImpl) JoinQuiz(ctx context.Context, req *quiz_service.JoinQuizRequest) (*quiz_service.JoinQuizResponse, error) {
	url := fmt.Sprintf("%s/join_quiz", qc.baseURL)
	var resp quiz_service.JoinQuizResponse
	err := qc.sendRequest(ctx, http.MethodPost, url, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (qc *clientImpl) SubmitAnswer(ctx context.Context, req *quiz_service.SubmitAnswerRequest) (*quiz_service.SubmitAnswerResponse, error) {
	url := fmt.Sprintf("%s/submit_answer", qc.baseURL)
	var resp quiz_service.SubmitAnswerResponse
	err := qc.sendRequest(ctx, http.MethodPost, url, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (qc *clientImpl) sendRequest(ctx context.Context, method, url string, requestBody, responseBody interface{}) error {
	var reqBody *bytes.Buffer
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

	resp, err := qc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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
