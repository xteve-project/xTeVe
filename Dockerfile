# First stage. Building a binary
# -----------------------------------------------------------------------------

# Base image for builder is debian 11 with golang 1.18+ pre-installed
FROM --platform=$BUILDPLATFORM golang:bullseye AS builder

# Download the source code
# Uncomment the below line to force git pull (no cache)
#ADD "https://www.random.org/cgi-bin/randbyte?nbytes=10&format=h" skipcache
RUN git clone https://github.com/SenexCrenshaw/xTeVe.git /src
WORKDIR /src

ARG TARGETOS TARGETARCH
# Install dependencies
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go mod download

# Compile
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build xteve.go

# Second stage. Creating an image
# -----------------------------------------------------------------------------
FROM alpine:latest

ARG BUILD_DATE
ARG VCS_REF
ARG XTEVE_PORT=34400
ARG XTEVE_VERSION=2.5.1

LABEL org.opencontainers.image.created="{$BUILD_DATE}" \
      org.opencontainers.image.url="https://hub.docker.com/r/SenexCrenshaw/xteve/" \
      org.opencontainers.image.source="https://github.com/SenexCrenshaw/xTeVe" \
      org.opencontainers.image.version="{$XTEVE_VERSION}" \
      org.opencontainers.image.revision="{$VCS_REF}" \
      org.opencontainers.image.vendor="SenexCrenshaw" \
      org.opencontainers.image.title="xTeVe" \
      org.opencontainers.image.description="Dockerized fork of xTeVe by SenexCrenshaw" \
      org.opencontainers.image.authors="SenexCrenshaw SenexCrenshaw@gmail.com"

ENV XTEVE_BIN=/home/xteve/bin
ENV XTEVE_CONF=/home/xteve/conf
ENV XTEVE_HOME=/home/xteve
ENV XTEVE_TEMP=/tmp/xteve

# Add binary to PATH
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:$XTEVE_BIN

# Set working directory
WORKDIR $XTEVE_HOME

# Update package lists
RUN apk update
RUN apk upgrade

# Install CA certificates
RUN apk add --no-cache ca-certificates
RUN apk add curl

# Timezone (TZ)
RUN apk update && apk add --no-cache tzdata
ENV TZ=America/New_York
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Add ffmpeg and vlc
RUN apk add ffmpeg
RUN apk add vlc

# Creat bin dir
RUN mkdir $XTEVE_BIN

# Copy built binary from builder image
COPY --from=builder [ "/src/xteve", "${XTEVE_BIN}/" ]

# Set binary permissions
RUN chmod +rx $XTEVE_BIN/xteve

# Create XML cache directory
RUN mkdir $XTEVE_HOME/cache

# Create working directories for xTeVe
RUN mkdir $XTEVE_CONF
RUN chmod a+rwX $XTEVE_CONF
RUN mkdir $XTEVE_TEMP
RUN chmod a+rwX $XTEVE_TEMP

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Configure container volume mappings
VOLUME $XTEVE_CONF
VOLUME $XTEVE_TEMP

# Expose Port
EXPOSE 34400

# Run the xTeVe executable
ENTRYPOINT ${XTEVE_BIN}/xteve -port=${XTEVE_PORT} -config=${XTEVE_CONF}
