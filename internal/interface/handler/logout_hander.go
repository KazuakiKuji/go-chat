package handler

import (
	"log"
	"net/http"

	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/interface/middleware"
)

// ログアウト処理を実行
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		session, err := middleware.ValidateSession(w, r)
		if err == nil && session != nil && session.User != nil {
			if err := firebase.UpdateField("users", session.User.ID, "IsOnline", false); err != nil {
				log.Printf("ユーザー状態の更新に失敗: %v", err)
			}
		}
		
		err = middleware.DeleteSession(w, r)
		if err != nil {
			log.Fatalf("ログアウトエラー: %v", err)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
