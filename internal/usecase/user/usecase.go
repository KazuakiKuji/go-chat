package user

import (
	"context"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	utils "security_chat_app/internal/utils/uuid"

	"golang.org/x/crypto/bcrypt"
)

// ユーザー登録
func CreateUser(name, email, password string) (*domain.User, error) {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// ユーザーを作成
	user := &domain.User{
		ID:        utils.GenerateUUID(), // UUIDを生成する関数は別途実装が必要
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Firestoreにユーザーを保存
	client, err := firebase.InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
