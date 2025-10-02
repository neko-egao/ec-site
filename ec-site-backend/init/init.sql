SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 開発用データベースがなければ作成（文字コードを明確に指定）
CREATE DATABASE IF NOT EXISTS ECsite DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- テスト用データベースがなければ作成（文字コードを明確に指定）
CREATE DATABASE IF NOT EXISTS ECsite_test DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 開発用DBにテーブルを作成
USE ECsite;

CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; -- テーブルの文字コードを指定

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; -- テーブルの文字コードを指定

-- サンプル商品を登録
INSERT INTO products (name, price) VALUES
('すごいTシャツ', 2500),
('すごいキャップ', 3800),
('すごいステッカーセット', 1200),
('超クールなマグカップ', 1800),
('もふもふパーカー', 7500),
('万能スマホケース', 3200),
('2トートバッグ', 4500);

-- サンプル管理者ユーザーを登録（password123をハッシュ化）
INSERT INTO users (name, email, password, role) VALUES
('Admin User', 'admin@example.com', '$2a$10$pC01EMCncl2AyRyTJBxqs.m9ZwiXlesfn3nRpHDxk5b5d5jIY5YNu', 'admin');


-- テスト用DBにも同じテーブルを作成
USE ECsite_test;
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; -- テーブルの文字コードを指定

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; -- テーブルの文字コードを指定