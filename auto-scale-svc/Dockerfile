ARG REGISTRY
ARG ALPINE_BASE_VERSION
FROM ${REGISTRY}/armorblox/alpine:${ALPINE_BASE_VERSION}

WORKDIR /app
COPY ./container/bin /app
ENV GOBIN=/usr/local/bin
EXPOSE 5555 8002
COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["bash", "-c", "/app/auto-scale-svc"]