package main

import (
	"log"
	"net/http"

	"security_chat_app/internal/config"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/infrastructure/router"
	"security_chat_app/internal/usecase/chat"
)

func main() {
	// Firebaseの初期化
	client, err := firebase.InitFirebase()
	if err != nil {
		log.Fatalf("Firebase初期化に失敗: %v", err)
	}
	defer client.Close()

	// チャットリポジトリの作成
	chatRepo := chat.NewChatRepository(client)
	chatUsecase := chat.NewChatUsecase(chatRepo)
	if chatUsecase == nil {
		log.Fatal("チャットのユースケースの実装に不備があります")
	}

	// ルーティングの設定
	httpRouter := router.SetupRouter(chatUsecase)
	if httpRouter == nil {
		log.Fatal("ルーティングの設定に不備があります")
	}

	// サーバーを起動
	log.Printf("サーバーを起動します。ポート: %s", config.Config.Port)
	if err := http.ListenAndServe(":"+config.Config.Port, httpRouter); err != nil {
		log.Fatal("サーバーの起動に失敗しました")
	}
}
