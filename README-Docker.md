# xTeVe Docker Configuration Guide

This guide explains how to properly configure and deploy xTeVe using Docker with environment variable substitution.

## Environment Variables

The following environment variables are supported:

| Variable | Default | Description |
|----------|---------|-------------|
| `XTEVE_HOME` | `/home/xteve` | Home directory for xTeVe configuration and data |
| `XTEVE_PORT` | `34400` | Port on which xTeVe will listen |
| `XTEVE_TEMP` | `/tmp/xteve` | Temporary directory for buffer files |
| `TZ` | `UTC` | Timezone for the container |

## Template Processing

The Docker image includes an entrypoint script that processes template files with variables in the format `{XTEVE_*}`. 

### Supported Template Formats:
- `{XTEVE_HOME}` - Replaced with the value of XTEVE_HOME environment variable
- `{XTEVE_PORT}` - Replaced with the value of XTEVE_PORT environment variable
- `{XTEVE_TEMP}` - Replaced with the value of XTEVE_TEMP environment variable

### Template Files:
- Files with `.template` extension are automatically processed
- `settings.json` is also processed for template variables

## Docker Compose Configuration

### Development Environment

```yaml
version: '3.8'

services:
  xteve:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        - XTEVE_HOME=/home/xteve
        - XTEVE_PORT=34400
        - XTEVE_TEMP=/tmp/xteve
        - TZ=Europe/Moscow
    container_name: xteve
    restart: unless-stopped
    ports:
      - "34400:34400"
    environment:
      - TZ=Europe/Moscow
      - XTEVE_PORT=34400
      - XTEVE_HOME=/home/xteve
      - XTEVE_TEMP=/tmp/xteve
    volumes:
      - ./data:/home/xteve
      - ./temp:/tmp/xteve
      - /etc/localtime:/etc/localtime:ro
    networks:
      - xteve-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:${XTEVE_PORT:-34400}/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  xteve-network:
    driver: bridge
```

### Production Environment

```yaml
services:
  xteve:
    image: ghcr.io/ekzoman/xteve:latest
    restart: always
    environment:
      - TZ=Europe/Moscow
      - XTEVE_PORT=34400
      - XTEVE_HOME=/home/xteve
      - XTEVE_TEMP=/tmp/xteve
    ports:
      - 34400:34400
    networks:
      - xteve-network
    volumes:
      - xteve_data:/home/xteve
      - xteve_temp:/tmp/xteve
      - /etc/localtime:/etc/localtime:ro
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:${XTEVE_PORT:-34400}/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  xteve-network:
    driver: bridge

volumes:
  xteve_data:
  xteve_temp:
```

## Building the Image

```bash
# Build with default values
docker build -t xteve-fixed .

# Build with custom arguments
docker build \
  --build-arg XTEVE_HOME=/custom/home \
  --build-arg XTEVE_PORT=35000 \
  --build-arg TZ=Europe/Moscow \
  -t xteve-custom .
```

## Running the Container

### Basic Usage

```bash
docker run -d \
  --name xteve-test \
  -p 34400:34400 \
  -v $(pwd)/data:/home/xteve \
  -v $(pwd)/temp:/tmp/xteve \
  -e TZ=Europe/Moscow \
  xteve-fixed
```

### With Custom Variables

```bash
docker run -d \
  --name xteve-test \
  -p 35000:35000 \
  -v $(pwd)/data:/custom/home \
  -v $(pwd)/temp:/custom/temp \
  -e TZ=Europe/Moscow \
  -e XTEVE_HOME=/custom/home \
  -e XTEVE_PORT=35000 \
  -e XTEVE_TEMP=/custom/temp \
  xteve-fixed
```

## Troubleshooting

### Variable Substitution Issues

If you encounter issues with variable substitution:

1. Check the container logs for entrypoint messages:
   ```bash
   docker logs xteve-test
   ```

2. Verify environment variables are set correctly:
   ```bash
   docker exec xteve-test env | grep XTEVE_
   ```

3. Check if template files were processed:
   ```bash
   docker exec xteve-test ls -la /home/xteve/
   ```

### Common Issues

1. **Permission Issues**: Ensure the volumes mounted have correct permissions for the xteve user (UID 1001).

2. **Port Conflicts**: Make sure the host port is not already in use.

3. **Template Variables Not Replaced**: Verify that template files use the correct format `{XTEVE_*}` and not `${XTEVE_*}`.

## Advanced Configuration

### Custom Template Files

You can add your own template files by mounting them to the container:

```bash
docker run -d \
  --name xteve-test \
  -p 34400:34400 \
  -v $(pwd)/data:/home/xteve \
  -v $(pwd)/templates:/home/xteve/templates \
  -e TZ=Europe/Moscow \
  xteve-fixed
```

### Environment Variable Precedence

Environment variables are processed in the following order:

1. Docker build arguments (ARG) - used during image build
2. Docker environment variables (ENV) - used during container runtime
3. Docker run `-e` flags - override ENV variables
4. docker-compose environment section - override both

## Security Considerations

1. **Non-root User**: The container runs as a non-root user (xteve:1001) for security.

2. **Read-only Volumes**: Consider mounting configuration volumes as read-only where possible.

3. **Resource Limits**: Apply resource limits in production environments:

```yaml
deploy:
  resources:
    limits:
      cpus: '1.0'
      memory: 512M
    reservations:
      cpus: '0.5'
      memory: 256M
```

## Support

For issues related to:
- xTeVe functionality: [xTeVe Documentation](https://github.com/xteve-project/xTeVe-Documentation)
- Docker configuration: Create an issue in this repository
- Variable substitution: Check container logs first