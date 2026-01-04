package graph

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

// Middleware HTTPリクエストからUserIDを抽出してContextに入れる
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			// ヘッダーがない場合はそのまま次へ（未ログイン状態）
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// "Bearer 1" の "1" の部分を取得
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				uid, _ := strconv.Atoi(parts[1])
				// Contextに値をセット
				ctx := context.WithValue(r.Context(), userIDKey, uint(uid))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserUIDFromContext ResolverなどでContextからIDを取り出すためのヘルパー
func GetUserUIDFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(userIDKey).(string)
	return uid, ok
}
