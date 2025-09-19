# DNS-DOCKER-HELPER

A tool to dynamically update DNS entries for Docker containers exposed
using [nginx-proxy/nginx-proxy](https://github.com/nginx-proxy/nginx-proxy) container.

This application watch for when containers are added or removed and will take the `VIRTUAL_HOST` environment variables on containers used to configure nginx-proxy/nginx-proxy and create a matching DNS CNAME entry.

## Build

```
docker build . -t dns-docker-helper
```

## Configuration

The following environment variables need to be in scope

- `DNS_SERVER` IP address and port e.g. `192.168.1.1:53`
- `CNAME_TARGET` hostname as target for CNAME
- `KEY_NAME` TSIG Key name
- `KEY_SECRET` Base64 encode TSIG Key
- `ZONE` The zone name on the DNS Server

### Docker

```bash
docker run -d -v /var/run/docker.sock:/var/run/docker.sock \
 -e DNS_SERVER=192.168.1.1:53 \
 -e CNAME_TARGET=example.com \
 -e KEY_NAME=docker \
 -e KEY_SECRET=ADssd12343== \
 -e ZONE=example.com \
 hardillb/dns-docker-helper
```

### Docker Compose

```yaml
services:
  nginx:
    image: nginxproxy/nginx-proxy:latest
    restart: always
    volumes:
      - "/var/run/docker.sock:/tmp/docker.sock:ro"
    ports:
      - "80:80"
      - "443:443"
    environment:
      - VIRTUAL_HOST=www.example.com
  dns:
    image: hardillb/dns-docker-helper:latest
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
    environment:
      - DNS_SERVer=192.168.1.1:53
      - CNAME_TARGET=example.com
      - KEY_NAME=docker
      - KEY_SECRET=ADssd12343==
      - ZONE=example.com
```

But would normally be run in a `docker-compose.yml` file along side nginx-proxy/nginx-proxy