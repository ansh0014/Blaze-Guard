package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"auth/config"
	"auth/models"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gorilla/sessions"
	"google.golang.org/api/option"
)

type AuthHandler struct {
	userRepo    *models.UserRepository
	config      *config.Config
	store       *sessions.CookieStore
	firebaseApp *firebase.App
	authClient  *auth.Client
}

func NewAuthHandler(userRepo *models.UserRepository, cfg *config.Config, store *sessions.CookieStore) (*AuthHandler, error) {
	opt := option.WithCredentialsFile(cfg.Firebase.CredentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &AuthHandler{
		userRepo:    userRepo,
		config:      cfg,
		store:       store,
		firebaseApp: app,
		authClient:  authClient,
	}, nil
}

type LoginRequest struct {
	IDToken string `json:"idToken"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.authClient.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	firebaseUID := token.UID
	email := ""
	if e, ok := token.Claims["email"].(string); ok {
		email = e
	}
	name := ""
	if n, ok := token.Claims["name"].(string); ok {
		name = n
	}
	picture := ""
	if p, ok := token.Claims["picture"].(string); ok {
		picture = p
	}

	user, err := h.userRepo.FindByFirebaseUID(firebaseUID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		user = &models.User{
			FirebaseUID: firebaseUID,
			Email:       email,
			Name:        name,
			Picture:     picture,
			Role:        "fire_authority",
		}
		if err := h.userRepo.Create(user); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else {
		user.Email = email
		user.Name = name
		user.Picture = picture
		if err := h.userRepo.Update(user); err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}

	session, _ := h.store.Get(r, h.config.Session.Name)
	session.Values["user_id"] = user.ID
	session.Values["firebase_uid"] = user.FirebaseUID
	session.Values["email"] = user.Email
	session.Values["name"] = user.Name
	session.Values["role"] = user.Role

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, h.config.Session.Name)
	session.Options.MaxAge = -1
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, h.config.Session.Name)

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	firebaseUID, _ := session.Values["firebase_uid"].(string)
	user, err := h.userRepo.FindByFirebaseUID(firebaseUID)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.ID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "fire-authority-auth",
	})
}

func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.store.Get(r, h.config.Session.Name)

		if _, ok := session.Values["user_id"].(int); !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *AuthHandler) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", h.config.Server.FrontendURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
