package handler

import (
	"log"
	"net/http"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/infrastructure/firebase"
	"security_chat_app/internal/interface/markup"
	utils "security_chat_app/internal/utils/uuid"
)

// パスワード再設定処理を実行
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := domain.TemplateData{
			IsLoggedIn: false,
			ResetForm:  domain.ResetForm{},
		}
		markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		form := domain.ResetForm{
			Email: r.FormValue("email"),
		}
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("password_confirm")

		var validationErrors []string
		if form.Email == "" {
			validationErrors = append(validationErrors, "メールアドレスを入力してください")
		}
		if password == "" {
			validationErrors = append(validationErrors, "新しいパスワードを入力してください")
		}
		if len(password) < 8 {
			validationErrors = append(validationErrors, "パスワードは8文字以上で入力してください")
		}
		if password != passwordConfirm {
			validationErrors = append(validationErrors, "パスワードが一致しません")
		}

		if len(validationErrors) > 0 {
			log.Printf("バリデーションエラー: %v", validationErrors)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				ResetForm:        form,
				ValidationErrors: validationErrors,
			}
			markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
			return
		}

		// ユーザー検索
		users, err := firebase.GetDataByQuery("users", "Email", "==", form.Email)
		if err != nil {
			log.Printf("ユーザー検索エラー: %v", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				ResetForm:        form,
				ValidationErrors: []string{"ユーザー検索エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
			return
		}

		if len(users) == 0 {
			log.Printf("ユーザーが見つかりません: %s", form.Email)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				ResetForm:        form,
				ValidationErrors: []string{"該当するユーザーが見つかりません"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
			return
		}

		// パスワードのハッシュ化
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			log.Printf("パスワードハッシュ化エラー: %v", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				ResetForm:        form,
				ValidationErrors: []string{"パスワード更新エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
			return
		}

		// パスワード更新
		userID := users[0]["ID"].(string)
		err = firebase.UpdateField("users", userID, "Password", hashedPassword)
		if err != nil {
			log.Printf("パスワード更新エラー: %v", err)
			data := domain.TemplateData{
				IsLoggedIn:       false,
				ResetForm:        form,
				ValidationErrors: []string{"パスワード更新エラーが発生しました"},
			}
			markup.GenerateHTML(w, data, "layout", "header", "reset-password", "footer")
			return
		}

		// 成功時はログインページにリダイレクト
		http.Redirect(w, r, "/login?success=パスワードを再設定しました", http.StatusSeeOther)
		return
	}

	// その他のHTTPメソッドは許可しない
	log.Fatalf("メソッドが許可されていません")
}
