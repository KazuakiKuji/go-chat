package domain

import "time"

// ユーザーの構造体
type User struct {
	ID            string    // ユーザーのID
	Name          string    // ユーザーの名前
	Email         string    // ユーザーのメールアドレス
	Password      string    // ユーザーのパスワード
	CreatedAt     time.Time // ユーザーの作成日時
	UpdatedAt     time.Time // ユーザーの更新日時
	IsOnline      bool      // ユーザーがオンラインかどうか
	Icon          string    // ユーザーのアイコン
	Contacts      []Contact // ユーザーの連絡先
}

// 連絡先を交換したユーザーの構造体
type Contact struct {
	ID       string    // 連絡先のID
	Username string    // 連絡先のユーザー名
	Icon     string    // 連絡先のアイコンのURL
	LastSeen time.Time // 連絡先の最終接続日時
	IsOnline bool      // 連絡先がオンラインかどうか
}
