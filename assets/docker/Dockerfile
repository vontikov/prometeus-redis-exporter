FROM scratch

ARG VERSION

LABEL \
  org.opencontainers.image.title="Prometheus Redis Exporter" \
  org.opencontainers.image.description="\
Converts Redis INFO command into Prometheus metrics" \
  org.opencontainers.image.version="${VERSION}" \
  org.opencontainers.image.authors="Vladimir Ontikov <vontikov@gmail.com>"

COPY redis_exporter /usr/local/bin/
CMD ["redis_exporter"]
