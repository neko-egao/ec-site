package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	jatrans "github.com/go-playground/validator/v10/translations/ja"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

// Productの構造体
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"  validate:"required,min=3,max=50"`
	Price int    `json:"price" validate:"required,gt=0"`
}

// ユーザー登録リクエスト専用の構造体
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Userのデータ構造
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"-"` // パスワードはレスポンスに含めない
	CreatedAt string `json:"created_at"`
	Role      string `json:"role"`
}

// ログインリクエスト用の構造体
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	db       *sql.DB
)

// /products のリクエストを処理するハンドラ
func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)
	case http.MethodPost:
		createProduct(w, r)
	default:
		http.Error(w, "許可されていないメソッドです", http.StatusMethodNotAllowed)
	}
}

// 商品リストを取得する (GET /products)
func getProducts(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// 新しい商品を登録する (POST /products)
func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(p); err != nil {
		errs := err.(validator.ValidationErrors)
		var errMap []string
		for _, e := range errs {
			errMap = append(errMap, e.Translate(trans))
		}
		http.Error(w, strings.Join(errMap, ", "), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO products (name, price) VALUES (?, ?)"
	result, err := db.Exec(query, p.Name, p.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "商品が正常に登録されました",
		"id":      lastID,
	})
}

// /products/{id} のリクエストを処理するハンドラ
func productDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	switch r.Method {
	case http.MethodGet:
		getProduct(w, r, id)
	case http.MethodPut:
		updateProduct(w, r, id)
	case http.MethodDelete:
		deleteProduct(w, r, id)
	default:
		http.Error(w, "許可されていないメソッドです", http.StatusMethodNotAllowed)
	}
}

// 特定の商品を取得する (GET /products/{id})
func getProduct(w http.ResponseWriter, _ *http.Request, id string) {
	var p Product
	err := db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "商品が見つかりません", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// 特定の商品を更新する (PUT /products/{id})
func updateProduct(w http.ResponseWriter, r *http.Request, id string) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(p); err != nil {
		errs := err.(validator.ValidationErrors)
		var errMap []string
		for _, e := range errs {
			errMap = append(errMap, e.Translate(trans))
		}
		http.Error(w, strings.Join(errMap, ", "), http.StatusBadRequest)
		return
	}
	query := "UPDATE products SET name = ?, price = ? WHERE id = ?"
	result, err := db.Exec(query, p.Name, p.Price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 更新された行がなかった場合（該当IDの商品が存在しなかった場合）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "更新対象の商品が見つかりません", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "商品が正常に更新されました"})
}

// 特定の商品を削除する (DELETE /products/{id})
func deleteProduct(w http.ResponseWriter, _ *http.Request, id string) {
	query := "DELETE FROM products WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 削除された行がなかった場合（該当IDの商品が存在しなかった場合）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "削除対象の商品が見つかりません", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "商品が正常に削除されました"})
}

// /users のリクエストを処理するハンドラ
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "許可されていないメソッドです", http.StatusMethodNotAllowed)
	}
}

// 新規ユーザーを登録する (POST /users)
func createUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		var errMap []string
		for _, e := range errs {
			errMap = append(errMap, e.Translate(trans))
		}
		http.Error(w, strings.Join(errMap, ", "), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	result, err := db.Exec(query, req.Name, req.Email, string(hashedPassword))
	if err != nil {
		// メールアドレスの重複エラーをハンドリング
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, "このメールアドレスは既に使用されています", http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "ユーザーが正常に登録されました",
		"id":      lastID,
	})
}

// ログイン処理を行う (POST /login)
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	query := "SELECT id, name, email, password, role FROM users WHERE email = ?"
	err := db.QueryRow(query, req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "ユーザーが見つからないか、パスワードが間違っています", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "ユーザーが見つからないか、パスワードが間違っています", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user":  user,
	})
}

// 管理者権限をチェックするミドルウェア
func adminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "リクエストヘッダーにトークンが含まれていません", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("予期しない署名方法です: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if role, ok := claims["role"].(string); ok && role == "admin" {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		http.Error(w, "アクセス権限がありません", http.StatusForbidden)
	})
}

func main() {
	// 日本語のロケールと翻訳機をセットアップ
	jaLocale := ja.New()
	uni = ut.New(jaLocale, jaLocale)
	trans, _ = uni.GetTranslator("ja")

	// バリデータを初期化し、日本語翻訳を登録
	validate = validator.New()
	jatrans.RegisterDefaultTranslations(validate, trans)

	if err := godotenv.Load(); err != nil {
		log.Println("警告: .envファイルが読み込めませんでした")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("環境変数 DB_USER, DB_PASSWORD, DB_NAME が設定されていません")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/%s?charset=utf8mb4&parseTime=true", dbUser, dbPassword, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("データベースへの接続に失敗しました: %v", err)
	}
	defer db.Close()

	// --- ルーティング設定 ---
	r := mux.NewRouter()

	// 認証が不要な公開ルート
	r.HandleFunc("/users", usersHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")

	// --- /products へのリクエストを処理するサブルーター ---
	productsRouter := r.PathPrefix("/products").Subrouter()

	// GET /products (認証不要)
	productsRouter.HandleFunc("", productsHandler).Methods("GET")
	// POST /products (要管理者権限)
	productsRouter.Handle("", adminOnlyMiddleware(http.HandlerFunc(productsHandler))).Methods("POST")

	// --- /products/{id} へのリクエストを処理するサブルーター ---
	productDetailRouter := r.PathPrefix("/products/{id:[0-9]+}").Subrouter()

	// GET /products/{id} (認証不要)
	productDetailRouter.HandleFunc("", productDetailHandler).Methods("GET")
	// PUT, DELETE /products/{id} (要管理者権限)
	productDetailRouter.Handle("", adminOnlyMiddleware(http.HandlerFunc(productDetailHandler))).Methods("PUT", "DELETE")

	// --- CORSの設定と最終的なハンドラの作成 ---
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(r)

	log.Println("サーバー起動中... http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
