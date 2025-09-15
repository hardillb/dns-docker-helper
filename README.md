# DNS-DOCKER-HELPER

A tool to dynamically update DNS entries for Docker containers exposed
using nginx-proxy/nginx-proxy container

## Configuration

The following environment variables need to be in scope

- `DNS_SERVER` IP address and port e.g. `192.168.1.1:53`
- `CNAME_TARGET` hostname as target for CNAME
- `KEY_NAME` TSIG Key name
- `KEY_SECRET` Base64 encode TSIG Key
- `ZONE` The zone name on the DNS Server

### Docker

```
docker run -d -v /var/run/docker.sock:/var/run/docker.sock \
 -e DNS_SERVER=192.168.1.1:53 \
 -e CNAME_TARGET=example.com \
 -e KEY_NAME=docker \
 -e KEY_SECRET=ADssd12343== \
 -e ZONE=example.com \
 hardillb/dns-docker-helper
```

But would normally be run in a `docker-compose.yml` file along side nginx-proxy/nginx-proxy