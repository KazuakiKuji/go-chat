package firebase

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"security_chat_app/internal/config"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func InitFirebase() (*firestore.Client, error) {
	opt := option.WithCredentialsFile(config.Config.ServiceKeyPath)

	firebaseConfig := &firebase.Config{
		ProjectID:     config.Config.ProjectId,
		StorageBucket: config.Config.StorageBucket,
	}

	app, err := firebase.NewApp(context.Background(), firebaseConfig, opt)
	if err != nil {
		log.Printf("Firebaseアプリの初期化に失敗: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("Firestoreクライアント作成に失敗: %v", err)
		return nil, err
	}

	// デフォルトアイコンの初期化
	if err := initDefaultIcons(app); err != nil {
		log.Printf("デフォルトアイコンの初期化に失敗: %v", err)
	}

	return client, nil
}

// デフォルトアイコンを初期化する
func initDefaultIcons(app *firebase.App) error {
	ctx := context.Background()
	storageClient, err := app.Storage(ctx)
	if err != nil {
		return fmt.Errorf("Storageクライアントの作成に失敗: %v", err)
	}
	bucket, err := storageClient.DefaultBucket()
	if err != nil {
		return fmt.Errorf("デフォルトバケットの取得に失敗: %v", err)
	}
	
	// デフォルトアイコンディレクトリの存在確認
	storagePrefix := "icons/default/"
	it := bucket.Objects(ctx, &storage.Query{Prefix: storagePrefix})
	
	hasObjects := false
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("オブジェクトの列挙に失敗: %v", err)
		}
		hasObjects = true
		break
	}
	
	// デフォルトアイコンが存在しない場合、アップロードする
	if !hasObjects {
		localIconDir := config.Config.DefaultIconDir
		if localIconDir == "" {
			return fmt.Errorf("デフォルトアイコンディレクトリのパスが設定されていません")
		}

		// デフォルトアイコンディレクトリが存在しない場合はエラー
		if _, err := os.Stat(localIconDir); os.IsNotExist(err) {
			log.Printf("デフォルトアイコンディレクトリが存在しません: %s", localIconDir)
			return fmt.Errorf("デフォルトアイコンディレクトリが存在しません: %s", localIconDir)
		}

		files, err := os.ReadDir(localIconDir)
		if err != nil {
			return fmt.Errorf("デフォルトアイコンディレクトリの読み込みに失敗: %v", err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			
			filePath := filepath.Join(localIconDir, file.Name())
			fileContent, err := os.Open(filePath)
			if err != nil {
				log.Printf("ファイル %s のオープンに失敗: %v", filePath, err)
				continue
			}
			defer fileContent.Close()

			objectPath := storagePrefix + file.Name()
			obj := bucket.Object(objectPath)
			writer := obj.NewWriter(ctx)
			writer.ObjectAttrs = storage.ObjectAttrs{
				Name:        objectPath,
				ContentType: "image/png",
				ACL:         []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}},
			}

			if _, err := io.Copy(writer, fileContent); err != nil {
				log.Printf("ファイル %s のアップロードに失敗: %v", filePath, err)
				writer.Close()
				continue
			}

			if err := writer.Close(); err != nil {
				log.Printf("ファイル %s のアップロード完了に失敗: %v", filePath, err)
				continue
			}			
		}
	}	
	return nil
}

