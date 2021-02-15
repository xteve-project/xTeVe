FROM alpine:latest
LABEL maintainer="No one special <noone@gmail.com>"

RUN apk update
RUN apk upgrade
RUN apk add --no-cache ca-certificates

# Dependencies
RUN apk add --no-cache curl bash tzdata
ENV TZ=Europe/Stockholm
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk add ffmpeg
RUN apk add vlc
RUN apk add xmltv 

# Remove APK cache
RUN rm -rf /var/cache/apk/*

# Add xteve binary
ADD https://github.com/martinvillysson/xTeVe-Downloads/raw/master/xteve_linux_amd64.zip /tmp/xteve_linux_amd64.zip

# Unzip the Binary
RUN mkdir -p /xteve
RUN unzip -o /tmp/xteve_linux_amd64.zip -d /xteve

# Clean up the .zip
RUN rm /tmp/xteve_linux_amd64.zip

# Add user for VLC and ffmpeg
RUN addgroup -S xteve && adduser -S xteve -G xteve

# Set executable permissions
RUN chmod +x /xteve/xteve
RUN chown xteve:xteve /xteve/xteve

# Set user contexts
USER xteve

#Create folder structure for backups and tmp files
RUN mkdir /home/xteve/.xteve/
RUN mkdir /home/xteve/.xteve/backup/
RUN mkdir /tmp/xteve
RUN mkdir /home/xteve/xmltvse
RUN mkdir /home/xteve/.xmltv

#Set Permission on folders
RUN chown xteve:xteve /home/xteve/.xteve/
RUN chown xteve:xteve /home/xteve/.xteve/backup/
RUN chown xteve:xteve /tmp/xteve
RUN chown xteve:xteve /home/xteve/xmltvse
RUN chown xteve:xteve /home/xteve/.xmltv

COPY tv_grab_se_tvzon.conf /home/xteve/.xmltv/
COPY grab_xml_tv_se.sh /home/xteve/xmltvse/

# Volumes
VOLUME /home/xteve/.xteve
VOLUME /home/xteve/xmltvse

# Expose Ports for Access
EXPOSE 34400

# Healthcheck
HEALTHCHECK --interval=30s --start-period=30s --retries=3 --timeout=10s \
  CMD curl -f http://localhost:34400/ || exit 1

# Entrypoint should be the base command
ENTRYPOINT ["/xteve/xteve"]

# Command should be the basic working
CMD ["-port=34400"]
