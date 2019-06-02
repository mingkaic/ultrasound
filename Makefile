.PHONY: up
up: build_db build_server
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: build_db
build_db:
	docker build -t ultra-db db

.PHONY: build_server
build_server:
	bazel run //server:ultrasound_server
