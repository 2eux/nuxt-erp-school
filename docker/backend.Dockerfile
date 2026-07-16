# =============================================================================
# Islamic School ERP - Backend Go API Dockerfile
# =============================================================================
# Multi-stage build producing a minimal distroless production image.
# Build: docker build --build-arg GO_VERSION=1.23 --build-arg APP_VERSION=1.0.0 -f docker/backend.Dockerfile -t erp-school-backend ./backend
# =============================================================================

# ---- Stage 1: Build dependencies cache (go modules) ----
FROM golang:1.23.4-alpine3.21 AS deps

ARG GO_VERSION=1.23

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Copy only go.mod/go.sum first for layer caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ---- Stage 2: Build the application ----
FROM deps AS builder

ARG APP_VERSION=1.0.0
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /src

# Copy source code
COPY . .

# Build with optimizations:
#   -ldflags="-s -w"     : strip debug info, reduce binary size (~30-40%)
#   -trimpath            : remove filesystem paths from binary
#   -buildvcs=false      : skip VCS stamping for reproducible builds
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=${CGO_ENABLED} \
    GOOS=${GOOS} \
    GOARCH=${GOARCH} \
    go build \
      -ldflags="-s -w -X main.version=${APP_VERSION} -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
      -trimpath \
      -buildvcs=false \
      -o /app/server \
      ./cmd/server

# ---- Stage 3: Run security vulnerability scan ----
FROM aquasec/trivy:0.57.1 AS trivy-scanner

COPY --from=builder /app/server /scan/server
RUN trivy filesystem --severity HIGH,CRITICAL --no-progress --ignore-unfixed /scan

# ---- Stage 4: Distroless production image ----
FROM gcr.io/distroless/static-debian12:nonroot AS production

ARG APP_VERSION=1.0.0

LABEL org.opencontainers.image.title="Islamic School ERP Backend"
LABEL org.opencontainers.image.description="Go API server for Islamic School ERP system"
LABEL org.opencontainers.image.version="${APP_VERSION}"
LABEL org.opencontainers.image.vendor="Islamic School ERP"
LABEL org.opencontainers.image.source="https://github.com/opencode/erp-school-backend"

# Copy CA certificates from deps stage for TLS connections
COPY --from=deps /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=deps /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

# Copy the compiled binary
COPY --from=builder /app/server /app/server

# Copy migrations (embedded or mounted)
COPY --from=builder /src/migrations /app/migrations

# Distroless nonroot user (UID 65532)
USER 65532:65532

EXPOSE 8080

# Verify binary is executable
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD ["/app/server", "health"] || exit 1

ENTRYPOINT ["/app/server"]
CMD ["serve"]
