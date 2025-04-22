package domain

// ログインフォームのデータ構造体
type LoginForm struct {
	Email    string // メールアドレス
	Password string // パスワード
}

// サインアップフォームのデータ構造体
type SignupForm struct {
	Name     string // 名前
	Email    string // メールアドレス
	Password string // パスワード
}

// パスワードリセットフォームのデータ構造体
type ResetForm struct {
	Email           string // メールアドレス
	Password        string // パスワード
	PasswordConfirm string // パスワード確認
}

// パスワード変更フォームのデータ構造体
type PasswordForm struct {
	CurrentPassword    string // 現在のパスワード
	NewPassword        string // 新しいパスワード
	NewPasswordConfirm string // 新しいパスワード確認
}
