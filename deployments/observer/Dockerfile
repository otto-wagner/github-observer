FROM golang:1.25-alpine AS builder

WORKDIR /build

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV CC=gcc

RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-w -s" -o observer ./cmd/observer

FROM alpine:3.23

WORKDIR /observer
COPY --from=builder /build/observer ./observer

RUN mkdir -p /backend/data

ENV PORT=80
EXPOSE $PORT

HEALTHCHECK --interval=30s --timeout=5s --retries=3 --start-period=20s \
  CMD wget -qO- http://localhost:$PORT/health || exit 1

CMD ["./observer", "server"]
