main_package_path = cmd/ublogging/main.go
bynary_name = bin/ublogging


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

