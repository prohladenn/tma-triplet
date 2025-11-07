# Frontend Dockerfile
# Multi-stage build for optimized production image

# Build stage
FROM node:20-alpine AS builder

WORKDIR /build

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files first for better layer caching
COPY app/frontend/package.json app/frontend/pnpm-lock.yaml ./

# Install dependencies (cached if package files unchanged)
RUN pnpm install --frozen-lockfile --prefer-offline

# Copy source code
COPY app/frontend/ ./

# Build arguments for environment variables
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

# Build the application
RUN pnpm run build

# Runtime stage with nginx
FROM nginx:alpine

# Copy nginx configuration
COPY infra/nginx.conf /etc/nginx/nginx.conf

# Copy built files from builder
COPY --from=builder /build/dist /usr/share/nginx/html

# Remove default nginx config
RUN rm -rf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# Run nginx in foreground
CMD ["nginx", "-g", "daemon off;"]
