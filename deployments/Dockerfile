#################################
## Builder Stage               ##
#################################
FROM golang:1.22-alpine3.20 AS builder

WORKDIR /src

COPY  go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -o bin/zoe cmd/zoe/main.go

#################################
## Final Stage                 ##
#################################
FROM alpine:3.20

COPY --from=builder /src/bin/zoe /usr/local/bin/zoe

WORKDIR /app
COPY deployments/entrypoint.sh /app/entrypoint.sh

CMD ["/app/entrypoint.sh"]
