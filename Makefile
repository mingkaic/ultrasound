.PHONY: up
up: build
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: restart
restart: down up

.PHONY: build
build: build_db build_server

.PHONY: push
push: push_db push_server

.PHONY: build_db
build_db:
	docker build -t mkaichen/ultra-db db

.PHONY: build_server
build_server:
	bazel run //server:ultrasound_server
	docker tag bazel/server:ultrasound_server mkaichen/ultrasound_server:latest

.PHONY: push_db
push_db: build_db
	docker push mkaichen/ultra-db:latest

.PHONY: push_server
push_server: build_server
	bazel run //server:ultrasound_push
