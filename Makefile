main_package_path = cmd/ublogging/main.go
binary_name = bin/ublogging

.PHONY: build
build:
	@echo "Building ublogging"
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

.PHONY: run-mongo
run-mongo:
	@echo "Running the mongo database"
	docker-compose --file docker-compose.dev.yml up -d mongo_db

.PHONY: run-mongo-no-daemon
run-mongo-no-daemon:
	@echo "Running the mongo database"
	docker-compose --file docker-compose.dev.yml up mongo_db

.PHONY: logs-mongo
logs-mongo:
	docker logs --follow ublogging_mongo

.PHONY: build-container
build-container:
	@echo "Here I should run my docker build command"

.PHONY: test
test: run-mongo
	@echo "Runinng Ublogging test suite"
	@echo "waiting 10 seconds for the containers to start running"
	@sleep 10
	go test ./...

.PHONY: http-create-user
http-create-user:
	@curl -X POST http://localhost:9099/api/users -H "Content-Type: application/json" -d '{ "username": "wenceslao1207", "email": "wmb1207@testing.com" }'

.PHONY: http-follow-user
http-follow-user:
	@curl -X POST http://localhost:9099/api/users/6732d083987d31f35eba4a6f/follow -H "Content-Type: application/json" -H "Authorization: 67341cf3237d4efe689e9771" | jq '.'

.PHONY: http-user
http-user:
	@curl -X GET "http://localhost:9099/api/users/feed?page=1&limit=2" -H "Content-Type: application/json" -H "Authorization: 67341cf3237d4efe689e9771" | jq '.'

.PHONY: http-post
http-post:
	@curl -X POST http://localhost:9099/api/posts -H "Content-Type: application/json" -H "Authorization: 6732d083987d31f35eba4a6f" -d '{"Content": "testing a second user that Im following" }' | jq '.'

.PHONY: http-comment
http-comment:
	@curl -X POST http://localhost:9099/api/posts/67349951a2bf4775831a9c86  -H "Content-Type: application/json" -H "Authorization: 67341cee237d4efe689e9770" -d '{"Content": "comment by lao"}' | jq '.'

.PHONY: http-comments
http-comments:
	@curl -X GET http://localhost:9099/api/posts/67349951a2bf4775831a9c86 -H "Content-Type: application/json" -H "Authorization: 67341cf3237d4efe689e9771" | jq '.'