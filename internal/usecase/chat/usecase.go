package chat

import (
	"context"
	"fmt"
	"time"

	"security_chat_app/internal/domain"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// チャットのコントローラー
type ChatController struct {
	chatUsecase domain.ChatUsecase
}

// チャットのリポジトリの実装
type chatRepository struct {
	// *: メソッド内で構造体を変更する
	client *firestore.Client
}

// チャットのユースケースの実装
type chatUsecaseImpl struct {
	repo domain.ChatRepository
}

// **************************************************
// main.go で使用するメソッド **************
// **************************************************

// チャットのリポジトリを生成する
func NewChatRepository(client *firestore.Client) domain.ChatRepository {
	return &chatRepository{client: client}
}

// チャットのユースケースの実装を生成する
func NewChatUsecase(repo domain.ChatRepository) domain.ChatUsecase {
	return &chatUsecaseImpl{repo: repo}
}


// **************************************************
// ChatUsecaseの定義 **************
// **************************************************

// チャット作成時のビジネスロジックを定義
func (c *chatUsecaseImpl) CreateChat(user, message string) error {
	return c.repo.AddChat(user, message)
}

// GetChatHistoryメソッドの実装
func (c *chatUsecaseImpl) GetChatHistory(user *domain.User) ([]domain.Chat, error) {
	return c.repo.GetChats(user.ID)
}

// GetContactsメソッドの実装
func (c *chatUsecaseImpl) GetContacts(user *domain.User) ([]domain.Contact, error) {
	// TODO: 実装
	return nil, fmt.Errorf("not implemented")
}

// **************************************************
// ChatRepositoryの定義 **************
// **************************************************

// AddChatメソッドの実装
func (r *chatRepository) AddChat(user, message string) error {
	// Firestoreにチャットを追加する処理
	_, _, err := r.client.Collection("chats").Add(context.Background(), map[string]interface{}{
		"user":       user,
		"message":    message,
		"created_at": time.Now(),
	})
	return err
}

// GetChatsメソッドの実装
func (r *chatRepository) GetChats(userID string) ([]domain.Chat, error) {
	// Firestoreからチャットを取得する処理
	chats := []domain.Chat{}
	iter := r.client.Collection("chats").Where("user", "==", userID).Documents(context.Background())

	// Firestoreから取得した生のデータを、アプリケーションで使用しやすい domain.Chat 構造体の形式に変換
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var chat domain.Chat
		// Firestoreのデータをドメインの構造体に変換
		if err := doc.DataTo(&chat); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

// GetMessagesメソッドの実装
func (r *chatRepository) GetMessages(chatID string) ([]domain.Message, error) {
	// Firestoreからメッセージを取得する処理
	messages := []domain.Message{}
	iter := r.client.Collection("chats").Doc(chatID).Collection("messages").Documents(context.Background())

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var message domain.Message
		if err := doc.DataTo(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// **************************************************
// ChatControllerの定義 **************
// **************************************************

// HandleCreateChatメソッドの実装
func (c *ChatController) HandleCreateChat(user, message string) error {
	return c.chatUsecase.CreateChat(user, message)
}

// HandleGetChatHistoryメソッドの実装
func (c *ChatController) HandleGetChatHistory(user *domain.User) ([]domain.Chat, error) {
	return c.chatUsecase.GetChatHistory(user)
}

// HandleGetContactsメソッドの実装
func (c *ChatController) HandleGetContacts(user *domain.User) ([]domain.Contact, error) {
	return c.chatUsecase.GetContacts(user)
}
