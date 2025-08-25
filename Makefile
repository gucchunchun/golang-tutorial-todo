# もしコマンド名と同じファイルがあったとしてもコマンドとして実行されるようにする

.PHONY: hello nfile docker-up docker-up-build docker-down docker-clean docker-logs docker-ps docker-restart app-sh db-sh cli

COMPOSE       ?= docker compose
COMPOSE_FILE  ?= build/docker/compose/docker-compose.yml

hello:
	./scripts/hello_world.sh

nfile:
	mkdir -p "$(DIR)"
	touch "$(DIR)/$(FILE)"

# コンテナをスタートする
docker-up:
	$(COMPOSE) -f $(COMPOSE_FILE) up -d

# コンテナをビルドしてスタートする
docker-up-build:
	$(COMPOSE) -f $(COMPOSE_FILE) up -d --build

# コンテナを停止する
docker-down:
	$(COMPOSE) -f $(COMPOSE_FILE) down

# ボリューム/削除済みのコンテナを停止して削除する
docker-clean:
	$(COMPOSE) -f $(COMPOSE_FILE) down -v --remove-orphans

# ログの確認
docker-logs:
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f

docker-logs-app:
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f app

# ステータスの表示
docker-ps:
	$(COMPOSE) -f $(COMPOSE_FILE) ps

# コンテナを再起動する
docker-restart:
	$(COMPOSE) -f $(COMPOSE_FILE) restart

# シェルに入る (appコンテナ)
app-sh:
	$(COMPOSE) -f $(COMPOSE_FILE) exec app bash

# シェルに入る (dbコンテナ)
db-sh:
	$(COMPOSE) -f $(COMPOSE_FILE) exec db sh

# CLIアプリケーションを使用
cli:
	$(COMPOSE) -f $(COMPOSE_FILE) run --rm --entrypoint /app/todo app $(ARGS)
