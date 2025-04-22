package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/interface/markup"
	"security_chat_app/internal/interface/middleware"
	"security_chat_app/internal/utils/uuid"
)

// 新規登録画面の表示と確認画面への遷移を処理
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := domain.TemplateData{
			IsLoggedIn: false,
		}
		markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
		return
	}

	// サインアップ処理
	if r.Method == http.MethodPost {
		r.ParseForm()
		form := domain.SignupForm{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		// バリデーション
		validationErrors := validateSignupForm(form)
		if len(validationErrors) > 0 {
			log.Printf("バリデーションエラー: %v", validationErrors)
			renderSignupError(w, form, validationErrors)
			return
		}

		// メールアドレスの重複チェック
		existingUsers, err := checkEmailDuplicate(form.Email)
		if err != nil {
			log.Printf("ユーザー検索エラー: %v", err)
			validationErrors := []string{"エラーが発生しました"}
			renderSignupError(w, form, validationErrors)
			return
		}

		if existingUsers {
			log.Printf("メールアドレス重複エラー: %s", form.Email)
			validationErrors := []string{"このメールアドレスは既に登録されています"}
			renderSignupError(w, form, validationErrors)
			return
		}

		// ユーザーの作成と保存
		user, err := createAndSaveUser(form)
		if err != nil {
			log.Printf("ユーザー作成エラー: %v", err)
			validationErrors := []string{"ユーザー作成エラーが発生しました"}
			renderSignupError(w, form, validationErrors)
			return
		}

		// セッションの作成
		session, err := middleware.CreateSession(user)
		if err != nil {
			log.Printf("セッション作成エラー: %v", err)
			validationErrors := []string{"セッション作成エラーが発生しました"}
			renderSignupError(w, form, validationErrors)
			return
		}

		middleware.SetSessionCookie(w, session)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// その他のHTTPメソッドは許可しない
	log.Fatalf("メソッドが許可されていません")
}

// 登録内容の確認とFirebaseへの保存を処理
func SignupConfirmHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		var form domain.SignupForm
		if r.Method == http.MethodGet {
			form = domain.SignupForm{
				Name:     r.URL.Query().Get("name"),
				Email:    r.URL.Query().Get("email"),
				Password: r.URL.Query().Get("password"),
			}
		} else {
			r.ParseForm()
			form = domain.SignupForm{
				Name:     r.FormValue("name"),
				Email:    r.FormValue("email"),
				Password: r.FormValue("password"),
			}
		}

		// バリデーション
		validationErrors := validateSignupForm(form)
		if len(validationErrors) > 0 {
			renderSignupError(w, form, validationErrors)
			return
		}

		// メールアドレスの重複チェック
		existingUsers, err := checkEmailDuplicate(form.Email)
		if err != nil {
			log.Printf("ユーザー検索エラー: %v", err)
			validationErrors := []string{"エラーが発生しました"}
			renderSignupError(w, form, validationErrors)
			return
		}

		if existingUsers {
			validationErrors := []string{"このメールアドレスは既に登録されています"}
			renderSignupError(w, form, validationErrors)
			return
		}

		if r.Method == http.MethodPost {
			_, err := createAndSaveUser(form)
			if err != nil {
				log.Printf("ユーザー作成エラー: %v", err)
				validationErrors := []string{"ユーザー作成エラーが発生しました"}
				renderSignupError(w, form, validationErrors)
				return
			}

			// 登録成功後、ログインページにリダイレクト
			http.Redirect(w, r, "/login?success=true", http.StatusSeeOther)
			return
		}

		// 確認画面の表示
		data := domain.TemplateData{
			IsLoggedIn: false,
			SignupForm: form,
		}
		markup.GenerateHTML(w, data, "layout", "header", "register_confirm", "footer")
	default:
		log.Fatalf("メソッドが許可されていません")
	}
}

// サインアップフォームのバリデーション
func validateSignupForm(form domain.SignupForm) []string {
	var validationErrors []string
	if form.Name == "" {
		validationErrors = append(validationErrors, "名前を入力してください")
	}
	if form.Email == "" {
		validationErrors = append(validationErrors, "メールアドレスを入力してください")
	}
	if !strings.Contains(form.Email, "@") {
		validationErrors = append(validationErrors, "有効なメールアドレスを入力してください")
	}
	if form.Password == "" {
		validationErrors = append(validationErrors, "パスワードを入力してください")
	}
	if len(form.Password) < 8 {
		validationErrors = append(validationErrors, "パスワードは8文字以上で入力してください")
	}
	return validationErrors
}

// メールアドレスの重複チェック
func checkEmailDuplicate(email string) (bool, error) {
	existingUsers, err := firebase.GetDataByQuery("users", "Email", "==", email)
	if err != nil {
		log.Printf("ユーザー検索エラー: %v", err)
		return false, err
	}
	return len(existingUsers) > 0, nil
}

// ユーザーデータの作成と保存
func createAndSaveUser(form domain.SignupForm) (*domain.User, error) {
	// パスワードのハッシュ化
	hashedPassword, err := uuid.HashPassword(form.Password)
	if err != nil {
		log.Printf("パスワードハッシュ化エラー: %v", err)
		return nil, err
	}

	// UUIDの生成
	userID, err := uuid.GenerateUUID()
	if err != nil {
		log.Printf("UUID生成エラー: %v", err)
		return nil, err
	}

	// ユーザーの作成
	user := &domain.User{
		ID:        userID,
		Name:      form.Name,
		Email:     form.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsOnline:  false,
	}

	// Firestoreにユーザーを保存
	err = firebase.AddData("users", user, user.ID)
	if err != nil {
		log.Printf("ユーザー作成エラー: %v", err)
		return nil, err
	}

	return user, nil
}

// エラー時のテンプレート表示
func renderSignupError(w http.ResponseWriter, form domain.SignupForm, errors []string) {
	data := domain.TemplateData{
		IsLoggedIn:       false,
		SignupForm:       form,
		ValidationErrors: errors,
	}
	markup.GenerateHTML(w, data, "layout", "header", "register", "footer")
}
