package firebase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// コレクションにデータを追加する
func AddData(collection string, data interface{}, docID string) error {
	client, err := InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
		return err
	}
	defer client.Close()

	ctx := context.Background()
	if docID != "" {
		_, err = client.Collection(collection).Doc(docID).Set(ctx, data)
	} else {
		_, _, err = client.Collection(collection).Add(ctx, data)
	}
	if err != nil {
		log.Printf("データ追加エラー: %v", err)
		return err
	}
	return nil
}

// コレクションとドキュメントIDから特定フィールドを更新する
func UpdateField(collection string, documentID string, field string, value interface{}) error {
	client, err := InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
		return err
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Collection(collection).Doc(documentID).Update(ctx, []firestore.Update{
		{
			Path:  field,
			Value: value,
		},
	})
	if err != nil {
		log.Printf("フィールド更新エラー: %v, collection=%s, documentID=%s, field=%s", err, collection, documentID, field)
		return err
	}
	return nil
}

// コレクションからデータを取得する
func GetData(collection string, documentID string) (map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	doc, err := client.Collection(collection).Doc(documentID).Get(ctx)
	if err != nil {
		return nil, err
	}

	return doc.Data(), nil
}

// コレクションから条件に合うデータを取得する
func GetDataByQuery(collection string, field string, operator string, value interface{}) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	query := client.Collection(collection).Where(field, operator, value)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, doc := range docs {
		results = append(results, doc.Data())
	}

	return results, nil
}

// コレクションからデータを削除する
func DeleteData(collection string, documentID string) error {
	client, err := InitFirebase()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Collection(collection).Doc(documentID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

// コレクションの全データを取得する
func GetAllData(collection string, userID string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	
	var docs []*firestore.DocumentSnapshot
	var err2 error
	
	if collection == "chats" {
		query := client.Collection(collection).Where("participants", "array-contains", userID)
		docs, err2 = query.Documents(ctx).GetAll()
	} else {
		docs, err2 = client.Collection(collection).Documents(ctx).GetAll()
	}
	
	if err2 != nil {
		return nil, err2
	}

	var results []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		data["id"] = doc.Ref.ID
		results = append(results, data)
	}

	return results, nil
}

// ユーザーを検索する
func SearchUser(searchQuery string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()

	// すべてのユーザーを取得
	usersRef := client.Collection("users")
	docs, err := usersRef.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("ユーザー検索エラー: %v", err)
		return nil, err
	}

	// 検索クエリを小文字に変換
	searchQueryLower := strings.ToLower(searchQuery)

	var results []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		data["ID"] = doc.Ref.ID

		// ユーザー名を取得
		name, ok := data["Name"].(string)
		if !ok {
			continue
		}

		// 大文字小文字を区別せずに部分一致検索
		if strings.Contains(strings.ToLower(name), searchQueryLower) {
			results = append(results, data)
		}
	}

	return results, nil
}

// チャットを開始する
func StartChat(userID string, targetUserID string) (string, error) {
	client, err := InitFirebase()
	if err != nil {
		return "", err
	}
	defer client.Close()

	ctx := context.Background()

	// チャットIDを生成
	chatID := fmt.Sprintf("chat_%d", time.Now().UnixNano())

	// チャットを作成
	chat := map[string]interface{}{
		"id":           chatID,
		"participants": []string{userID, targetUserID},
		"createdAt":    time.Now(),
		"updatedAt":    time.Now(),
	}

	_, err = client.Collection("chats").Doc(chatID).Set(ctx, chat)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

// チャットメッセージを追加する
func AddChatMessage(chatID string, message map[string]interface{}) error {
	client, err := InitFirebase()
	if err != nil {
		log.Printf("Firebase初期化エラー: %v", err)
		return err
	}
	defer client.Close()

	ctx := context.Background()

	// メッセージIDを生成
	messageID := fmt.Sprintf("msg_%d", time.Now().UnixNano())
	message["id"] = messageID

	// メッセージを保存
	_, err = client.Collection("chats").Doc(chatID).Collection("messages").Doc(messageID).Set(ctx, message)
	if err != nil {
		log.Printf("メッセージ保存エラー: %v", err)
		return err
	}

	// チャットの更新時刻を更新
	_, err = client.Collection("chats").Doc(chatID).Update(ctx, []firestore.Update{
		{
			Path:  "updated_at",
			Value: time.Now(),
		},
	})
	if err != nil {
		log.Printf("チャット更新時刻の更新エラー: %v", err)
		return err
	}

	return nil
}

// チャットのメッセージを取得する
func GetChatMessages(chatID string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	docs, err := client.Collection("chats").Doc(chatID).Collection("messages").OrderBy("created_at", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var messages []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		data["id"] = doc.Ref.ID
		messages = append(messages, data)
	}

	return messages, nil
}

// チャットの存在確認
func CheckChatExists(chatID string) (bool, error) {
	client, err := InitFirebase()
	if err != nil {
		return false, err
	}
	defer client.Close()

	ctx := context.Background()
	doc, err := client.Collection("chats").Doc(chatID).Get(ctx)
	if err != nil {
		return false, err
	}

	return doc.Exists(), nil
}

// チャットの参加者を取得
func GetChatParticipants(chatID string) ([]string, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()
	doc, err := client.Collection("chats").Doc(chatID).Get(ctx)
	if err != nil {
		return nil, err
	}

	data := doc.Data()
	participants, ok := data["participants"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("participants field is invalid")
	}

	var result []string
	for _, p := range participants {
		if str, ok := p.(string); ok {
			result = append(result, str)
		}
	}

	return result, nil
}

// 指定されたユーザーIDが参加者として含まれるチャットを全て取得します
func GetAllChats(userID string) ([]map[string]interface{}, error) {
	client, err := InitFirebase()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()

	// チャットコレクションを参照
	chatsRef := client.Collection("chats")

	// ユーザーIDが参加者に含まれるチャットを検索
	query := chatsRef.Where("participants", "array-contains", userID)
	iter := query.Documents(ctx)
	defer iter.Stop()

	var chats []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("チャットデータの取得に失敗: %v", err)
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		chats = append(chats, data)
	}

	return chats, nil
}
