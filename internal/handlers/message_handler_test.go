package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the MessageService interface
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) CreateMessage(ctx context.Context, senderID int64, recipientID int64, content string) (int64, error) {
	args := m.Called(ctx, senderID, recipientID, content)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMessageService) GetMessagesByUserId(ctx context.Context, userID int64) ([]domain.Message, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func TestCreateMessageHandler_Success(t *testing.T) {
	mockService := new(MockMessageService)

	messageID := int64(123)
	mockService.On("CreateMessage", mock.Anything, int64(1), int64(2), "Hello").Return(messageID, nil)

	reqBody, _ := json.Marshal(handlers.CreateMessageRequest{Content: "Hello"})
	req := httptest.NewRequest(http.MethodPost, "/messages/user/2", bytes.NewReader(reqBody))
	req.Header.Set("User-Id", "1")

	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /messages/user/{user_id}", handlers.CreateMessageHandler(mockService))
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp handlers.CreateMessageResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	assert.Equal(t, strconv.Itoa(int(messageID)), resp.MessageId)
	mockService.AssertExpectations(t)
}

func TestCreateMessageHandler_InvalidRequest(t *testing.T) {
	mockService := new(MockMessageService)
	handler := handlers.CreateMessageHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/messages/user/2", bytes.NewReader([]byte("invalid")))
	req.Header.Set("User-Id", "1")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateMessageHandler_MissingUserId(t *testing.T) {
	mockService := new(MockMessageService)
	handler := handlers.CreateMessageHandler(mockService)

	reqBody, _ := json.Marshal(handlers.CreateMessageRequest{Content: "Hello"})
	req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(reqBody))
	req.Header.Set("User-Id", "1")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateMessageHandler_ServiceError(t *testing.T) {
	mockService := new(MockMessageService)

	mockService.On("CreateMessage", mock.Anything, int64(1), int64(2), "Hello").Return(int64(0), errors.New("service error"))

	reqBody, _ := json.Marshal(handlers.CreateMessageRequest{Content: "Hello"})
	req := httptest.NewRequest(http.MethodPost, "/messages/user/2", bytes.NewReader(reqBody))
	req.Header.Set("User-Id", "1")

	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /messages/user/{user_id}", handlers.CreateMessageHandler(mockService))
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockService.AssertExpectations(t)
}
