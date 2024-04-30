ENV_PATH = $(PWD)/.env
include $(ENV_PATH)


DOWN =
UP =
FORCE =
SEQ = 

dev:
	go run cmd/main.go
build:
	go build -o bin/app cmd/main.go
start:
	./bin/app
swag:
	swag init -g cmd/main.go 
migrate-up:
	migrate -path ./migrations -database '$(DB_URL)?sslmode=disable' up $(UP)
migrate-down:
	migrate -path ./migrations -database '$(DB_URL)?sslmode=disable' down $(DOWN)
migrate-force:
	migrate -path ./migrations -database '$(DB_URL)?sslmode=disable' force $(FORCE)
create-migration:
	migrate create -ext sql -dir ./migrations -seq $(SEQ)