#################################
## Builder Stage - SASS        ##
#################################
FROM dart:stable AS css-builder

ARG DART_SASS_VERSION=1.71.1
ARG DART_SASS_TAR=dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz
ARG DART_SASS_URL=https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/${DART_SASS_TAR}

WORKDIR /opt

ADD ${DART_SASS_URL} .
RUN tar -xzf ${DART_SASS_TAR} && rm ${DART_SASS_TAR}

WORKDIR /app
COPY pkg/monitor/web/static/zoe.scss  .
COPY pkg/service/web/web/static/dashboard.scss  .
RUN << __EOF__
/opt/dart-sass/sass zoe.scss zoe.css
/opt/dart-sass/sass dashboard.scss dashboard.css
__EOF__

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

COPY --from=css-builder /app/zoe.css pkg/monitor/web/static/zoe.css
COPY --from=css-builder /app/dashboard.css pkg/service/web/web/static/dashboard.css
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
