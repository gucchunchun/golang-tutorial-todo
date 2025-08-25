# golang-tutorial-todo

## 概要 

**golang-tutorial-todo** は Go 言語で作成されたシンプルな **Todo アプリケーション**です。
このプロジェクトは、Go のクリーンアーキテクチャ、モジュール設計、CLI/REST API/Web UI の実装を学習するためのチュートリアルとして設計されています。

### 特徴

* **CLI**: ターミナルからタスクを追加・一覧・更新できる
* **REST API**: タスク管理用のエンドポイントを提供
* **ストレージ**:
  * MySQL (デフォルト)
  * JSON ファイル（ローカル用・学習用）
* **その他の機能**:
  * 構造化ロギング
  * 入力バリデーション
  * 名言 API との連携（ランダムな名言表示）

### プロジェクト構成

```
golang-tutorial-todo/
├── cmd/                # エントリーポイント (CLI & Server)
│   ├── todo/           # CLI アプリ
│   └── todo-server/    # API サーバ
├── internal/
│   ├── api/            # ルーター & ミドルウェア
│   ├── app/            # 初期化処理、サーバ、CLI ハンドラ
│   ├── db/             # データベース設定 & マイグレーション
│   ├── models/         # ドメインモデル
│   ├── storage/        # ストレージ層 (MySQL, JSON)
│   ├── quote/          # 名言クライアント
│   └── utils/          # ユーティリティ関数
├── web/                # Web テンプレート & 静的ファイル
├── .env.example        # 環境変数サンプル
├── Makefile            # ビルド & 実行用タスク
├── tasks.json          # JSON ストレージ用サンプルタスク
└── go.mod              # Go モジュール
```

### 使用方法

#### cli
- タスク一覧取得

  コマンド：
  ```
  go run golang/tutorial/todo/cmd/todo list
  ```
  結果：
  ```
    ID                                   | Name                 | Status   | CreatedAt           | DueDate     | TimeLeft
  -----------------------------------------------------------------------------------------
  11               | Read book            | Pending  | 2025-08-24 18:40:00 | 2025-12-31  | 128 days       
  10               | Pay bills            | Pending  | 2025-08-24 18:40:00 | -           | -              
  9                | Buy milk             | Pending  | 2025-08-24 18:40:00 | 2025-09-01  | 7 days         
  8                | Read book            | Pending  | 2025-08-24 17:54:16 | 2025-12-31  | 128 days       
  ```

- ファイル暗号化・複合化
  ```
  go run golang/tutorial/todo/cmd/todo encrypt -i {SrcFilePath} -o {DestFilePath} -p {PW}
  ```
  ```
  go run golang/tutorial/todo/cmd/todo decrypt -i {SrcFilePath} -o {DestFilePath} -p {PW}
  ```

#### エンドポイント

* `GET /tasks/{id}` - タスク取得
* `GET /tasks` - タスク一覧
* `POST /tasks` - タスク追加
* `GET /quote` - ランダム名言取得

---

## Overview

**golang-tutorial-todo** is a simple **Todo Application** written in Go.
It is designed as a tutorial project to learn Go's clean architecture, modular design, and implementation of CLI, REST API, and Web UI.

### Features

* **CLI**: Manage tasks (add/list/update) from terminal
* **REST API**: Provides endpoints for task management
* **Web UI**: Simple UI with templates and static assets
* **Storage Backends**:

  * MySQL (default)
  * JSON file (local/demo use)
* **Other Utilities**:

  * Structured logging
  * Input validation
  * Random quote API integration

### Project Structure

```
golang-tutorial-todo/
├── cmd/                # Entry points (CLI & Server)
│   ├── todo/           # CLI app
│   └── todo-server/    # API server
├── internal/
│   ├── api/            # HTTP router & middleware
│   ├── app/            # Bootstrap, server, CLI handlers
│   ├── db/             # Database setup & migrations
│   ├── models/         # Domain models
│   ├── storage/        # Storage adapters (MySQL, JSON)
│   ├── quote/          # Quote client
│   └── utils/          # Utility functions
├── web/                # Web templates & static files
├── .env.example        # Example environment config
├── Makefile            # Common tasks (build, run, etc.)
├── tasks.json          # Sample tasks for JSON storage
└── go.mod              # Go modules
```

### How to Use

#### cli
- Fetch Task List

  Command：
  ```
  go run golang/tutorial/todo/cmd/todo list
  ```
  Result：
  ```
    ID                                   | Name                 | Status   | CreatedAt           | DueDate     | TimeLeft
  -----------------------------------------------------------------------------------------
  11               | Read book            | Pending  | 2025-08-24 18:40:00 | 2025-12-31  | 128 days       
  10               | Pay bills            | Pending  | 2025-08-24 18:40:00 | -           | -              
  9                | Buy milk             | Pending  | 2025-08-24 18:40:00 | 2025-09-01  | 7 days         
  8                | Read book            | Pending  | 2025-08-24 17:54:16 | 2025-12-31  | 128 days       
  ```

- File En/Decryption
  ```
  go run golang/tutorial/todo/cmd/todo encrypt -i {SrcFilePath} -o {DestFilePath} -p {PW}
  ```
  ```
  go run golang/tutorial/todo/cmd/todo decrypt -i {SrcFilePath} -o {DestFilePath} -p {PW}
  ```

#### Endpoints

* `GET /tasks/{id}` - Fetch a task
* `GET /tasks` - List tasks
* `POST /tasks` - Create task
* `GET /quote` - Get random quote
