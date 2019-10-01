FROM golang AS builder

WORKDIR /root

# Dependencies
RUN go get github.com/koron/go-ssdp && \
	go get github.com/gorilla/websocket && \
	go get github.com/kardianos/osext

# Copy of source files
COPY . .

# Disable CGO Tool
ENV CGO_ENABLED=0

# xTeVe build
RUN go build xteve.go

#                    #
# xTeVe docker image #
#                    #
FROM alpine:latest  

CMD ["/usr/local/bin/xteve", "-config", "/data"] 

ENV UID=1000
ENV GID=1000

EXPOSE 34400   

# User creation and installation of ca-certificates
RUN addgroup -g $UID -S xteve  && \
	adduser -u $GID -S xteve -G xteve && \
	mkdir /data && \
	chown $UID:$GID /data && \
	apk add --no-cache ca-certificates && \
	rm -rf /var/cache/apk/*

# Copy binary from build stage
COPY --from=builder /root/xteve /usr/local/bin/

USER xteve

HEALTHCHECK --interval=1m --timeout=3s \
  CMD curl --fail http://localhost:33440/web || exit 1

VOLUME ["/data"]