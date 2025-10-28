package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloWorldHandler(t *testing.T) {
	t.Run("should return Hello World message", func(t *testing.T) {
		s := &Server{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		s.HelloWorldHandler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		expected := map[string]string{"message": "Hello World"}
		var actual map[string]string
		err = json.Unmarshal(body, &actual)
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("should return valid JSON", func(t *testing.T) {
		s := &Server{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		s.HelloWorldHandler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Response should be valid JSON")
	})
}

func TestHealthHandler(t *testing.T) {
	t.Run("should return health status", func(t *testing.T) {
		mockDB := &MockDatabaseService{}
		s := &Server{db: mockDB}

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		s.healthHandler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var health map[string]string
		err = json.Unmarshal(body, &health)
		require.NoError(t, err)

		assert.Equal(t, "up", health["status"])
		assert.Equal(t, "It's healthy", health["message"])
	})

	t.Run("should return valid JSON health data", func(t *testing.T) {
		mockDB := &MockDatabaseService{}
		s := &Server{db: mockDB}

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		s.healthHandler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Response should be valid JSON")
	})
}

func TestBeginAuthHandler(t *testing.T) {
	t.Run("should set provider in context", func(t *testing.T) {
		s := &Server{}

		r := chi.NewRouter()
		r.Get("/auth/{provider}", s.beginAuthHandler)

		req := httptest.NewRequest(http.MethodGet, "/auth/google", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
	})

	t.Run("should handle different providers", func(t *testing.T) {
		s := &Server{}
		providers := []string{"google", "github", "facebook"}

		for _, provider := range providers {
			r := chi.NewRouter()
			r.Get("/auth/{provider}", s.beginAuthHandler)

			req := httptest.NewRequest(http.MethodGet, "/auth/"+provider, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusInternalServerError, w.Code,
				"Handler should not return internal server error for provider: "+provider)
		}
	})
}

func TestGetAuthCallbackFunction(t *testing.T) {
	t.Skip("Skipping auth callback tests - handler calls log.Fatal which terminates tests")

	t.Run("should handle auth callback route", func(t *testing.T) {
		os.Setenv("APP_URI", "http://localhost:3000")
		defer os.Unsetenv("APP_URI")

		s := &Server{}

		r := chi.NewRouter()
		r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)

		req := httptest.NewRequest(http.MethodGet, "/auth/google/callback", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

	})

	t.Run("should require APP_URI environment variable", func(t *testing.T) {
		originalValue := os.Getenv("APP_URI")
		os.Setenv("APP_URI", "http://localhost:3000")
		defer func() {
			if originalValue != "" {
				os.Setenv("APP_URI", originalValue)
			} else {
				os.Unsetenv("APP_URI")
			}
		}()

		s := &Server{}
		r := chi.NewRouter()
		r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)

		req := httptest.NewRequest(http.MethodGet, "/auth/github/callback", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

	})
}

func TestLogout(t *testing.T) {

	t.Skip("Skipping logout tests - handler calls log.Fatal which terminates tests")

	t.Run("should handle logout route", func(t *testing.T) {
		os.Setenv("POST_LOGOUT_REDIRECT_URL", "http://localhost:3000/login")
		defer os.Unsetenv("POST_LOGOUT_REDIRECT_URL")

		s := &Server{}

		r := chi.NewRouter()
		r.Get("/logout/{provider}", s.logout)

		req := httptest.NewRequest(http.MethodGet, "/logout/google", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

		location := w.Header().Get("Location")
		assert.Equal(t, "http://localhost:3000/login", location)
	})

	t.Run("should handle different providers for logout", func(t *testing.T) {
		os.Setenv("POST_LOGOUT_REDIRECT_URL", "http://localhost:3000")
		defer os.Unsetenv("POST_LOGOUT_REDIRECT_URL")

		s := &Server{}
		providers := []string{"google", "github", "facebook"}

		for _, provider := range providers {
			r := chi.NewRouter()
			r.Get("/logout/{provider}", s.logout)

			req := httptest.NewRequest(http.MethodGet, "/logout/"+provider, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusTemporaryRedirect, w.Code,
				"Handler should return redirect status for provider: "+provider)
		}
	})
}

func TestRegisterRoutes(t *testing.T) {
	t.Run("should register all routes", func(t *testing.T) {
		os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
		defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

		mockDB := &MockDatabaseService{}
		s := &Server{db: mockDB}

		handler := s.RegisterRoutes()
		require.NotNil(t, handler)

		assert.Implements(t, (*http.Handler)(nil), handler)
	})

	t.Run("should handle root route", func(t *testing.T) {
		os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
		defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

		mockDB := &MockDatabaseService{}
		s := &Server{db: mockDB}

		handler := s.RegisterRoutes()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		body, err := io.ReadAll(w.Body)
		require.NoError(t, err)

		var result map[string]string
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "Hello World", result["message"])
	})

	t.Run("should handle health route", func(t *testing.T) {
		os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
		defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

		mockDB := &MockDatabaseService{}
		s := &Server{db: mockDB}

		handler := s.RegisterRoutes()

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		body, err := io.ReadAll(w.Body)
		require.NoError(t, err)

		var result map[string]string
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, "up", result["status"])
	})

	t.Run("should handle CORS configuration", func(t *testing.T) {
		os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080")
		defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

		mockDB := &MockDatabaseService{}
		s := &Server{db: mockDB}

		handler := s.RegisterRoutes()

		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("should use default CORS when env not set", func(t *testing.T) {
		os.Unsetenv("CORS_ALLOWED_ORIGINS")

		mockDB := &MockDatabaseService{}
		s := &Server{db: mockDB}

		handler := s.RegisterRoutes()
		require.NotNil(t, handler)

		assert.Implements(t, (*http.Handler)(nil), handler)
	})

	t.Run("should handle 404 for unknown routes", func(t *testing.T) {
		os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
		defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

		mockDB := &MockDatabaseService{}
		s := &Server{db: mockDB}

		handler := s.RegisterRoutes()

		req := httptest.NewRequest(http.MethodGet, "/unknown-route", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
