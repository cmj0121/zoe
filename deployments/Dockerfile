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

RUN apk add --no-cache envsubst

COPY --from=builder /src/bin/zoe /usr/local/bin/zoe

# downgrade the non-root user
RUN addgroup -g 1000 -S app && adduser -u 1000 -S app -G app
USER app

WORKDIR /home/app
COPY deployments/zoe.yml /app/zoe.yml
COPY deployments/entrypoint.sh /app/entrypoint.sh

CMD ["/app/entrypoint.sh"]
