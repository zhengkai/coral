start:
	go run *.go

lint:
	cd src && golangci-lint run ./...
