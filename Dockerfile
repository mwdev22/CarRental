FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

# ---------------------------------------------

FROM alpine:latest

# install deoendencies
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# copy binary from builder
COPY --from=builder /app/main .

# copy migrations
COPY --from=builder /app/migrations ./migrations


CMD ["./main"]
