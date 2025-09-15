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

