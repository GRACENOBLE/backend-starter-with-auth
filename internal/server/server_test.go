package server

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockDatabaseService struct{}

func (m *MockDatabaseService) Health() map[string]string {
	return map[string]string{
		"status":  "up",
		"message": "It's healthy",
	}
}

func (m *MockDatabaseService) Close() error {
	return nil
}

func TestNewServer(t *testing.T) {
	t.Run("should create server with correct configuration", func(t *testing.T) {
	
		os.Setenv("PORT", "3000")
		defer os.Unsetenv("PORT")

		server := NewServer()

		require.NotNil(t, server)
		assert.Equal(t, ":3000", server.Addr)
		assert.NotNil(t, server.Handler)
		assert.Equal(t, time.Minute, server.IdleTimeout)
		assert.Equal(t, 10*time.Second, server.ReadTimeout)
		assert.Equal(t, 30*time.Second, server.WriteTimeout)
	})

	t.Run("should handle missing PORT environment variable", func(t *testing.T) {
		os.Unsetenv("PORT")

		server := NewServer()

		require.NotNil(t, server)
		assert.Equal(t, ":0", server.Addr)
	})

	t.Run("should handle invalid PORT environment variable", func(t *testing.T) {
		os.Setenv("PORT", "invalid")
		defer os.Unsetenv("PORT")

		server := NewServer()

		require.NotNil(t, server)
		assert.Equal(t, ":0", server.Addr)
	})

	t.Run("should use different port values", func(t *testing.T) {
		testCases := []struct {
			name     string
			port     string
			expected string
		}{
			{"port 8080", "8080", ":8080"},
			{"port 5000", "5000", ":5000"},
			{"port 80", "80", ":80"},
			{"port 443", "443", ":443"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				os.Setenv("PORT", tc.port)
				defer os.Unsetenv("PORT")

				server := NewServer()

				assert.Equal(t, tc.expected, server.Addr)
			})
		}
	})
}

func TestServerStruct(t *testing.T) {
	t.Run("should create Server with port and database", func(t *testing.T) {
		mockDB := &MockDatabaseService{}
		server := &Server{
			port: 3000,
			db:   mockDB,
		}

		assert.Equal(t, 3000, server.port)
		assert.NotNil(t, server.db)
	})

	t.Run("database health should be accessible through Server", func(t *testing.T) {
		mockDB := &MockDatabaseService{}
		server := &Server{
			// port: 3000,
			db:   mockDB,
		}

		health := server.db.Health()

		assert.Equal(t, "up", health["status"])
		assert.Equal(t, "It's healthy", health["message"])
	})
}

func TestServerTimeouts(t *testing.T) {
	t.Run("should have correct timeout configurations", func(t *testing.T) {
		os.Setenv("PORT", "3000")
		defer os.Unsetenv("PORT")

		server := NewServer()

		assert.Equal(t, time.Minute, server.IdleTimeout, "IdleTimeout should be 1 minute")
		assert.Equal(t, 10*time.Second, server.ReadTimeout, "ReadTimeout should be 10 seconds")
		assert.Equal(t, 30*time.Second, server.WriteTimeout, "WriteTimeout should be 30 seconds")
	})
}

func TestServerHandler(t *testing.T) {
	t.Run("should have a valid handler", func(t *testing.T) {
		os.Setenv("PORT", "3000")
		defer os.Unsetenv("PORT")

		server := NewServer()

		require.NotNil(t, server.Handler)
		assert.Implements(t, (*http.Handler)(nil), server.Handler)
	})
}

func TestServerIntegration(t *testing.T) {
	t.Run("should create a fully functional HTTP server", func(t *testing.T) {
		os.Setenv("PORT", "8888")
		defer os.Unsetenv("PORT")

		server := NewServer()

		assert.NotNil(t, server)
		assert.Equal(t, ":8888", server.Addr)
		assert.NotNil(t, server.Handler)

		assert.Greater(t, server.IdleTimeout, time.Duration(0))
		assert.Greater(t, server.ReadTimeout, time.Duration(0))
		assert.Greater(t, server.WriteTimeout, time.Duration(0))
	})
}

// Benchmark tests
func BenchmarkNewServer(b *testing.B) {
	os.Setenv("PORT", "3000")
	defer os.Unsetenv("PORT")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewServer()
	}
}

func BenchmarkServerCreationWithDifferentPorts(b *testing.B) {
	ports := []string{"3000", "8080", "5000", "9000"}

	for _, port := range ports {
		b.Run("port_"+port, func(b *testing.B) {
			os.Setenv("PORT", port)
			defer os.Unsetenv("PORT")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = NewServer()
			}
		})
	}
}
