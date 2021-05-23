
run: dockerbuild
	docker run --rm -it --publish 8080:8080 tiger

dockerbuild: fmt test
	docker build -t tiger .

build: fmt test
	go build -v ./...

test:
	go test -v ./...

fmt:
	go fmt ./...

startdev:
	go run main.go
