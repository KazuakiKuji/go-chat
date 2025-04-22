package firebase

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"security_chat_app/internal/config"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// アイコンをアップロードする
func UploadIcon(userID string, filePath string) (string, error) {
	opt := option.WithCredentialsFile(config.Config.ServiceKeyPath)
	config := &firebase.Config{
		ProjectID:     config.Config.ProjectId,
		StorageBucket: config.Config.StorageBucket,
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return "", fmt.Errorf("firebaseアプリの初期化に失敗しました: %v", err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		return "", fmt.Errorf("storageクライアントの作成に失敗しました: %v", err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", fmt.Errorf("バケットの取得に失敗しました: %v", err)
	}

	// アップロード先のパスを設定
	objectPath := fmt.Sprintf("icons/%s%s", userID, filepath.Ext(filePath))
	object := bucket.Object(objectPath)

	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("ファイルのオープンに失敗しました: %v", err)
	}
	defer file.Close()

	// ファイルをアップロード
	wc := object.NewWriter(context.Background())

	// メタデータを設定
	wc.ObjectAttrs = storage.ObjectAttrs{
		Name:        objectPath,
		ContentType: "image/jpeg",
		ACL:         []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}},
	}

	_, err = io.Copy(wc, file)
	if err != nil {
		return "", fmt.Errorf("ファイルのアップロードに失敗しました: %v", err)
	}

	if err = wc.Close(); err != nil {
		return "", fmt.Errorf("ライターのクローズに失敗しました: %v", err)
	}

	// 公開URLを取得
	attrs, err := object.Attrs(context.Background())
	if err != nil {
		return "", fmt.Errorf("オブジェクトの属性取得に失敗しました: %v", err)
	}

	return attrs.MediaLink, nil
}

// デフォルトアイコンのURLを取得
func GetDefaultIconURL(objectPath string) (string, error) {
	// 公開URLを生成
	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", config.Config.StorageBucket, url.PathEscape(objectPath))
	return url, nil
}
