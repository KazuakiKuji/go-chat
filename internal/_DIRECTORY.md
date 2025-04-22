# internal ディレクトリについて

この `internal/` ディレクトリは、アプリケーションの内部実装を格納するためのディレクトリです。
クリーンアーキテクチャの原則に基づき、依存関係が内側に向かうように設計されています。

Go では、`internal/` 以下のパッケージは **このプロジェクト内からのみインポート可能** であり、他のプロジェクトからのアクセスを防ぐことができます。

## ディレクトリ構成と依存関係

依存の方向: infrastructure → interface → usecase → domain

```
internal/
├── config/          # アプリケーション設定
│   └── config.go
├── domain/          # エンティティ、ビジネスルール（最も内側のレイヤー）
│   ├── user.go
│   ├── chat.go
│   ├── post.go
│   ├── session.go
│   ├── form.go
│   └── template.go
├── usecase/         # ビジネスロジック、ユースケース
│   ├── user/
│   │   └── service.go
│   └── chat/
│       └── usecase.go
├── interface/       # 外部とのインターフェース、アダプター
│   ├── handler/
│   │   ├── chat_handler.go
│   │   ├── login_handler.go
│   │   ├── logout_hander.go
│   │   ├── profile_handler.go
│   │   ├── reset_password_handler.go
│   │   ├── search_handler.go
│   │   ├── settings_handler.go
│   │   └── signup_handler.go
│   ├── middleware/
│   │   ├── middleware.go
│   │   └── session.go
│   └── markup/
│       └── template.go
├── infrastructure/ # 外部技術の具体的な実装（最も外側のレイヤー）
│   ├── firebase/
│   │   ├── firestore.go
│   │   ├── setup.go
│   │   └── storage.go
│   ├── repository/
│   │   └── user_repository.go
│   └── router/
│       └── router.go
├── web/           # Web関連の静的ファイル
│   ├── static/
│   ├── templates/
│   ├── css/
│   ├── scss/
│   ├── js/
│   └── images/
└── utils/         # 汎用ユーティリティ
    └── log/
        └── logger.go
```

## レイヤー別の責務

### domain/

- ビジネスエンティティの定義
- ビジネスルールの実装
- インターフェースの定義
- 他のレイヤーへの依存を持たない
- 値オブジェクトやドメインサービスの定義

### usecase/

- アプリケーションのビジネスロジック
- ドメインオブジェクトの操作
- トランザクション管理
- domain レイヤーのみに依存
- ユースケース固有のデータ構造の定義

### interface/

- 外部とのやりとりを抽象化
- コントローラーの実装
- リポジトリインターフェースの実装
- usecase と domain レイヤーに依存
- データの変換（DTO とドメインオブジェクト間）

### infrastructure/

- 具体的な技術の実装
- データベースアクセス
- 外部 API クライアント
- フレームワーク統合
- interface レイヤーで定義されたインターフェースの実装

## クロスカッティングコンサーン

以下の機能は複数のレイヤーにまたがる可能性があります：

- **config/**: アプリケーション設定の管理

  - 環境変数の読み込み
  - 設定ファイルの管理
  - 外部サービスの認証情報

- **utils/**: 共通ユーティリティ
  - ロギング
  - エラーハンドリング
  - 日付処理
  - バリデーション

## 依存性注入

- 依存性は内側に向かって注入
- インターフェースを使用して実装を抽象化
- 具体的な実装は最も外側のレイヤーで提供

## テスト戦略

各レイヤーに応じたテスト方針：

- **domain/**: ユニットテスト

  - ビジネスルールの検証
  - 依存のないピュアな状態でテスト

- **usecase/**: ユニットテスト、統合テスト

  - モックを使用したテスト
  - ユースケースのフロー検証

- **interface/**: 統合テスト

  - 外部とのインターフェースの検証
  - データ変換の正確性確認

- **infrastructure/**: 結合テスト
  - 外部システムとの連携確認
  - 実際の環境に近い状態でのテスト

## 新機能追加のガイドライン

1. まずドメインモデルとビジネスルールを定義（domain/）
2. ユースケースとしてビジネスロジックを実装（usecase/）
3. 外部とのインターフェースを定義（interface/）
4. 具体的な技術実装を提供（infrastructure/）

## 注意事項

- 依存関係は必ず内側に向かうように保つ
- 各レイヤーの責務を明確に分離
- インターフェースを活用して疎結合を維持
- ビジネスロジックは必ず usecase レイヤーに実装
- 外部技術の詳細は infrastructure レイヤーに隠蔽

---

## 参考資料

- [The Clean Architecture（Uncle Bob）](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
