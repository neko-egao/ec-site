package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// テスト実行前に呼ばれる特別な関数
func TestMain(m *testing.M) {
	// -------------------
	// テストの準備 (Setup)
	// -------------------
	// .envファイルを読み込む
	if err := godotenv.Load(); err != nil {
		log.Println("警告: .envファイルが読み込めませんでした")
	}

	// ★ 環境変数からテスト用のDB接続情報を取得するように変更
	dbUser := os.Getenv("TEST_DB_USER")
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	dbName := os.Getenv("TEST_DB_NAME")

	// 接続先を"db"（docker-compose.ymlで定義したサービス名）に変更
	// parseTime=true の前に &charset=utf8mb4 を追加
	dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/%s?charset=utf8mb4&parseTime=true&loc=Local&collation=utf8mb4_unicode_ci", dbUser, dbPassword, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("テストDBへの接続に失敗しました: %v", err)
	}

	// バリデータを初期化
	validate = validator.New()

	// -------------------
	// 全てのテストを実行
	// -------------------
	code := m.Run()

	// -------------------
	// テストの後片付け (Teardown)
	// -------------------
	// (今回は特に不要)

	os.Exit(code)
}

func TestGetProducts(t *testing.T) {
	// テストの前に、テーブルを空にしてテストデータを投入
	db.Exec("DELETE FROM products")
	db.Exec("INSERT INTO products (id, name, price) VALUES (1, 'テストTシャツ', 1000)")

	// 1. テスト用のリクエストとレスポンスレコーダーを作成
	req, _ := http.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()

	// 2. ルーターとハンドラを準備
	// (main関数から持ってくる)
	r := mux.NewRouter()
	r.HandleFunc("/products", productsHandler)

	// 3. リクエストを実行
	r.ServeHTTP(rr, req)

	// 4. 結果を検証
	if rr.Code != http.StatusOK {
		t.Errorf("期待と違うステータスコード: got %d want %d", rr.Code, http.StatusOK)
	}

	var products []Product
	json.Unmarshal(rr.Body.Bytes(), &products)
	if len(products) != 1 {
		t.Errorf("期待と違う商品数: got %d want %d", len(products), 1)
	}
	if products[0].Name != "テストTシャツ" {
		t.Errorf("期待と違う商品名: got %s want %s", products[0].Name, "テストTシャツ")
	}
}

// ★★★★★ この3つの関数を追記します ★★★★★

// POST /products のテスト
func TestCreateProduct(t *testing.T) {
	// 正常なリクエストボディ
	body := strings.NewReader(`{"name":"テストジュース","price":150}`)

	// 1. テスト用のリクエストとレスポンスレコーダーを作成
	req, _ := http.NewRequest("POST", "/products", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// 2. ルーターとハンドラを準備
	r := mux.NewRouter()
	r.HandleFunc("/products", productsHandler)

	// 3. リクエストを実行
	r.ServeHTTP(rr, req)

	// 4. 結果を検証
	if rr.Code != http.StatusCreated { // 201 Createdが返ってくることを期待
		t.Errorf("期待と違うステータスコード: got %d want %d", rr.Code, http.StatusCreated)
	}

	// バリデーションエラーのテスト
	body = strings.NewReader(`{"name":"短","price":-100}`) // 不正なデータ
	req, _ = http.NewRequest("POST", "/products", body)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest { // 400 Bad Requestが返ってくることを期待
		t.Errorf("期待と違うステータスコード: got %d want %d", rr.Code, http.StatusBadRequest)
	}
}

// PUT /products/{id} のテスト
func TestUpdateProduct(t *testing.T) {
	// 事前にテストデータを1件登録
	db.Exec("DELETE FROM products")
	db.Exec("INSERT INTO products (id, name, price) VALUES (99, '更新前', 10)")

	body := strings.NewReader(`{"name":"更新後テスト","price":123}`)
	req, _ := http.NewRequest("PUT", "/products/99", body) // ID:99を更新
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/products/{id}", productDetailHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("期待と違うステータスコード: got %d want %d", rr.Code, http.StatusOK)
	}
}

// DELETE /products/{id} のテスト
func TestDeleteProduct(t *testing.T) {
	// 事前にテストデータを1件登録
	db.Exec("DELETE FROM products")
	db.Exec("INSERT INTO products (id, name, price) VALUES (101, '削除対象', 10)")

	req, _ := http.NewRequest("DELETE", "/products/101", nil) // ID:101を削除
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/products/{id}", productDetailHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("期待と違うステータスコード: got %d want %d", rr.Code, http.StatusOK)
	}
}

// ★★★★★ この2つの関数を追記します ★★★★★

// POST /login のテスト
func TestLogin(t *testing.T) {
	// 事前にテスト用のユーザーを準備
	db.Exec("DELETE FROM users")
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)", "テストユーザー", "test@example.com", string(hashedPassword), "admin")

	// 正常なログインリクエスト
	body := strings.NewReader(`{"email":"test@example.com","password":"password123"}`)
	req, _ := http.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler)
	r.ServeHTTP(rr, req)

	// ステータスコードが200 OKであることを確認
	if rr.Code != http.StatusOK {
		t.Errorf("期待と違うステータスコード: got %d want %d", rr.Code, http.StatusOK)
	}

	// レスポンスにtokenが含まれていることを確認
	var resBody map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resBody)
	if _, ok := resBody["token"]; !ok {
		t.Errorf("レスポンスにトークンが含まれていません")
	}
}

// adminOnlyMiddlewareのテスト
func TestAdminOnlyMiddleware(t *testing.T) {
	// テスト用のダミーハンドラ
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// --- シナリオ1：トークンなしでアクセス ---
	req, _ := http.NewRequest("POST", "/protected", nil)
	rr := httptest.NewRecorder()
	// ミドルウェアを適用したハンドラを作成
	protectedHandler := adminOnlyMiddleware(dummyHandler)
	protectedHandler.ServeHTTP(rr, req)

	// 401 Unauthorizedが返ってくることを期待
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("トークンなしの場合、期待と違うステータスコード: got %d want %d", rr.Code, http.StatusUnauthorized)
	}

	// --- シナリオ2：有効な管理者トークンでアクセス ---
	// まずログインしてトークンを取得 (この部分はTestLoginから簡略化)
	claims := jwt.MapClaims{"user_id": 1, "role": "admin", "exp": time.Now().Add(time.Hour * 1).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	req, _ = http.NewRequest("POST", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr = httptest.NewRecorder()
	protectedHandler.ServeHTTP(rr, req)

	// 200 OKが返ってくることを期待
	if rr.Code != http.StatusOK {
		t.Errorf("有効なトークンの場合、期待と違うステータスコード: got %d want %d", rr.Code, http.StatusOK)
	}
}
