FROM bitnami/minideb
COPY ./build/bin/gateway /
CMD ["/gateway", "-env=product"]