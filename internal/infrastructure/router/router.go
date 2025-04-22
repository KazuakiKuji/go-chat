package router

import (
	"net/http"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/interface/handler"
	"security_chat_app/internal/interface/middleware"
)

// ルーティングの設定
func SetupRouter(chatUsecase domain.ChatUsecase) *http.ServeMux {
	rootDir := "internal/web/"
	httpRouter := http.NewServeMux()
	// 静的ファイル (CSS/JS)
	httpRouter.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(rootDir+"css"))))
	httpRouter.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(rootDir+"js"))))
	httpRouter.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(rootDir+"images"))))
	// ルーティング
	httpRouter.Handle("/", middleware.Middleware(http.HandlerFunc(handler.SearchHandler)))
	httpRouter.Handle("/login", http.HandlerFunc(handler.LoginHandler))
	httpRouter.Handle("/logout", http.HandlerFunc(handler.LogoutHandler))
	httpRouter.Handle("/signup", http.HandlerFunc(handler.SignupHandler))
	httpRouter.Handle("/signup/confirm", http.HandlerFunc(handler.SignupConfirmHandler))
	httpRouter.Handle("/reset-password", http.HandlerFunc(handler.ResetPasswordHandler))
	httpRouter.Handle("/profile", middleware.Middleware(http.HandlerFunc(handler.ProfileHandler)))
	httpRouter.Handle("/profile/", middleware.Middleware(http.HandlerFunc(handler.ProfileHandler)))
	httpRouter.Handle("/profile/icon", middleware.Middleware(http.HandlerFunc(handler.ProfileIconHandler)))
	httpRouter.Handle("/chat/", middleware.Middleware(http.HandlerFunc(handler.StartChatHandler)))
	httpRouter.Handle("/chat", middleware.Middleware(http.HandlerFunc(handler.ChatHandler)))
	httpRouter.Handle("/search", middleware.Middleware(http.HandlerFunc(handler.SearchHandler)))
	httpRouter.Handle("/settings", middleware.Middleware(http.HandlerFunc(handler.SettingsHandler)))
	httpRouter.Handle("/settings/username", middleware.Middleware(http.HandlerFunc(handler.SettingsHandler)))

	return httpRouter
}
