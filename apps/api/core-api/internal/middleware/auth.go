package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/pkg/response"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	claimsKey contextKey = "claims"
)

// KeycloakClaims represents the JWT claims issued by Keycloak.
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

// HasRole returns true if the claims contain the given role.
func (c *KeycloakClaims) HasRole(role string) bool {
	for _, r := range c.RealmAccess.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// Authenticate validates the Bearer JWT from the Authorization header.
// On success it stores *KeycloakClaims in the request context.
// On failure it responds 401 and aborts the chain.
func Authenticate(jwks keyfunc.Keyfunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r)
			if token == "" {
				response.Unauthorized(w)
				return
			}

			claims := &KeycloakClaims{}
			_, err := jwt.ParseWithClaims(token, claims, jwks.Keyfunc,
				jwt.WithValidMethods([]string{"RS256"}),
				jwt.WithExpirationRequired(),
			)
			if err != nil {
				response.Unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole returns 403 if the authenticated user does not hold one of the
// given roles. Must be used after Authenticate.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromCtx(r.Context())
			if !ok {
				response.Unauthorized(w)
				return
			}
			for _, role := range roles {
				if claims.HasRole(role) {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.Forbidden(w)
		})
	}
}

// ClaimsFromCtx extracts KeycloakClaims from the request context.
func ClaimsFromCtx(ctx context.Context) (*KeycloakClaims, bool) {
	c, ok := ctx.Value(claimsKey).(*KeycloakClaims)
	return c, ok
}

// bearerToken extracts the token string from "Authorization: Bearer <token>".
func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(h, "Bearer ")
}
