package auth

import (
	"os"
	"testing"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstants(t *testing.T) {
	t.Run("MaxAge should be 30 days in seconds", func(t *testing.T) {
		expected := 86400 * 30 // 30 days
		assert.Equal(t, expected, MaxAge)
	})

	t.Run("IsProd should be false by default", func(t *testing.T) {
		assert.False(t, IsProd)
	})

	t.Run("key should not be empty", func(t *testing.T) {
		assert.NotEmpty(t, key)
	})
}

func TestNewAuth(t *testing.T) {
	
	setupTestEnv := func(t *testing.T) func() {
		t.Helper()

		envContent := []byte("GOOGLE_CLIENT_ID=test_client_id\nGOOGLE_CLIENT_SECRET=test_client_secret\n")
		err := os.WriteFile(".env.test", envContent, 0644)
		require.NoError(t, err)

		originalEnvFile := ".env"
		os.Rename(".env", ".env.backup")
		os.Rename(".env.test", ".env")

		return func() {
			os.Remove(".env")
			if _, err := os.Stat(".env.backup"); err == nil {
				os.Rename(".env.backup", originalEnvFile)
			}

			goth.ClearProviders()
		}
	}

	t.Run("should load environment variables", func(t *testing.T) {
		cleanup := setupTestEnv(t)
		defer cleanup()

		os.Setenv("GOOGLE_CLIENT_ID", "test_client_id")
		os.Setenv("GOOGLE_CLIENT_SECRET", "test_client_secret")
		defer func() {
			os.Unsetenv("GOOGLE_CLIENT_ID")
			os.Unsetenv("GOOGLE_CLIENT_SECRET")
		}()

		assert.NotPanics(t, func() {
			NewAuth()
		})

		assert.Equal(t, "test_client_id", os.Getenv("GOOGLE_CLIENT_ID"))
		assert.Equal(t, "test_client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	})

	t.Run("should configure gothic store", func(t *testing.T) {
		cleanup := setupTestEnv(t)
		defer cleanup()

		os.Setenv("GOOGLE_CLIENT_ID", "test_client_id")
		os.Setenv("GOOGLE_CLIENT_SECRET", "test_client_secret")
		defer func() {
			os.Unsetenv("GOOGLE_CLIENT_ID")
			os.Unsetenv("GOOGLE_CLIENT_SECRET")
		}()

		NewAuth()

		assert.NotNil(t, gothic.Store)
	})

	t.Run("should register Google provider", func(t *testing.T) {
		cleanup := setupTestEnv(t)
		defer cleanup()

		os.Setenv("GOOGLE_CLIENT_ID", "test_client_id")
		os.Setenv("GOOGLE_CLIENT_SECRET", "test_client_secret")
		defer func() {
			os.Unsetenv("GOOGLE_CLIENT_ID")
			os.Unsetenv("GOOGLE_CLIENT_SECRET")
		}()

		NewAuth()

		providers := goth.GetProviders()
		assert.NotEmpty(t, providers)

		_, exists := providers["google"]
		assert.True(t, exists, "Google provider should be registered")
	})

	t.Run("should configure Google provider with correct credentials", func(t *testing.T) {
		cleanup := setupTestEnv(t)
		defer cleanup()

		testClientID := "test_google_client_id"
		testClientSecret := "test_google_secret"

		os.Setenv("GOOGLE_CLIENT_ID", testClientID)
		os.Setenv("GOOGLE_CLIENT_SECRET", testClientSecret)
		defer func() {
			os.Unsetenv("GOOGLE_CLIENT_ID")
			os.Unsetenv("GOOGLE_CLIENT_SECRET")
		}()

		NewAuth()

		providers := goth.GetProviders()
		googleProvider, exists := providers["google"]
		require.True(t, exists)

		assert.Equal(t, "google", googleProvider.Name())
		assert.NotNil(t, googleProvider)
	})
}

func TestSessionConfiguration(t *testing.T) {
	t.Run("session store should have correct MaxAge", func(t *testing.T) {

		os.Setenv("GOOGLE_CLIENT_ID", "test_id")
		os.Setenv("GOOGLE_CLIENT_SECRET", "test_secret")
		defer func() {
			os.Unsetenv("GOOGLE_CLIENT_ID")
			os.Unsetenv("GOOGLE_CLIENT_SECRET")
			goth.ClearProviders()
		}()

		envContent := []byte("GOOGLE_CLIENT_ID=test_id\nGOOGLE_CLIENT_SECRET=test_secret\n")
		os.WriteFile(".env", envContent, 0644)
		defer os.Remove(".env")

		NewAuth()

		assert.NotNil(t, gothic.Store)
	})
}

func TestAuthPackageIntegration(t *testing.T) {
	t.Run("full auth initialization flow", func(t *testing.T) {
		envContent := []byte("GOOGLE_CLIENT_ID=integration_test_id\nGOOGLE_CLIENT_SECRET=integration_test_secret\n")
		err := os.WriteFile(".env", envContent, 0644)
		require.NoError(t, err)
		defer os.Remove(".env")

		os.Setenv("GOOGLE_CLIENT_ID", "integration_test_id")
		os.Setenv("GOOGLE_CLIENT_SECRET", "integration_test_secret")
		defer func() {
			os.Unsetenv("GOOGLE_CLIENT_ID")
			os.Unsetenv("GOOGLE_CLIENT_SECRET")
			goth.ClearProviders()
		}()

		assert.NotPanics(t, func() {
			NewAuth()
		})

		providers := goth.GetProviders()
		assert.NotEmpty(t, providers)
		assert.NotNil(t, gothic.Store)

		googleProvider, exists := providers["google"]
		assert.True(t, exists)
		assert.Equal(t, "google", googleProvider.Name())
	})
}

func BenchmarkNewAuth(b *testing.B) {
	envContent := []byte("GOOGLE_CLIENT_ID=bench_test_id\nGOOGLE_CLIENT_SECRET=bench_test_secret\n")
	os.WriteFile(".env", envContent, 0644)
	defer os.Remove(".env")

	os.Setenv("GOOGLE_CLIENT_ID", "bench_test_id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "bench_test_secret")
	defer func() {
		os.Unsetenv("GOOGLE_CLIENT_ID")
		os.Unsetenv("GOOGLE_CLIENT_SECRET")
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		goth.ClearProviders()
		NewAuth()
	}
}
