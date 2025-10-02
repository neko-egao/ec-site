# --- ビルドステージ ---
# Goの公式イメージをベースにする
FROM golang:1.23-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係のファイルを先にコピーして、キャッシュを効かせる
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションをビルドする
# CGO_ENABLED=0 は、外部ライブラリに依存しない静的バイナリを生成するためのおままじない
# -o /app/main で、ビルドした実行ファイルを/app/mainという名前で出力
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# --- 実行ステージ ---
# alpineという非常に軽量なOSイメージをベースにする
FROM alpine:latest

# ビルドステージから、ビルド済みの実行ファイルだけをコピー
COPY --from=builder /app/main /app/main

# コンテナが起動した時に、このコマンドを実行する
CMD ["/app/main"]