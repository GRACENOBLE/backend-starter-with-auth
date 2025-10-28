package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env not loaded: %v", err)
	}

	allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		for _, o := range strings.Split(allowedOriginsEnv, ",") {
			if s := strings.TrimSpace(o); s != "" {
				allowedOrigins = append(allowedOrigins, s)
			}
		}
	}

	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"https://*", "http://*"}
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	r.Get("/auth/{provider}", s.beginAuthHandler)

	r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)

	r.Get("/logout/{provider}", s.logout)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) beginAuthHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothic.BeginAuthHandler(w, r)
}

func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	provider := chi.URLParam(r, "provider")
	redirectURL := os.Getenv("APP_URI")

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("Auth error: %v", err)
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	log.Printf("User authenticated: %s (%s)", user.Name, user.Email)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env File")
	}
	postLogoutRedirectURL := os.Getenv("POST_LOGOUT_REDIRECT_URL")
	gothic.Logout(w, r)
	w.Header().Set("Location", postLogoutRedirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
