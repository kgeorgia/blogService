all:
	docker-compose up

down:
	docker-compose down

re:
	docker-compose up --build

build:
	@go build -o blogService cmd/app/main.go
	@echo "blogService is compiled!"

run:
	@go run -o blogService cmd/app/main.go

test:
	@echo "TODO"

clean:
	@rm blogService

.PHONY: down, re, build, test, clean