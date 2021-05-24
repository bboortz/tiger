
run: dockerbuild
	docker run --rm -it --publish 8080:8080 tiger

dockerbuild: fmt test
	docker build -t tiger .

build: fmt test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -trimpath ./...

test:
	go test -v ./...

fmt:
	go fmt ./...

startdev:
	go run main.go
