# =============================================================================
# Islamic School ERP - AI Gateway Dockerfile
# =============================================================================
# AI orchestration and proxy service for managing AI provider requests with
# caching, rate limiting, prompt routing, and request fallback.
# =============================================================================

# ---- Stage 1: Dependencies ----
FROM golang:1.23.4-alpine3.21 AS deps

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ---- Stage 2: Build ----
FROM deps AS builder

ARG APP_VERSION=1.0.0
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /src

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=${CGO_ENABLED} \
    GOOS=${GOOS} \
    GOARCH=${GOARCH} \
    go build \
      -ldflags="-s -w -X main.version=${APP_VERSION}" \
      -trimpath \
      -buildvcs=false \
      -o /app/ai-gateway \
      ./cmd/gateway

# ---- Stage 3: Production ----
FROM gcr.io/distroless/static-debian12:nonroot AS production

ARG APP_VERSION=1.0.0

LABEL org.opencontainers.image.title="Islamic School ERP AI Gateway"
LABEL org.opencontainers.image.description="AI orchestration gateway for Islamic School ERP"
LABEL org.opencontainers.image.version="${APP_VERSION}"

COPY --from=deps /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=deps /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
COPY --from=builder /app/ai-gateway /app/ai-gateway

USER 65532:65532

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD ["/app/ai-gateway", "health"] || exit 1

ENTRYPOINT ["/app/ai-gateway"]
CMD ["serve"]
