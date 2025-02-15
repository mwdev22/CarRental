package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/mwdev22/CarRental/internal/config"
	"github.com/mwdev22/CarRental/internal/types"
)

type ctxKey string

const (
	userIdKey ctxKey = "user_id"
	roleKey   ctxKey = "role"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

// base handler func to handle errors
func makeHandler(h apiFunc, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		if err := h(w, r); err != nil {
			if e, ok := err.(types.ApiError); ok {
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

// Middleware to authenticate user and extract user ID
func authMiddleware(h apiFunc, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := parseToken(r.Header.Get("Authorization"))
		if err != nil {
			types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "user ID not found in token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, int(userID))
		r = r.WithContext(ctx)

		// If authenticated, process the request
		makeHandler(h, logger)(w, r)
	}
}

// Middleware to authenticate and extract role-specific claims
func roleMiddleware(h apiFunc, role types.UserRole, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := parseToken(r.Header.Get("Authorization"))
		if err != nil {
			types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "user ID not found in token",
			})
			return
		}

		roleFloat, ok := claims["role"].(float64)
		if !ok {
			http.Error(w, "invalid role type", http.StatusInternalServerError)
			return
		}

		tokenRole := types.UserRole(int(roleFloat))
		if tokenRole != role && tokenRole != types.UserTypeAdmin {
			types.WriteJSON(w, http.StatusForbidden, map[string]string{
				"error": "forbidden: insufficient role permissions",
				"role":  fmt.Sprintf("%v", tokenRole),
			})
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, int(userID))
		ctx = context.WithValue(ctx, roleKey, role)
		r = r.WithContext(ctx)

		makeHandler(h, logger)(w, r)
	}
}

func parseToken(authHeader string) (jwt.MapClaims, error) {
	if authHeader == "" {
		return nil, fmt.Errorf("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid Authorization header format")
	}

	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return config.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
