include .env

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]
## run/api: run the cmd/api application
run/api:
	go run ./cmd/api

## db/migrations/new name=$1: create a new database migration
db/migrations/new:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	@migrate -path ./migrations -database ${DATABASE_URL} up

## build/api: build the cmd/api application
build/api:
	@echo 'Building cmd/api...'
	go build -o=./bin/api ./cmd/api