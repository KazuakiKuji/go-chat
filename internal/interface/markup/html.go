package markup

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"security_chat_app/internal/domain"
	"security_chat_app/internal/utils/icons"
	"security_chat_app/internal/utils/random"
)

// テンプレートで使用する関数
var templateFuncs = template.FuncMap{
	"sub": func(a, b int) int {
		return a - b
	},
	"len": func(slice any) int {
		switch v := slice.(type) {
		case []domain.Message:
			return len(v)
		case []domain.Contact:
			return len(v)
		case []domain.Chat:
			return len(v)
		default:
			return 0
		}
	},
	"substr": func(s string, start, length int) string {
		if start < 0 {
			start = 0
		}
		if length < 0 {
			length = len(s)
		}
		if start > len(s) {
			return ""
		}
		end := start + length
		if end > len(s) {
			end = len(s)
		}
		return s[start:end]
	},
	"getRandomDefaultIcon": func() string {
		// 0から6までのランダムな数字を生成
		randomNum := random.LocalRand.Intn(icons.DefaultIconCount)
		// デフォルトアイコンのパスを生成
		return fmt.Sprintf("%s/%s.png", icons.DefaultIconPath, icons.DefaultIconNames[randomNum])
	},
}

// GenerateHTML layout.htmlをベースとしたHTMLを生成し、レスポンスに書きだす
func GenerateHTML(writer http.ResponseWriter, data any, filenames ...string) {
	var files []string
	for _, file := range filenames {
		path := fmt.Sprintf("internal/web/templates/%s.html", file)
		files = append(files, path)
	}

	templates, err := template.New("layout").Funcs(templateFuncs).ParseFiles(files...)
	if err != nil {
		log.Fatalf("テンプレートの読み込みに失敗: %v", err)
		return
	}

	// テンプレートをバッファに出力
	var buf bytes.Buffer
	err = templates.ExecuteTemplate(&buf, "layout", data)
	if err != nil {
		log.Fatalf("テンプレートの実行に失敗: %v", err)
		return
	}

	// 成功したらまとめて出力
	buf.WriteTo(writer)
}
