package handler

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/infrastructure/repository"
	"security_chat_app/internal/interface/markup"
	"security_chat_app/internal/interface/middleware"
	"security_chat_app/internal/utils/icons"
)

// プロフィールページのデータ構造体
type ProfileData struct {
	IsLoggedIn     bool
	LoggedInUserID string
	User           *domain.User
}

// プロフィールページの表示
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.ValidateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// URLからユーザーIDを取得
	path := r.URL.Path
	var targetUserID string
	if path == "/profile" || path == "/profile/" {
		targetUserID = session.User.ID
	} else {
		targetUserID = path[len("/profile/"):]
		if targetUserID == "" {
			log.Fatalf("ユーザーIDが指定されていません")
			return
		}
	}

	// ユーザー情報の取得
	user, err := repository.GetUserByID(targetUserID)
	if err != nil {
		log.Fatalf("ユーザー情報の取得に失敗: %v", err)
		return
	}

	// アイコンが設定されていない場合はデフォルトアイコンを設定
	if user.Icon == "" {
		randomNum := rand.Intn(7)
		defaultIconPath := fmt.Sprintf(icons.DefaultIconPath+"/default_icon_%s.png", icons.DefaultIconNames[randomNum])
		iconURL, er := firebase.GetDefaultIconURL(defaultIconPath)
		if er != nil {
			log.Fatalf("デフォルトアイコンの取得に失敗: %v", er)
			return
		}

		// ユーザーのIconURLを更新
		user.Icon = iconURL
		err = firebase.UpdateField("users", user.ID, "Icon", iconURL)
		if err != nil {
			log.Fatalf("アイコンURLの更新に失敗: %v", err)
			return
		}
	}

	// 最終更新日時を現在時刻に更新 (自分のプロフィールの場合のみ更新すべきか検討)
	if targetUserID == session.User.ID {
		user.UpdatedAt = time.Now()
		err = firebase.UpdateField("users", user.ID, "UpdatedAt", user.UpdatedAt)
		if err != nil {
			log.Fatalf("最終更新日時の更新に失敗: %v", err)
			return
		}
	}

	// プロフィールデータの作成
	data := ProfileData{
		IsLoggedIn:     true,
		LoggedInUserID: session.User.ID,
		User:           user,
	}

	// テンプレートを描画
	markup.GenerateHTML(w, data, "layout", "header", "profile", "footer")
}

// アイコンアップロードハンドラ
func ProfileIconHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.ValidateSession(w, r)
	if err != nil {
		log.Fatalf("セッションが無効: %v", err)
		return
	}

	// URLからユーザーIDを取得
	path := r.URL.Path
	prefix := "/profile/icon/"
	var targetUserID string
	if strings.HasPrefix(path, prefix) && len(path) > len(prefix) {
		targetUserID = path[len(prefix):]
	}

	// 自分のプロフィール以外での変更を防止
	if targetUserID != "" && targetUserID != session.User.ID {
		log.Fatalf("他のユーザーのアイコンは変更できません")
		return
	}

	// マルチパートフォームの解析
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Fatalf("フォームの解析に失敗: %v", err)
		return
	}

	// アイコンファイルを取得
	file, header, err := r.FormFile("icon")
	if err != nil {
		log.Fatalf("アイコンファイルの取得に失敗: %v", err)
		return
	}
	defer file.Close()

	// ファイルサイズの制限（5MB）
	const maxFileSize = 5 * 1024 * 1024
	if header.Size > maxFileSize {
		log.Fatalf("ファイルサイズは5MB以下にしてください: %v", err)
		return
	}

	// ファイルの拡張子を取得と検証
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedExts[ext] {
		http.Redirect(w, r, "/profile?error=アップロードできるファイル形式は.jpg、.jpeg、.pngのみです", http.StatusSeeOther)
		return
	}

	// 画像ファイルの検証
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Fatalf("ファイルの読み込みに失敗: %v", err)
		http.Redirect(w, r, "/profile?error=ファイルの読み込みに失敗しました", http.StatusSeeOther)
		return
	}
	filetype := http.DetectContentType(buff)
	if !strings.HasPrefix(filetype, "image/") {
		log.Fatalf("画像ファイルのみアップロード可能: %v", err)
		http.Redirect(w, r, "/profile?error=画像ファイルのみアップロード可能です", http.StatusSeeOther)
		return
	}
	file.Seek(0, 0)

	// 一時ファイルを作成
	tempFile, err := os.CreateTemp("", "icon-*"+ext)
	if err != nil {
		log.Fatalf("一時ファイルの作成に失敗: %v", err)
		return
	}
	defer tempFile.Close()

	// ファイルをコピー
	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Fatalf("ファイルの保存に失敗: %v", err)
		return
	}

	// 一時ファイルのパスを取得
	tempFilePath := tempFile.Name()

	// Firebase Storageにアップロード
	iconURL, err := firebase.UploadIcon(session.User.ID, tempFilePath)
	if err != nil {
		log.Fatalf("アイコンのアップロードに失敗: %v", err)
		http.Redirect(w, r, "/profile?error=アイコンのアップロードに失敗しました", http.StatusSeeOther)
		return
	}

	// 一時ファイルを削除
	os.Remove(tempFilePath)

	// ユーザードキュメントを更新
	err = firebase.UpdateField("users", session.User.ID, "Icon", iconURL)
	if err != nil {
		log.Fatalf("ユーザー情報の更新に失敗: %v", err)
		http.Redirect(w, r, "/profile?error=ユーザー情報の更新に失敗しました", http.StatusSeeOther)
		return
	}

	// プロフィールページにリダイレクト
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
