package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, username, password string) error {
	args := m.Called(ctx, username, password)
	return args.Error(0)
}

func (m *MockUserService) AddContact(ctx context.Context, userID int64, contactUsername string) error {
	args := m.Called(ctx, userID, contactUsername)
	return args.Error(0)
}

func (m *MockUserService) GetUserContacts(ctx context.Context, userID int64) ([]domain.Contact, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Contact), args.Error(1)
}

func (m *MockUserService) GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]domain.User), args.Error(1)
}

func TestCreateUserHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.CreateUserHandler(mockService)

		reqBody := `{"username":"testuser", "password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
		rr := httptest.NewRecorder()

		mockService.On("CreateUser", mock.Anything, "testuser", "password123").Return(nil)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.CreateUserHandler(mockService)

		reqBody := `{"username":"testuser", "password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
		rr := httptest.NewRecorder()

		mockService.On("CreateUser", mock.Anything, "testuser", "password123").Return(errors.New("service error"))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetUsersHandler(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.GetUsersHandler(mockService)
		req := httptest.NewRequest(http.MethodGet, "/users?username=testuser", nil)
		rr := httptest.NewRecorder()
		expectedUsers := []domain.User{{UserId: 1, Username: "testuser"}}

		mockService.On("GetUsersByUsername", mock.Anything, "testuser").Return(expectedUsers, nil)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var users []domain.User
		err := json.NewDecoder(rr.Body).Decode(&users)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		mockService.AssertExpectations(t)
	})

	t.Run("missing username query parameter", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.GetUsersHandler(mockService)
		req := httptest.NewRequest(http.MethodGet, "/get_users", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.GetUsersHandler(mockService)
		req := httptest.NewRequest(http.MethodGet, "/get_users?username=testuser", nil)
		rr := httptest.NewRecorder()
		var users []domain.User

		mockService.On("GetUsersByUsername", mock.Anything, "testuser").Return(users, errors.New("service error"))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCreateContactHandler(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.CreateContactHandler(mockService)
		reqBody := `{"contactUsername":"contactuser"}`
		req := httptest.NewRequest(http.MethodPost, "/create_contact", bytes.NewBufferString(reqBody))
		req.Header.Set("User-Id", "1")
		rr := httptest.NewRecorder()

		mockService.On("AddContact", mock.Anything, int64(1), "contactuser").Return(nil)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.CreateContactHandler(mockService)
		reqBody := `{"contactUsername":}`
		req := httptest.NewRequest(http.MethodPost, "/create_contact", bytes.NewBufferString(reqBody))
		req.Header.Set("User-Id", "1")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.CreateContactHandler(mockService)
		reqBody := `{"contactUsername":"contactuser"}`
		req := httptest.NewRequest(http.MethodPost, "/create_contact", bytes.NewBufferString(reqBody))
		req.Header.Set("User-Id", "invalid")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.CreateContactHandler(mockService)
		reqBody := `{"contactUsername":"contactuser"}`
		req := httptest.NewRequest(http.MethodPost, "/create_contact", bytes.NewBufferString(reqBody))
		req.Header.Set("User-Id", "1")
		rr := httptest.NewRecorder()

		mockService.On("AddContact", mock.Anything, int64(1), "contactuser").Return(errors.New("service error"))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetContactsHandler(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.GetContactsHandler(mockService)
		req := httptest.NewRequest(http.MethodGet, "/get_contacts", nil)
		req.Header.Set("User-Id", "1")
		rr := httptest.NewRecorder()
		expectedContacts := []domain.Contact{{UserId: 1, ContactUserId: 2}}

		mockService.On("GetUserContacts", mock.Anything, int64(1)).Return(expectedContacts, nil)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var contacts []domain.Contact
		err := json.NewDecoder(rr.Body).Decode(&contacts)
		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.GetContactsHandler(mockService)
		req := httptest.NewRequest(http.MethodGet, "/get_contacts", nil)
		req.Header.Set("User-Id", "invalid")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := handlers.GetContactsHandler(mockService)
		req := httptest.NewRequest(http.MethodGet, "/get_contacts", nil)
		req.Header.Set("User-Id", "1")
		rr := httptest.NewRecorder()
		var contacts []domain.Contact
		mockService.On("GetUserContacts", mock.Anything, int64(1)).Return(contacts, errors.New("service error"))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})
}
