package middleware

import (
	"context"
	"net/http"
	"strings"

	jwtutil "pet-shelter/pkg/jwt"
)

type contextKey string

const claimsContextKey contextKey = "jwt_claims"

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := tokenFromRequest(r)
			if tokenString == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			claims, err := jwtutil.ValidateToken(tokenString, m.jwtSecret)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			if requiredRole != "" && claims.Role != requiredRole {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*jwtutil.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*jwtutil.Claims)
	return claims, ok
}

func tokenFromRequest(r *http.Request) string {
	if cookie, err := r.Cookie("access_token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	const prefix = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, prefix) {
		return strings.TrimPrefix(authHeader, prefix)
	}

	return ""
}
