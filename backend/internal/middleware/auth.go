package middleware

import (
	"context"
	"net/http"
	"strings"

	"EMS/internal/cognito"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey = contextKey("userSub")

func Auth(verifier *cognito.JWTVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions || r.URL.Path == "/api/" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			// Verify token using Cognito verifier
			token, err := verifier.VerifyAccessToken(tokenStr)
			if err != nil {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Unauthorized: No sub in token", http.StatusUnauthorized)
				return
			}

			// Add to context
			ctx := context.WithValue(r.Context(), UserContextKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
