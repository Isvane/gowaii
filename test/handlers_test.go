package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Isvane/gomen/internal/api"
	models "github.com/Isvane/gomen/internal/model"
	"github.com/Isvane/gomen/internal/repository"
)

func TestEchoHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /echo/{arg}", api.EchoHandler)

	req := httptest.NewRequest(http.MethodGet, "/echo/world", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedBody := `Hello, "world"`
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestRegisterUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		setup          func(repo *repository.UserRepo)
		expectedStatus int
	}{
		{
			name:           "successful registration",
			payload:        `{"name": "Alice", "age": 25}`,
			setup:          func(repo *repository.UserRepo) {},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json syntax",
			payload:        `{"name": "Alice", "age":}`,
			setup:          func(repo *repository.UserRepo) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "validation failure: empty name",
			payload:        `{"name": "", "age": 25}`,
			setup:          func(repo *repository.UserRepo) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "validation failure: negative age",
			payload:        `{"name": "Bob", "age": -5}`,
			setup:          func(repo *repository.UserRepo) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "user already exists",
			payload: `{"name": "Charlie", "age": 30}`,
			setup: func(repo *repository.UserRepo) {
				_ = repo.Create(models.User{Name: "Charlie", Age: 30})
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewUserRepo()
			tt.setup(repo)
			db := &api.Database{Repo: repo}

			mux := http.NewServeMux()
			mux.HandleFunc("POST /user", db.RegisterUserHandler)

			req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		targetUrl      string
		setup          func(repo *repository.UserRepo)
		expectedStatus int
		expectedUser   models.User
	}{
		{
			name:      "user found",
			targetUrl: "/user/Alice",
			setup: func(repo *repository.UserRepo) {
				_ = repo.Create(models.User{Name: "Alice", Age: 28})
			},
			expectedStatus: http.StatusOK,
			expectedUser:   models.User{Name: "Alice", Age: 28},
		},
		{
			name:           "user not found",
			targetUrl:      "/user/John",
			setup:          func(repo *repository.UserRepo) {},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewUserRepo()
			tt.setup(repo)
			db := &api.Database{Repo: repo}

			mux := http.NewServeMux()
			mux.HandleFunc("GET /user/{name}", db.GetUserHandler)

			req := httptest.NewRequest(http.MethodGet, tt.targetUrl, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var gotUser models.User
				if err := json.NewDecoder(rec.Body).Decode(&gotUser); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if gotUser != tt.expectedUser {
					t.Errorf("expected user %+v, got %+v", tt.expectedUser, gotUser)
				}
			}
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		targetUrl      string
		payload        string
		setup          func(repo *repository.UserRepo)
		expectedStatus int
		expectedUser   models.User
	}{
		{
			name:      "successful update",
			targetUrl: "/user/Alice",
			payload:   `{"name": "Alicia", "age": 29}`,
			setup: func(repo *repository.UserRepo) {
				_ = repo.Create(models.User{Name: "Alice", Age: 28})
			},
			expectedStatus: http.StatusOK,
			expectedUser:   models.User{Name: "Alicia", Age: 29},
		},
		{
			name:           "user not found",
			targetUrl:      "/user/Ghost",
			payload:        `{"name": "Ghost", "age": 30}`,
			setup:          func(repo *repository.UserRepo) {},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "validation failure: negative age",
			targetUrl: "/user/Alice",
			payload:   `{"name": "Alice", "age": -5}`,
			setup: func(repo *repository.UserRepo) {
				_ = repo.Create(models.User{Name: "Alice", Age: 28})
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewUserRepo()
			tt.setup(repo)
			db := &api.Database{Repo: repo}

			mux := http.NewServeMux()
			mux.HandleFunc("PUT /user/{name}", db.UpdateUserHandler)

			req := httptest.NewRequest(http.MethodPut, tt.targetUrl, bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var gotUser models.User
				if err := json.NewDecoder(rec.Body).Decode(&gotUser); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if gotUser != tt.expectedUser {
					t.Errorf("expected user %+v, got %+v", tt.expectedUser, gotUser)
				}
			}
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		targetUrl      string
		setup          func(repo *repository.UserRepo)
		expectedStatus int
	}{
		{
			name:      "successful delete",
			targetUrl: "/user/Alice",
			setup: func(repo *repository.UserRepo) {
				_ = repo.Create(models.User{Name: "Alice", Age: 28})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user not found",
			targetUrl:      "/user/John",
			setup:          func(repo *repository.UserRepo) {},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewUserRepo()
			tt.setup(repo)
			db := &api.Database{Repo: repo}

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /user/{name}", db.DeleteUserHandler)

			req := httptest.NewRequest(http.MethodDelete, tt.targetUrl, nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				_, err := repo.Get("Alice")
				if err != repository.ErrUserNotFound {
					t.Errorf("expected user to be deleted from repository, but found: %v", err)
				}
			}
		})
	}
}
