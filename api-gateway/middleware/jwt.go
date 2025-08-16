package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RibunLoc/WebPersonalBackend/api-gateway/util"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret string
}

type ctxKey string

const userIDKey ctxKey = "user_id"

func (m JWT) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			util.Error(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (any, error) {
			return []byte(m.Secret), nil
		})
		if err != nil || !token.Valid {
			util.Error(w, http.StatusUnauthorized, "invalid token")
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			util.Error(w, http.StatusUnauthorized, "invalid claims")
			return
		}
		uid, _ := claims["user_id"].(string)
		ctx := context.WithValue(r.Context(), userIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper lấy user_id từ context trong handler sau này
func UserIDFromCtx(r *http.Request) (string, bool) {
	uid, ok := r.Context().Value(userIDKey).(string)
	return uid, ok
}
