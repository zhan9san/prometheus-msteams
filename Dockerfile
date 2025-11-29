FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk --no-cache add git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
ARG VERSION
ARG COMMIT
ARG BRANCH
ARG BUILD_DATE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X github.com/zhan9san/prometheus-msteams/pkg/version.VERSION=${VERSION} \
              -X github.com/zhan9san/prometheus-msteams/pkg/version.COMMIT=${COMMIT} \
              -X github.com/zhan9san/prometheus-msteams/pkg/version.BRANCH=${BRANCH} \
              -X github.com/zhan9san/prometheus-msteams/pkg/version.BUILDDATE=${BUILD_DATE}" \
    -o prometheus-msteams ./cmd/server

FROM alpine:latest AS certs

RUN apk --no-cache add ca-certificates tzdata

FROM scratch
LABEL description="A lightweight Go Web Server that accepts POST alert message from Prometheus Alertmanager and sends it to Microsoft Teams Channels using an incoming webhook url."
EXPOSE 2000

# Copy required cert and zoneinfo from previous stage
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certs /usr/share/zoneinfo /usr/share/zoneinfo

COPY ./default-message-card.tmpl /default-message-card.tmpl
COPY ./default-message-workflow-card.tmpl /default-message-workflow-card.tmpl
COPY --from=builder /build/prometheus-msteams /promteams

ENTRYPOINT ["/promteams"]
