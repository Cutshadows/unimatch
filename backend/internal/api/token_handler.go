package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"unimatch-back/internal/store"
	"unimatch-back/internal/tokens"
	"unimatch-back/util"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "Invalid request body"})
		return
	}
	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil {
		h.logger.Printf("Error retrieving user: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "Internal server error"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: passwordsDoMatch: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "Internal server error"})
		return
	}
	if !passwordsDoMatch {
		util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "Invalid username or password"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)

	if err != nil {
		h.logger.Printf("ERROR: Creating token: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "Internal server error"})
		return
	}

	util.WriteJSON(w, http.StatusCreated, util.Envelope{
		"token": token,
	})

}
