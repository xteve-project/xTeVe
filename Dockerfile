# -----------------------------------------------------------------------------
# xTeVe Docker Image
#
# Builds a docker image for xTeVe M3U proxy server
# (https://xteve.de/)
#
# Multi-stage build for optimized image size
# -----------------------------------------------------------------------------

# First-stage build for xteve builder image
# -----------------------------------------------------------------------------
FROM golang:1.19-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o xteve .

# -----------------------------------------------------------------------------
# Second-stage build for xteve runtime image
# -----------------------------------------------------------------------------
FROM alpine:latest

# Install runtime dependencies including envsubst for template processing
RUN apk --no-cache add ca-certificates tzdata ffmpeg vlc gettext

# Create non-root user
RUN addgroup -g 1001 -S xteve && \
    adduser -u 1001 -S xteve -G xteve

# Set build arguments with default values
ARG XTEVE_HOME=/home/xteve
ARG XTEVE_PORT=34400
ARG XTEVE_TEMP=/tmp/xteve
ARG TZ=UTC

# Set environment variables
ENV XTEVE_HOME=${XTEVE_HOME}
ENV XTEVE_PORT=${XTEVE_PORT}
ENV XTEVE_TEMP=${XTEVE_TEMP}
ENV TZ=${TZ}

# Create necessary directories
RUN mkdir -p ${XTEVE_HOME} ${XTEVE_TEMP} && \
    chown -R xteve:xteve ${XTEVE_HOME} ${XTEVE_TEMP}

# Set working directory
WORKDIR ${XTEVE_HOME}

# Copy binary from builder stage
COPY --from=builder /build/xteve /usr/local/bin/xteve

# Copy static files
COPY --from=builder /build/html ./html

# Copy configuration templates
COPY config/settings.json.template ${XTEVE_HOME}/settings.json.template

# Copy entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

# Make entrypoint script executable
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Change ownership
RUN chown -R xteve:xteve /usr/local/bin/xteve ${XTEVE_HOME} /usr/local/bin/docker-entrypoint.sh

# Switch to non-root user
USER xteve

# Expose port
EXPOSE ${XTEVE_PORT}

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${XTEVE_PORT}/ || exit 1

# Run the application through entrypoint script
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["-config", "${XTEVE_HOME}"]