package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"unimatch-back/internal/store"
	"unimatch-back/util"
)

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateResgisterUserRequest(req *RegisterUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) > 50 {
		return errors.New("username must be at most 50 characters long")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error: decoding register request: %v", err)
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": "invalid request body"})
		return
	}

	err = h.validateResgisterUserRequest(&req)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, util.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error: setting password hash: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error: creating user: %v", err)
		util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal server error"})
		return
	}
	util.WriteJSON(w, http.StatusCreated, util.Envelope{"user": user})

}
