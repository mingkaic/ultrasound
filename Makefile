.PHONY: up
up: build_db build_server
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: restart
restart: down up

.PHONY: build_db
build_db:
	docker build -t mkaichen/ultra-db db

.PHONY: build_server
build_server:
	bazel run //server:ultrasound_server

.PHONY: push_db
push_db: build_db
	docker push mkaichen/ultra-db:latest

.PHONY: push_server
push_server: build_server
	bazel run //server:ultrasound_push
