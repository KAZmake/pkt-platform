package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/pkg/response"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type KeycloakClaims struct {
	jwt.RegisteredClaims
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}

type ctxKey string

const claimsKey ctxKey = "jwt_claims"

func ClaimsFromCtx(ctx context.Context) (*KeycloakClaims, bool) {
	c, ok := ctx.Value(claimsKey).(*KeycloakClaims)
	return c, ok && c != nil
}

func Authenticate(jwks keyfunc.Keyfunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				response.Unauthorized(w)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims := &KeycloakClaims{}
			_, err := jwt.ParseWithClaims(tokenStr, claims, jwks.KeyfuncCtx(r.Context()))
			if err != nil {
				response.Unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromCtx(r.Context())
			if !ok {
				response.Forbidden(w)
				return
			}
			for _, role := range claims.RealmAccess.Roles {
				if allowed[role] {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.Forbidden(w)
		})
	}
}
