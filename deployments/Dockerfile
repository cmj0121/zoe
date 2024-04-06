#################################
## Builder Stage               ##
#################################
FROM golang:alpine3.19 AS builder

WORKDIR /src
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

RUN apk add --no-cache build-base musl-dev sqlite-dev

COPY  go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -o zoe cmd/zoe/main.go

#################################
## Final Stage                 ##
#################################
FROM alpine:3.19

COPY --from=builder /src/zoe /usr/local/bin/zoe

WORKDIR /app

COPY assets/migrations assets/migrations
COPY deployments/launch.sh launch.sh
COPY config/example.yaml zoe.yaml

CMD ["/app/launch.sh"]
