package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
)


// セッションを検証
func ValidateSession(w http.ResponseWriter, r *http.Request) (*domain.Session, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("セッションクッキー取得エラー: %v", err)
		return nil, err
	}

	sessionID := cookie.Value

	// Firestoreからセッションを取得
	client, err := firebase.InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
		return nil, err
	}
	defer client.Close()

	ctx := r.Context()
	doc, err := client.Collection("sessions").Doc(sessionID).Get(ctx)
	if err != nil {
		log.Printf("セッション取得エラー: %v, sessionID=%s", err, sessionID)
		return nil, err
	}

	var session domain.Session
	if err := doc.DataTo(&session); err != nil {
		log.Printf("セッションデータ変換エラー: %v", err)
		return nil, err
	}

	if !session.CheckSession() {
		log.Printf("セッションが無効です: sessionID=%s", sessionID)
		return nil, fmt.Errorf("セッションが無効です")
	}
	return &session, nil
}

// セッションを作成
func CreateSession(user *domain.User) (*domain.Session, error) {
	// セッションIDの生成
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	sessionID := base64.URLEncoding.EncodeToString(bytes)

	// セッションの作成
	session := &domain.Session{
		ID:        sessionID,                           // セッションID
		User:      user,                                // ユーザー
		Token:     sessionID,                           // セッショントークン
		CreatedAt: time.Now(),                          // セッションの作成日時
		UpdatedAt: time.Now(),                          // セッションの更新日時
		ExpiredAt: time.Now().Add(30 * 24 * time.Hour), // 30日間有効
		IsValid:   true,                                // セッションが有効かどうか
	}

	// Firestoreにセッションを保存
	err := firebase.AddData("sessions", session, sessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// セッションクッキーを設定
func SetSessionCookie(w http.ResponseWriter, session *domain.Session) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                // 開発環境ではfalseに設定
		SameSite: http.SameSiteLaxMode, // 開発環境ではLaxに設定
		MaxAge:   86400 * 30,           // 30日
	}
	http.SetCookie(w, cookie)
}

// セッションを更新
func UpdateSession(w http.ResponseWriter, r *http.Request, session *domain.Session) error {
	// Firestoreにセッションを保存（セッションIDをドキュメントIDとして使用）
	err := firebase.AddData("sessions", session, session.ID)
	if err != nil {
		return err
	}

	// セッションクッキーを更新
	SetSessionCookie(w, session)

	return nil
}

// セッションを削除
func DeleteSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return err
	}

	// Firestoreからセッションを削除
	client, err := firebase.InitFirebase()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := r.Context()
	_, err = client.Collection("sessions").Doc(cookie.Value).Delete(ctx)
	if err != nil {
		return err
	}

	// クッキーを削除
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	return nil
}
