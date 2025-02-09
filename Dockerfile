# Step 1: Build the Go application
FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN make build

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/file_storage .

COPY .env .env

EXPOSE 8080

CMD ["./file_storage"]
