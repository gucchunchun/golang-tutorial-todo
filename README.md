# 📌 **Golang ToDo CLI アプリ**

**Go 言語**で作るシンプルな **コマンドライン ToDo アプリケーション**
以下のスキルを練習できます：

* ✅ CLI プログラミング
* ✅ 入力バリデーション
* ✅ 時間処理
* ✅ 構造化されたコード設計

---

## 🎯 **プロジェクト目標**

JSON 永続化とオプションで CSV サポートを備えた **ToDo CLI アプリ**を構築する。

---

## ✅ **基本機能**

### 1. **タスク追加**

**コマンド例：**

```bash
todo add "買い物をする" --due 2025-07-28
```

**保存する情報：**

* **ID**（自動採番）
* **タスク名**
* **ステータス**（`pending` または `done`）
* **作成日時**
* **期限日**

✔ **タスク名**を **正規表現**でバリデーション（*Task 27*）

---

### 2. **タスク一覧表示**

**コマンド例：**

```bash
todo list
```

**表示フォーマット：**

```
ID | Name | Status | CreatedAt | DueDate | TimeLeft
```

✔ \*\*残り時間（TimeLeft）\*\*を Go の `time` パッケージで計算（*Task 28*）

---

### 3. **タスクを完了にする**

**コマンド例：**

```bash
todo done 3
```

➡ **ステータス**を `pending` → `done` に変更。

---

### 4. **タスク削除**

**コマンド例：**

```bash
todo delete 3
```

➡ 指定したタスクを削除。

---

## ✅ **データ永続化**

* **JSON ファイル**（`tasks.json`）にタスクを保存
* **encoding/json** を使用して読み書き

---

## 🔍 **追加機能**

* **タスクを CSV にエクスポート**

```bash
todo export tasks.csv
```

* **CSV からタスクをインポート**

```bash
todo import tasks.csv
```
---
---
# 📌 **Golang ToDo CLI App**

A simple **command-line ToDo application** in **Go** for practicing:

* ✅ CLI programming
* ✅ Input validation
* ✅ Time handling
* ✅ Structured code organization

---

## 🎯 **Project Goal**

Build a functional **ToDo CLI App** with JSON persistence and optional CSV support.

---

## ✅ **Core Features**

### 1. **Add Task**

**Command:**

```bash
todo add "Buy groceries" --due 2025-07-28
```

**Stored Fields:**

* **ID** (auto-increment)
* **Task Name**
* **Status** (`pending` or `done`)
* **CreatedAt** (timestamp)
* **DueDate**

✔ Validate **task name** using **regexp** (*Task 27*)

---

### 2. **List Tasks**

**Command:**

```bash
todo list
```

**Display Format:**

```
ID | Name | Status | CreatedAt | DueDate | TimeLeft
```

✔ Calculate **TimeLeft** using Go `time` package (*Task 28*)

---

### 3. **Mark Task as Done**

**Command:**

```bash
todo done 3
```

➡ Changes **status** from `pending` → `done`.

---

### 4. **Delete Task**

**Command:**

```bash
todo delete 3
```

➡ Removes the task from the list.

---

## ✅ **Data Persistence**

* Store all tasks in **JSON file** (`tasks.json`)
* Use **encoding/json** for read/write operations

---

## 🔍 **Additional Features**

* **Export tasks to CSV**

```bash
todo export tasks.csv
```

* **Import tasks from CSV**

```bash
todo import tasks.csv
```
