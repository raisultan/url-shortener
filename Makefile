run-tests:
	@echo "Running e2e tests..."
	go test -v ./services/main/tests

run-postgres:
	@echo "Running PostgreSQL container..."
	docker run --name alias-gen-postgres \
	-e POSTGRES_USER=alias-gen \
	-e POSTGRES_PASSWORD=alias-gen \
	-e POSTGRES_DB=url-aliases -d -p 5432:5432 postgres

run-redis:
	@echo "Running Redis container..."
	docker run --name url-shortener-redis -d -p 6379:6379 redis

create-docker-network:
	@echo "Creating Docker network for ClickHouse and Metabase..."
	docker network create url-shortener

run-clickhouse:
	@echo "Running ClickHouse container..."
	docker run -d --name clickhouse-server \
	--ulimit nofile=262144:262144 \
	-p 9000:9000 -p 8123:8123 \
	-e CLICKHOUSE_DB=testing \
	-e CLICKHOUSE_SERVER__LISTEN_HOST='0.0.0.0' \
	--network url-shortener \
	yandex/clickhouse-server

setup-metabase-plugins:
	@echo "Prepping Clickhouse plugin for Metabase..."
	mkdir -p mb/plugins
	curl -L -o mb/plugins/clickhouse.jar https://github.com/ClickHouse/metabase-clickhouse-driver/releases/download/1.2.2/clickhouse.metabase-driver.jar

run-metabase:
	@echo "Running Metabase container..."
	docker run -d -p 3000:3000 \
	--network url-shortener \
	--mount type=bind,source=$(shell pwd)/mb/plugins/clickhouse.jar,destination=/plugins/clickhouse.jar \
	metabase/metabase:v0.47.2

run-alias-gen:
	@echo "Running alias-gen service..."
	go run services/alias-gen/cmd/alias-gen/main.go

run-main:
	@echo "Running main service..."
	go run services/main/cmd/url-shortener/main.go
