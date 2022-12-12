all: lint test

lint:
	go mod download
	docker run --rm \
		-v ${GOPATH}/pkg/mod:/go/pkg/mod \
 		-v ${PWD}:/app \
 		-w /app \
	    golangci/golangci-lint:v1.50 \
	    golangci-lint run -v --modules-download-mode=readonly

test:
	GOARCH=amd64 go test -gcflags='-N -l' ./...

doc:
	docker run --rm \
        -p 127.0.0.1:6060:6060 \
        -v ${PWD}:/go/src/github.com/evsamsonov/tinkoff-broker \
        -w /go/src/github.com/evsamsonov/tinkoff-broker \
        golang:latest \
        bash -c "go install golang.org/x/tools/cmd/godoc@latest && /go/bin/godoc -http=:6060"