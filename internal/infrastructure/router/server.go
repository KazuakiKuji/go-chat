package router

import (
	"net/http"

	"security_chat_app/internal/config"
	"security_chat_app/internal/domain"
)

// メインサーバーを起動する
func StartMainServer(chatUsecase domain.ChatUsecase) error {
	mux := SetupRouter(chatUsecase)
	return http.ListenAndServe(":"+config.Config.Port, mux)
}
