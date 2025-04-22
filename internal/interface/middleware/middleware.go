package middleware

import (
	"context"
	"log"
	"net/http"
	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
)

// コンテキストのキーとして使用するカスタム型
type contextKey string

// テンプレートデータのキー
const templateDataKey contextKey = "templateData"

// セッション管理のミドルウェア
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := ValidateSession(w, r)
		if err != nil {
			// セッションが無効な場合は、ログインしていない状態として処理
			data := domain.TemplateData{IsLoggedIn: false}
			r = r.WithContext(context.WithValue(r.Context(), templateDataKey, data))
		} else {
			// セッションが有効な場合は、ログインしている状態として処理
			data := domain.TemplateData{
				IsLoggedIn: true,
				User:       session.User,
			}
			r = r.WithContext(context.WithValue(r.Context(), templateDataKey, data))
			
			// Firebaseのユーザー状態をオンラインに更新
			if err := firebase.UpdateField("users", session.User.ID, "IsOnline", true); err != nil {
				log.Printf("ユーザー状態の更新に失敗: %v", err)
			}
		}
		next.ServeHTTP(w, r)
	})
}
