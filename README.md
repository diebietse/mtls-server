# mTLS Server Example

[![build-status][build-badge]][build-link]

This is an example HTTPS server that generates valid [Let's Encrypt][lets-encrypt] certificates and validates connection with [Mutual TLS authentication (mTLS)][mtls].

The docker repository can be found at [diebietse/mtls-server][diebietse-docker].

![mtls-verified][art-mtls-verified] ![mtls-not-verified][art-mtls-not-verified]

## Quick Start

```sh
docker run --rm -e DEMO_FQDN=example.com -p80:80 -p443:443 docker pull diebietse/mtls-server
```

[build-badge]: https://github.com/diebietse/mtls-server/workflows/build/badge.svg
[build-link]: https://github.com/diebietse/mtls-server/actions?query=workflow%3Abuild
[mtls]: https://en.wikipedia.org/wiki/Mutual_authentication
[lets-encrypt]: https://letsencrypt.org/
[diebietse-docker]: https://hub.docker.com/repository/docker/diebietse/mtls-server
[art-mtls-verified]: https://github.com/diebietse/mtls-server/raw/master/art/mtls-verified-1024x768.png
[art-mtls-not-verified]: https://github.com/diebietse/mtls-server/raw/master/art/mtls-not-verified-1024x768.png
