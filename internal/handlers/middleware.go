package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/mwdev22/FileStorage/internal/config"
	"github.com/mwdev22/FileStorage/internal/types"
)

type ctxKey string

const userIdKey ctxKey = "user_id"

type apiFunc func(w http.ResponseWriter, r *http.Request) error

// middleware on default http handler
// allows to assert occured error as one of our known errors (check errors.go)
func makeHandler(h apiFunc, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		if err := h(w, r); err != nil {
			// assert the error from handler to known error type
			if e, ok := err.(ApiError); ok {
				types.WriteJSON(w, e.StatusCode, map[string]string{
					"error": e.Error(),
				})
			} else {
				types.WriteJSON(w, http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
			}
		}
	}
}

// auth middleware for private routes, if authorized returns handler made of makeHandler
func authMiddleware(h apiFunc, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "missing Authorization header",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid Authorization header format",
			})
			return
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return config.SecretKey, nil
		})

		if err != nil || !token.Valid {
			types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid or expired token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid token claims",
			})
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "user ID not found in token",
			})
		}

		ctx := context.WithValue(r.Context(), userIdKey, int(userID))
		r = r.WithContext(ctx)

		// if authenticated, make handler for the request
		makeHandler(h, logger)(w, r)
	}
}
