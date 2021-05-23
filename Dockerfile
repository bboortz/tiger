FROM golang:alpine as builder
WORKDIR /app/
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM scratch
EXPOSE 8080
WORKDIR /app/
COPY --from=builder /app/main /app/
CMD ["/app/main"]
