FROM docker.io/golang:1.23 AS build

ENV CGO_ENABLED=0

ARG PROJECT_PKG="github.com/Karzoug/meower-post-service"
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT_HASH=unknown

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

RUN --mount=type=cache,target=/root/.cache/go-build \
go build -trimpath -ldflags="-s -w \
-X '${PROJECT_PKG}/pkg/buildinfo.version=$VERSION' \
-X '${PROJECT_PKG}/pkg/buildinfo.buildDate=$BUILD_DATE' \
-X '${PROJECT_PKG}/pkg/buildinfo.commitHash=$COMMIT_HASH' \
" -o /bin/follow-service cmd/main.go

FROM scratch

USER 1001:1001

COPY --from=build --chown=1001 /bin/post-service /

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/post-service"]