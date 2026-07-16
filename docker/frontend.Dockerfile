# =============================================================================
# Islamic School ERP - Frontend Nuxt 4 Dockerfile
# =============================================================================
# Multi-stage build: dependencies → build → production with Nginx + SSR
# =============================================================================

# ---- Stage 1: Install dependencies ----
FROM node:22.12.0-alpine3.21 AS deps

WORKDIR /app

# Copy package files for layer caching
COPY package.json package-lock.json* yarn.lock* pnpm-lock.yaml* .npmrc* ./

# Use corepack to respect packageManager field
RUN corepack enable && corepack prepare --activate

RUN --mount=type=cache,target=/root/.npm \
    if [ -f pnpm-lock.yaml ]; then \
      pnpm install --frozen-lockfile --prod; \
    elif [ -f yarn.lock ]; then \
      yarn install --frozen-lockfile --production; \
    else \
      npm ci --omit=dev; \
    fi

# ---- Stage 2: Build the application ----
FROM node:22.12.0-alpine3.21 AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Set build-time env vars
ARG APP_VERSION=1.0.0
ARG API_BASE_URL=https://api.erp-school.id
ENV NUXT_PUBLIC_API_BASE=$API_BASE_URL
ENV NUXT_PUBLIC_APP_VERSION=$APP_VERSION

RUN corepack enable && corepack prepare --activate

# Build Nuxt app (generates .output directory with Nitro server)
RUN --mount=type=cache,target=/root/.npm \
    if [ -f pnpm-lock.yaml ]; then \
      pnpm run build; \
    elif [ -f yarn.lock ]; then \
      yarn build; \
    else \
      npm run build; \
    fi

# ---- Stage 3: Production - Node.js SSR server ----
FROM node:22.12.0-alpine3.21 AS production

ARG APP_VERSION=1.0.0

LABEL org.opencontainers.image.title="Islamic School ERP Frontend"
LABEL org.opencontainers.image.description="Nuxt 4 frontend for Islamic School ERP system"
LABEL org.opencontainers.image.version="${APP_VERSION}"
LABEL org.opencontainers.image.source="https://github.com/opencode/erp-school-frontend"

RUN apk add --no-cache tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    echo "Asia/Jakarta" > /etc/timezone

# Create non-root user
RUN addgroup -g 1001 -S nuxt && \
    adduser -S nuxt -u 1001 -G nuxt

WORKDIR /app

# Copy only the runtime output from the build stage
COPY --from=builder /app/.output ./.output

# Set Node.js production optimizations
ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000
ENV HOST=0.0.0.0
ENV PORT=3000

USER nuxt

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
  CMD curl -f http://localhost:3000/api/_health || exit 1

ENTRYPOINT ["node", ".output/server/index.mjs"]
