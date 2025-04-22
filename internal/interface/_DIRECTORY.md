# interface ディレクトリ

クリーンアーキテクチャにおける「インターフェース層(Interface Layer)」にあたる階層です。

以下の役割を担います。

* 外部との接点（I/O）
  * Web(HTTP)サーバーのルーティング
  * CLI
  * DB接続
  * API通信

## インターフェース層のルール

1. 基本は外部との橋渡し
  * ビジネスロジックは記述しない
  * 処理の中心はusecase層に渡す
  * 結果を受け取って返すだけの処理を定義する
2. 「依存される側」ではなく「依存する側」
    ```:go
    // NG：usecase が interface に依存している
    type UserController struct {
        userRepo *UserRepository // NG
    }

    // OK：interface 側が usecase の定義に依存している
    type UserController struct {
        userUsecase usecase.UserUsecase
    }
    ```
3. 抽象化(インターフェース)を明示し、テストしやすくする

## interface と infrastructure の違い

interfaceの役割は「こういう窓口があるよ」と定義する

infrastructureが「実際はMySQLで実装するよ」と定義する
