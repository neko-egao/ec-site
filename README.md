# Go + Next.js ECサイト (ポートフォリオ)

## 概要

Go言語によるAPIと、Next.jsによるフロントエンドで構成された、ECサイトをイメージして作成したアプリケーションです。

## 主な機能

*   **商品管理 (CRUD)**
    *   商品一覧・詳細情報の取得 (`GET`)
    *   管理者権限による商品の登録・更新・削除 (`POST`, `PUT`, `DELETE`)
*   **ユーザー機能**
    *   ユーザー新規登録 (`POST /users`)
    *   入力値バリデーション
*   **認証・認可システム**
    *   ログイン機能 (`POST /login`)
    *   JWT (JSON Web Token) を用いた認証情報の管理
    *   ミドルウェアによる、管理者専用APIのアクセス制御 (ロールベース認可)
*   **開発環境**
    *   Docker / Docker Compose による、完全コンテナ化された開発環境
    *   Goの標準ライブラリによるAPI機能の自動テスト

## 技術スタック

### バックエンド
*   **言語**: Go
*   **データベース**: MySQL
*   **ルーティング**: `gorilla/mux`
*   **認証**: `golang-jwt/jwt`, `bcrypt`
*   **設定管理**: `.env`
*   **バリデーション**: `go-playground/validator`
*   **CORS**: `rs/cors`

### フロントエンド
*   **言語**: TypeScript
*   **フレームワーク**: Next.js (App Router)
*   **スタイリング**: Tailwind CSS
*   **認証**: `jwt-decode`

### インフラ・その他
*   **コンテナ化**: Docker, Docker Compose
*   **テスト**: Go `testing`パッケージ

## APIエンドポイント一覧

| HTTPメソッド | URL | 認証 | 説明 |
| :--- | :--- | :--- | :--- |
| `GET` | `/products` | 不要 | 全ての商品リストを取得します。 |
| `GET` | `/products/{id}` | 不要 | 指定したIDの商品情報を取得します。 |
| `POST` | `/products` | **要** (管理者) | 新しい商品を登録します。 |
| `PUT` | `/products/{id}` | **要** (管理者) | 指定したIDの商品情報を更新します。 |
| `DELETE` | `/products/{id}` | **要** (管理者) | 指定したIDの商品を削除します。 |
| `POST` | `/users` | 不要 | 新しいユーザーを登録します。 |
| `POST` | `/login` | 不要 | ログインし、JWTトークンを取得します。 |

## セットアップと実行手順

このプロジェクトはバックエンドがDocker化されているため、あなたのPCにGoやMySQLを直接インストールする必要はありません。**Docker**と**Git**がインストールされていれば、以下の手順で簡単に起動できます。

1.  **リポジトリをクローン**
    ```bash
    git clone https://github.com/neko-egao/ec-site.git
    cd ec-site
    ```

2.  **環境変数ファイルを作成**
    バックエンドの環境変数ファイルを作成します。


    **バックエンド:**
    ```bash
    cp ec-site-backend/.env.example ec-site-backend/.env
    # ec-site-backend/.env を必要に応じて編集（.env.exampleを参考に作成）
    ```

3.  **アプリケーションの起動**
    バックエンド（APIサーバーとデータベース）をDockerコンテナで起動し、フロントエンドをローカルで起動します。

    **バックエンドの起動:**
    ```bash
    cd ec-site-backend
    docker compose up --build -d
    ```
    **フロントエンドの起動:**
    ```bash
    cd ../ec-site-frontend
    npm install
    npm run dev
    ```

    起動後、以下のURLにアクセスできます。
    *   フロントエンド: `http://localhost:3000`
    *   バックエンドAPI: `http://localhost:8081`

4.  **テストの実行**
    バックエンドの自動テストは、以下のコマンドで実行できます。
    ```bash
    cd ec-site-backend
    docker compose exec api go test ./...
    ```

5.  **コンテナの停止**
    ```bash
    cd ec-site-backend
    docker compose down
    ```

## 注意点

*   初期時点でデフォルトの管理者ユーザー
   （メールアドレス：admin@example.com，パスワード：password123）
   が設定されています。
*   初期時点でいくつかの商品情報が登録されています。

## 今後の展望

*   商品画像の追加
*   商品検索機能の追加
*   カート機能の追加

