FROM golang:1.25-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o sentinel ./cmd/agent

FROM alpine:3.20

RUN adduser -D -u 10001 appuser && apk add --no-cache ca-certificates 

WORKDIR /app
COPY --from=builder /app/sentinel /usr/local/bin/sentinel

USER appuser

# Set the executable as the container entrypoint
ENTRYPOINT ["sentinel"]