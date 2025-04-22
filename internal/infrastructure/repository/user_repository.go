package repository

import (
	"context"
	"log"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
)

// メールアドレスでユーザーを検索する
func GetUserByEmail(email string) (*domain.User, error) {
	client, err := firebase.InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	query := client.Collection("users").Where("Email", "==", email)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Firestoreクエリエラー: %v", err)
		return nil, err
	}

	// ユーザーが見つからない場合
	if len(docs) == 0 {
		log.Printf("ユーザーが見つかりません: email=%s", email)
		return nil, nil
	}

	var user domain.User
	if err := docs[0].DataTo(&user); err != nil {
		log.Printf("ユーザーデータ変換エラー: %v", err)
		return nil, err
	}

	// ドキュメントIDをユーザーIDとして設定
	user.ID = docs[0].Ref.ID
	return &user, nil
}

// ユーザーIDからユーザー情報を取得する
func GetUserByID(userID string) (*domain.User, error) {
	client, err := firebase.InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	doc, err := client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	user.ID = doc.Ref.ID

	return &user, nil
}
