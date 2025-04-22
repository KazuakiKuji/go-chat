# go-chat-app

Goを使用した基本的なチャットアプリになります。

## アプリURL
https://go-chat-app.kuji-server.com/

### ■ Basic認証
ユーザー名：ey9n5cjiq5ye

パスワード：ma(5k9bv.dtn

## スクリーンショット(2025.04.19時点)
<img src="https://github.com/user-attachments/assets/2db48f5c-f321-47e0-8f1b-c85e8f208684" width="400">
<img src="https://github.com/user-attachments/assets/6bbdf874-0680-4566-b7da-768c4dbc60a8" width="400">
<img src="https://github.com/user-attachments/assets/8a7b27b7-41b7-4093-9f73-47089b381892" width="400">
<img src="https://github.com/user-attachments/assets/37d12558-07b3-4cdb-9380-76e7ca8d1981" width="400">

## 実装済み機能
- 認証機能（登録/ログイン/ログアウト）
- プロフィール（ユーザー名・画像・パスワードなどの変更）
- 検索機能（登録済みユーザーのフィルタリング）
- チャット機能（他ユーザーと連絡）

## 使用技術
- Go
- Firebase(firestore, storage)
- HTML/CSS(SCSS)/JavaScript
- stylelint, prettier, gulp

## ローカルへの導入手順

1. **プロジェクトのclone**
   
    ```bash
    git clone git@github.com:Kazu-K0032/go-chat-app.git
    ```

2. **FirebaseからFirebase Admin SDKの認証ファイルを取り込む**

    * [Firebase](https://console.firebase.google.com/u/1/?hl=ja)からプロジェクトを作成
    * 作成したプロジェクトにアクセスし、「プロジェクトの設定」⇒「サービスアカウント」⇒「新しい鍵を生成」
      
      <img src="https://github.com/user-attachments/assets/c0820422-87d5-4490-80aa-cfe02c564456" width="400">
      <img src="https://github.com/user-attachments/assets/de34f37d-d44b-40a4-8e6f-44ec215f11c9" height="300">
    * ダウンロードしたファイル名を「serviceAccountKey.json」に変更し、クローンしたプロジェクトの`internal/config/`に配置してください

3. **Firestoreの設定**

    * 左サイドバー「構築」⇒「Firestore Database」から、「データベースを作成」
      
      <img src="https://github.com/user-attachments/assets/85cb2709-e414-4e69-84d1-abbeda4f10f7" height="300">
      <img src="https://github.com/user-attachments/assets/8bcdd85f-75ac-4771-88ca-8d73fa04e35a" height="300">
  
    * 「Cloud Firestore」⇒「ルール」タブから、以下のルールであることを確認
      
      ```js
      rules_version = '2';
      
      service cloud.firestore {
        match /databases/{database}/documents {
          match /{document=**} {
            allow read, write: if false;
          }
        }
      }
      ```

4. **Storageの設定**

      * 左サイドバー「構築」⇒「Storage」を選択
         * Storageを始める場合、請求先設定が必要になります。
        
         <img src="https://github.com/user-attachments/assets/6e107ac9-c117-45a3-8307-3c3b494b9b57" height="300">
         <img src="https://github.com/user-attachments/assets/fde43faf-2f26-4693-bf28-5c8ca77ca917" height="300">

      * 「Storage」⇒「ルール」タブから、以下のルールに変更

         ```js
         rules_version = '2';
         service firebase.storage {
           match /b/{bucket}/o {
             match /icons/default/{fileName} {
               allow read: if true;
               allow write: if false;
             }
             match /icons/{userId}/{fileName} {
               allow read: if true;
               allow write: if request.auth != null && request.auth.uid == userId;
             }
             match /{allPaths=**} {
               allow read, write: if request.auth != null;
             }
           }
         }
         ```


6. **設定ファイルの修正（`config.ini`）**

    * または `config.ini` をコピーし `config.local.ini`に変更してください。

    * ポートの設定をしています。ご自身の環境に合わせて、随時修正してください。

    ```txt
    [web]
    port = 8050
    logfile = debug.log
    static = app/views
    
    [firebase]
    defaultIconDir = icons/default/
    serviceKeyPath = internal/config/serviceAccountKey.json // serviceAccountKey.jsonの相対パス
    projectId = // <プロジェクトの設定> -> <全般> -> <プロジェクトID> の値
    storageBucket = // <Storage> -> <バケット ex: testa87e4.firebasestorage.app>
    ```

    ### 参考(projectId)
    <img src="https://github.com/user-attachments/assets/ee00624a-0634-4f30-8ccf-65b1eedb23d7" height="300">
    
    ### 参考(storageBucket)
    <img src="https://github.com/user-attachments/assets/7918ecc4-1617-478d-82db-65b183fbbf33" height="300">


8. **モジュール初期化および依存解決**

   * 事前に、Go及びNode.jsをダウンロードしてください。
      * バージョンは、Goは最低1.21以上, Node.jsはv16.0.0以上を目安に更新してください。
      
      ```bash
      go version
      node -v
      ```

   * 以下を実行してください。
      ```:bash
      cd go-chat-app/
      
      # Go モジュールの初期化
      go mod tidy
      
      # Node.jsの依存解決
      npm install
      ```

9. **サーバーの起動**

     ```bash
     go run cmd/app/main.go
     ```
     * 実行後、`debug.log`が生成されます。
     * デフォルトだと、`localhost:8050`にアクセスできるようになります。
