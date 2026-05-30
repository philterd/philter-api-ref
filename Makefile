build:
	go build -o philter-api-ref philter-api-ref.go

run:
	go run philter-api-ref.go

test:
	go test ./...
