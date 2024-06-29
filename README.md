# factorio-docker

Docker container for running [Factorio].

## Features

* Discord Bot via [Factocord]
* [Attested Docker images]

## Quickstart

### Docker CLI

```bash
# Create the storage directory for the server files and ensure it's
# owned by the container user.
# - data is for your server data
# - server_files is for the downloaded server only
sudo mkdir -p /opt/factorio{data,server_files}
sudo chown 845:845 /opt/factorio

docker run -d \
  -p 34197:34197/udp \
  -p 27015:27015/tcp \
  -v /opt/factorio/server_files:/opt/factorio \
  -v /opt/factorio/data:/data \
  -e VERSION=stable \
  --name factorio \
  --restart=unless-stopped \
  ghcr.io/jaredallard/factorio
```

### docker compose

```yaml
services:
  factorio:
    image: ghcr.io/jaredallard/factorio
    restart: unless-stopped
    environment:
      # Version can be: experimental, stable, or a specific version.
      VERSION: stable
      # Optional: Checksum to expect for the downloaded tar.xz. If not
      # set, it will be fetched from Factorio's site.
      SHA256SUM: ""
    ports:
      - 34197:34197/udp
      - 27015:27015/tcp
    volumes:
    - data:/data
    - server_files:/opt/factorio

volumes:
  data:
  server_files:
```

**Note**: Alternatively, you can remove the top-level `volumes` key and
replace `data:` with a path on your host machine. For more information
about Docker volumes, see [the docs](https://docs.docker.com/storage/volumes/).

## Versions

To see the available versions, check out the [Github Packages UI] for this
repository.

## Development

Prerequisites:

* [mise](https://mise.jdx.dev) (One can try to use host Go, but some
  scripts may not work without `mise` being installed)

First time/after-pulling:

```bash
mise install
```

Building:

```bash
mise run build
```

## FAQ

### Why not [factoriotools/factorio-docker]?

I don't have an amazing reason other than I building on top of that
image originally and found that using Bash was hitting it's limits
pretty quickly for adding support for things like Factocord. I'd prefer
to use Go to implement that, but upstream was Python and I didn't have a
ton of interest in using Python (sorry!). I also really believe that all
Docker images should be attested thanks to how easy Github has made it,
but their images are built with a Python script that sadly can't easily
do this :(

## Differences between [factoriotools/factorio-docker]

1. Docker images do NOT contain Factorio's server code.
  a) Why not? While it was considered to do this, it was decided to not
  because that would leave the base docker images vulnerable as they
  would likely rarely be updated. As such, it was decided that the
  version should always be downloaded at runtime once instead.

## License

AGPL-3.0

[Factorio]: https://www.factorio.com/
[Factocord]: https://github.com/maxsupermanhd/FactoCord-3.0
[Github Packages UI]: https://github.com/jaredallard/factorio-docker/packages
[Attested Docker Images]: https://docs.github.com/en/actions/security-guides/using-artifact-attestations-to-establish-provenance-for-builds#about-artifact-attestations
[factoriotools/factorio-docker]: https://github.com/factoriotools/factorio-docker
