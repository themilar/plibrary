run:
	go run ./cmd/api

up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database {} up

