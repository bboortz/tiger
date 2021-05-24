FROM golang:alpine as builder
WORKDIR /app/
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o main
RUN wget -O upx.tar.xz https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz && \
    tar xvf upx.tar.xz && \
    mv upx-3.96-amd64_linux upx && \
    ./upx/upx ./main

FROM scratch
EXPOSE 8080
WORKDIR /app/
COPY --from=builder /app/main /app/
CMD ["/app/main"]
