#FROM golang:1.18.9 AS builder
FROM --platform=$BUILDPLATFORM golang:1.19.9-alpine AS builder

# Install tzdata and set the timezone
RUN apk update && apk add --no-cache tzdata
#ENV TZ=Asia/Jakarta

# Configure the time zone in the system
#RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go build -o bin/apps-users-management main.go

CMD ["/app/bin/apps-users-management"]

FROM builder AS dev-envs

RUN <<EOF
apk update
apk add git
EOF

RUN <<EOF
addgroup -S docker
adduser -S --shell /bin/bash --ingroup docker vscode
adduser -S --shell /bin/bash --ingroup docker vscode
EOF

# install Docker tools (cli, buildx, compose)
COPY --from=gloursdocker/docker / /

CMD ["go", "run", "main.go"]

FROM scratch AS app-release
COPY --from=builder /app/bin/apps-users-management /usr/local/bin/go-user-management/apps
COPY --from=builder /app/.env .

CMD ["/usr/local/bin/go-user-management/apps"]