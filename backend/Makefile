format:
	go fmt ./...

lint:
	go vet ./...
	golangci-lint run

migrate:
	echo ${GOOSE_DBSTRING}
	goose up