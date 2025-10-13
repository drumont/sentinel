FROM golang:1.25-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o sentinel ./cmd/agent

FROM alpine:3.20

RUN adduser -D -u 10001 sentinel && apk add --no-cache ca-certificates nmap

WORKDIR /app
COPY --from=builder /app/sentinel /usr/local/bin/sentinel

RUN chown sentinel:sentinel -R /app
USER sentinel

# Set the executable as the container entrypoint
ENTRYPOINT ["sentinel"]