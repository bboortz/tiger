
run: build
	docker run --rm -it --publish 8080:8080 tiger

build:
	go fmt .
	docker build -t tiger .

dev:
	go run main.go
