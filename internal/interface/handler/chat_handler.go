package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/infrastructure/repository"
	"security_chat_app/internal/interface/markup"
	"security_chat_app/internal/interface/middleware"
)

// チャット開始ハンドラ
func StartChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Fatalf("メソッドが許可されていません")
		return
	}

	// セッションの検証
	session, err := middleware.ValidateSession(w, r)
	if err != nil {
		log.Fatalf("認証されていません: %v", err)
		return
	}

	// セッションからユーザー情報を取得
	user, err := repository.GetUserByID(session.User.ID)
	if err != nil {
		log.Fatalf("ユーザー情報の取得に失敗: %v", err)
		return
	}

	// URLから対象ユーザーIDを取得
	targetUserID := r.URL.Path[len("/chat/"):]
	if targetUserID == "" {
		log.Fatalf("ユーザーIDが指定されていません")
		return
	}

	// 対象ユーザーの存在確認
	_, err = GetUserData(targetUserID)
	if err != nil {
		log.Fatalf("対象ユーザーが見つかりません: %v", err)
		return
	}

	// チャットを開始
	chatID, err := firebase.StartChat(user.ID, targetUserID)
	if err != nil {
		log.Fatalf("チャットの開始に失敗: %v", err)
		return
	}

	// チャットページにリダイレクト
	redirectURL := fmt.Sprintf("/chat?chat_id=%s", chatID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// チャットページのハンドラ
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの検証
	session, err := middleware.ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// セッションからユーザー情報を取得
	user, err := repository.GetUserByID(session.User.ID)
	if err != nil {
		log.Fatalf("ユーザー情報の取得に失敗: %v", err)
		return
	}

	// POSTリクエストの場合はメッセージ送信処理
	if r.Method == http.MethodPost {
		// フォームデータから情報を取得
		chatID := r.FormValue("chatID")
		content := r.FormValue("content")

		if chatID == "" || content == "" {
			log.Fatalf("チャットIDとメッセージ内容が必要です")
			return
		}

		// メッセージIDを生成
		messageID := generateMessageID()

		// メッセージを作成
		message := map[string]interface{}{
			"id":         messageID,
			"sender_id":  user.ID,
			"sender_name": user.Name,
			"content":    content,
			"created_at": time.Now(),
			"is_read":    false,
			"type":       "text",
		}

		// メッセージを保存
		err = firebase.AddChatMessage(chatID, message)
		if err != nil {
			log.Fatalf("メッセージの送信に失敗: %v", err)
			return
		}

		// チャットの最終更新時刻を更新
		err = firebase.UpdateField("chats", chatID, "updated_at", time.Now())
		if err != nil {
			log.Printf("チャットの更新時刻の更新に失敗: %v", err)
		}

		// JSONレスポンスを返す
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         messageID,
			"content":    content,
			"sender_id":  user.ID,
			"sender_name": user.Name,
			"created_at": time.Now().Format("15:04"),
			"is_read":    false,
		})
		return
	}

	// チャット履歴を取得
	chats, err := getChatHistory(user)
	if err != nil {
		log.Fatalf("チャット一覧の取得に失敗: %v", err)
		return
	}

	// URLからチャットIDを取得
	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		// チャットIDがない場合は、チャット一覧を表示
		data := domain.TemplateData{
			IsLoggedIn: true,
			User:       user,
			Chats:      chats,
			ChatID:     "", // 空のチャットIDを設定
		}
		markup.GenerateHTML(w, data, "layout", "header", "chat", "footer")
		return
	}

	// チャットの存在確認
	exists, err := firebase.CheckChatExists(chatID)
	if err != nil {
		log.Fatalf("チャットの確認に失敗: %v", err)
		return
	}
	if !exists {
		log.Fatalf("チャットが見つかりません")
		return
	}

	// チャットの参加者を取得
	participants, err := firebase.GetChatParticipants(chatID)
	if err != nil {
		log.Fatalf("チャットの参加者情報の取得に失敗: %v", err)
		return
	}

	// 対象ユーザーを特定
	var targetUserID string
	for _, p := range participants {
		if p != user.ID {
			targetUserID = p
			break
		}
	}

	// 対象ユーザーの存在確認
	_, err = GetUserData(targetUserID)
	if err != nil {
		log.Fatalf("対象ユーザーが見つかりません: %v", err)
		return
	}

	// 対象ユーザーの情報を取得
	targetUser, err := GetUserData(targetUserID)
	if err != nil {
		log.Fatalf("対象ユーザーの情報の取得に失敗: %v", err)
		return
	}

	// メッセージを取得
	messagesData, err := firebase.GetChatMessages(chatID)
	if err != nil {
		log.Fatalf("メッセージの取得に失敗: %v", err)
		return
	}

	// メッセージの型変換
	var messages []domain.Message
	for _, msg := range messagesData {
		// 各フィールドを安全に取得する関数
		getString := func(key string) string {
			if val, exists := msg[key]; exists && val != nil {
				if str, ok := val.(string); ok {
					return str
				}
			}
			// 大文字のキーも試す
			if val, exists := msg[strings.ToUpper(key)]; exists && val != nil {
				if str, ok := val.(string); ok {
					return str
				}
			}
			return ""
		}

		// 時刻の取得
		var createdAt time.Time
		if t, ok := msg["created_at"].(time.Time); ok {
			createdAt = t
		} else if t, ok := msg["CreatedAt"].(time.Time); ok {
			createdAt = t
		} else {
			createdAt = time.Now() // デフォルト値
		}

		// 既読状態の取得
		isRead := false
		if r, ok := msg["is_read"].(bool); ok {
			isRead = r
		} else if r, ok := msg["IsRead"].(bool); ok {
			isRead = r
		}

		message := domain.Message{
			ID:         getString("id"),
			Content:    getString("content"),
			SenderID:   getString("sender_id"),
			SenderName: getString("sender_name"),
			CreatedAt:  createdAt,
			IsRead:     isRead,
		}
		messages = append(messages, message)
	}

	// 現在のチャットを特定
	var currentChat *domain.Chat
	for _, chat := range chats {
		if chat.ID == chatID {
			currentChat = &chat
			break
		}
	}

	// チャットページのデータを取得
	data := domain.TemplateData{
		IsLoggedIn:  true,
		User:        user,
		Messages:    messages,
		Contacts:    []domain.Contact{{ID: targetUser.ID, Username: targetUser.Name, Icon: targetUser.Icon}},
		Chats:       chats,
		CurrentChat: currentChat,
		ChatID:      chatID,
	}

	// テンプレートのレンダリング
	markup.GenerateHTML(w, data, "layout", "header", "chat", "footer")
}

// チャット履歴を取得
func getChatHistory(user *domain.User) ([]domain.Chat, error) {
	// チャット履歴を取得
	chats, err := firebase.GetAllChats(user.ID)
	if err != nil {
		return nil, fmt.Errorf("チャット履歴の取得に失敗しました: %v", err)
	}

	var chatHistory []domain.Chat
	seenChats := make(map[string]bool) // 重複チェック用のマップ

	for _, chatData := range chats {
		// チャットIDの取得
		chatID, ok := chatData["id"].(string)
		if !ok {
			log.Printf("チャットIDの取得に失敗: %v", chatData)
			continue
		}

		// 参加者の取得
		participants, ok := chatData["participants"].([]interface{})
		if !ok || len(participants) != 2 {
			continue
		}

		// 自分が参加者に含まれているか確認
		isParticipant := false
		var targetUserID string
		for _, p := range participants {
			participantID, ok := p.(string)
			if !ok {
				continue
			}
			if participantID == user.ID {
				isParticipant = true
			} else {
				targetUserID = participantID
			}
		}

		// 自分が参加者でない場合はスキップ
		if !isParticipant || targetUserID == "" || seenChats[chatID] {
			continue
		}

		seenChats[chatID] = true

		// メッセージの取得
		messagesData, err := firebase.GetChatMessages(chatID)
		if err != nil {
			log.Printf("メッセージの取得に失敗: chatID=%s, error=%v", chatID, err)
			continue
		}

		// メッセージの型変換
		var messages []domain.Message
		var lastMessageTime time.Time
		for _, msg := range messagesData {
			// 各フィールドを安全に取得する関数
			getString := func(key string) string {
				if val, exists := msg[key]; exists && val != nil {
					if str, ok := val.(string); ok {
						return str
					}
				}
				// 大文字のキーも試す
				if val, exists := msg[strings.ToUpper(key)]; exists && val != nil {
					if str, ok := val.(string); ok {
						return str
					}
				}
				return ""
			}

			// 時刻の取得
			var createdAt time.Time
			if t, ok := msg["created_at"].(time.Time); ok {
				createdAt = t
			} else if t, ok := msg["CreatedAt"].(time.Time); ok {
				createdAt = t
			} else {
				createdAt = time.Now() // デフォルト値
			}

			// 既読状態の取得
			isRead := false
			if r, ok := msg["is_read"].(bool); ok {
				isRead = r
			} else if r, ok := msg["IsRead"].(bool); ok {
				isRead = r
			}

			message := domain.Message{
				ID:         getString("id"),
				Content:    getString("content"),
				SenderID:   getString("sender_id"),
				SenderName: getString("sender_name"),
				CreatedAt:  createdAt,
				IsRead:     isRead,
			}
			messages = append(messages, message)

			// 最新のメッセージ時刻を更新
			if createdAt.After(lastMessageTime) {
				lastMessageTime = createdAt
			}
		}

		// チャット相手の情報を取得
		targetUser, err := GetUserData(targetUserID)
		if err != nil {
			log.Printf("チャット相手の情報取得に失敗: targetUserID=%s, error=%v", targetUserID, err)
			continue
		}

		// チャット履歴に追加
		chatHistory = append(chatHistory, domain.Chat{
			ID: chatID,
			Contact: domain.Contact{
				ID:       targetUser.ID,
				Username: targetUser.Name,
				Icon:     targetUser.Icon,
				LastSeen: time.Now(),
				IsOnline: targetUser.IsOnline,
			},
			Messages:  messages,
			UpdatedAt: lastMessageTime,
		})
	}

	// 更新時刻でソート（新しい順）
	sort.Slice(chatHistory, func(i, j int) bool {
		return chatHistory[i].UpdatedAt.After(chatHistory[j].UpdatedAt)
	})

	return chatHistory, nil
}

// ChatControllerImpl チャットコントローラーの実装
type ChatControllerImpl struct {
	chatUsecase domain.ChatUsecase
}

// NewChatController チャットコントローラーを作成する
func NewChatController(chatUsecase domain.ChatUsecase) *ChatControllerImpl {
	return &ChatControllerImpl{
		chatUsecase: chatUsecase,
	}
}

// Create チャットを作成する
func (c *ChatControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Fatalf("メソッドが許可されていません")
		return
	}

	user := r.FormValue("user")
	message := r.FormValue("message")

	if err := c.chatUsecase.CreateChat(user, message); err != nil {
		log.Fatalf("チャットの作成に失敗: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// メッセージ送信ハンドラ
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Fatalf("メソッドが許可されていません")
		return
	}

	// セッションの検証
	session, err := middleware.ValidateSession(w, r)
	if err != nil {
		log.Fatalf("認証されていません: %v", err)
		return
	}

	// セッションからユーザー情報を取得
	user, err := repository.GetUserByID(session.User.ID)
	if err != nil {
		log.Fatalf("ユーザー情報の取得に失敗: %v", err)
		return
	}

	// フォームデータから情報を取得
	chatID := r.FormValue("chatID")
	content := r.FormValue("content")

	if chatID == "" || content == "" {
		log.Fatalf("チャットIDとメッセージ内容が必要です")
		return
	}

	// メッセージを作成
	message := map[string]interface{}{
		"sender_id":   user.ID,
		"sender_name": user.Name,
		"content":     content,
		"created_at":  time.Now(),
		"is_read":     false,
	}

	// メッセージを保存
	err = firebase.AddChatMessage(chatID, message)
	if err != nil {
		log.Fatalf("メッセージの送信に失敗: %v", err)
		return
	}

	// チャットページにリダイレクト
	http.Redirect(w, r, fmt.Sprintf("/chat?chat_id=%s", chatID), http.StatusSeeOther)
}

// メッセージIDを生成する
func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}
