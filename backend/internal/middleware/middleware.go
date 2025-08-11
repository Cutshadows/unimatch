package middleware

import (
	"context"
	"net/http"
	"strings"
	"unimatch-back/internal/store"
	"unimatch-back/internal/tokens"
	"unimatch-back/util"
)

type Middleware struct {
	UserStore store.UserStore
}
type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("user not found in request")
	}
	return user
}

func (um *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, "")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "invalid authorization header format"})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "invalid token"})
			return
		}
		if user == nil {
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "token expired or invalid"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *Middleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "user not authenticated"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
