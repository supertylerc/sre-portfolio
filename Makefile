BINARY_NAME=scheduler

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

run: build
	cp ${BINARY_NAME}-linux ${BINARY_NAME}
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME} 2> /dev/null || true
	rm ${BINARY_NAME}-darwin 2> /dev/null || true
	rm ${BINARY_NAME}-linux 2> /dev/null || true
	rm ${BINARY_NAME}-windows 2> /dev/null || true

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet ./...

lint:
	golangci-lint run --enable-all

format:
	go fmt ./...

fmt: format
	gci write --skip-vendor --skip-generated -s standard -s default . > /dev/null

tidy:
	go mod tidy

vendor: tidy
	go mod vendor
