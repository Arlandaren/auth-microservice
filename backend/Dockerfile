FROM golang:1.23.2-alpine AS builder

RUN apk update && apk add --no-cache ca-certificates git gcc g++ libc-dev binutils

WORKDIR /opt

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ../.. .

RUN go build -v -o bin/grpc_server ./cmd/grpc_server

FROM alpine AS runner

RUN apk update && apk add --no-cache ca-certificates libc6-compat bash && rm -rf /var/cache/apk/**

WORKDIR /opt

COPY --from=builder /opt/bin/grpc_server ./

CMD ["./grpc_server"]
