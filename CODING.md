# コーディング規約

## クリーンアーキテクチャの原則

このプロジェクトはクリーンアーキテクチャの原則に基づいて設計されています。クリーンアーキテクチャでは、依存関係が内側に向かうように設計し、各レイヤーの責務を明確に分離します。

### 依存関係の方向

依存関係は常に内側に向かいます：

- 外側のレイヤーは内側のレイヤーに依存する
- 内側のレイヤーは外側のレイヤーに依存しない
- 依存関係は抽象化（インターフェース）を通じて行う

### レイヤー構成

1. **Domain Layer** (最も内側)

   - ビジネスエンティティとルール
   - 外部に依存しない純粋なビジネスロジック

2. **Use Case Layer**

   - アプリケーション固有のビジネスルール
   - ドメインレイヤーのみに依存

3. **Interface Adapter Layer**

   - コントローラー、プレゼンター、ゲートウェイ
   - ユースケースレイヤーとドメインレイヤーに依存

4. **Frameworks & Drivers Layer** (最も外側)
   - データベース、Web、デバイス、フレームワーク
   - インターフェースアダプターレイヤーに依存

## コーディング規約

### 命名規則

- **パッケージ名**: 小文字の単数形（例: `user`, `chat`）
- **インターフェース名**: 動詞または名詞+er（例: `UserRepository`, `ChatService`）
- **構造体名**: 名詞（例: `User`, `Chat`）
- **メソッド名**: 動詞で始める（例: `GetUser`, `CreateChat`）
- **変数名**: キャメルケース（例: `userID`, `chatHistory`）
- **定数名**: 大文字のスネークケース（例: `MAX_RETRY_COUNT`）

### ファイル構成

- 各ファイルは単一の責任を持つ
- インターフェースと実装は別ファイルに分ける
- テストファイルは対象ファイルと同じディレクトリに配置し、`_test.go`サフィックスを付ける

### コメント

- 公開 API には必ずドキュメントコメントを付ける
- 複雑なロジックには説明コメントを付ける
- TODO コメントには担当者と期限を明記する

### エラーハンドリング

- エラーは適切なレベルで処理する
- カスタムエラー型を使用して意味のあるエラーメッセージを提供する
- エラーログには十分なコンテキスト情報を含める

### テスト

- ユニットテストは各レイヤーで実装する
- 依存関係はモックまたはスタブで置き換える
- テストカバレッジは 80%以上を目標とする

## ディレクトリ構成

```
project-root/
├── cmd/
│   └── app/
│       └── main.go          // アプリのエントリーポイント
├── internal/
│   ├── domain/              // エンティティやドメインモデル
│   │   ├── user.go
│   │   ├── chat.go
│   │   └── post.go
│
│   ├── usecase/             // ユースケース層（ビジネスロジック）
│   │   ├── user/
│   │   │   └── service.go
│   │   └── chat/
│   │       └── usecase.go
│
│   ├── interface/           // インターフェース層（外部との接続）
│   │   ├── handler/         // HTTPハンドラー
│   │   │   ├── chat_handler.go
│   │   │   ├── login_handler.go
│   │   │   └── user_handler.go
│   │   ├── middleware/      // ミドルウェア
│   │   │   └── auth.go
│   │   └── repository/      // リポジトリインターフェース
│   │       └── user_repository.go
│
│   ├── infrastructure/      // データベースや外部APIとのやり取り
│   │   ├── firebase/
│   │   │   ├── client.go    // Firebaseの初期化や設定
│   │   │   └── firestore.go // Firestore操作
│   │   └── router/
│   │       └── router.go    // ルーティング定義
│
│   ├── config/              // 設定の読み込み
│   │   └── config.go
│   │
│   └── web/                 // Web関連の静的ファイル
│       ├── templates/
│       ├── static/
│       ├── css/
│       ├── scss/
│       ├── js/
│       └── images/
│
├── pkg/                     // 公開パッケージ（必要であれば）
│   └── utils/
│       └── validator.go
│
├── test/                    // 単体・統合テスト
│   └── handler_test.go
│
├── go.mod
└── go.sum
```

## 参考資料

### クリーンアーキテクチャ関連

- [The Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Clean Architecture in Go](https://github.com/evrone/go-clean-template)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

### Go 言語関連

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go by Example](https://gobyexample.com/)

### テスト関連

- [Go Testing](https://golang.org/pkg/testing/)
- [Go Test Coverage](https://blog.golang.org/cover)
- [Go Testing Best Practices](https://github.com/crytal/go-testing-best-practices)
