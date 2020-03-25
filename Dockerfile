FROM golang:1.13-alpine as golang

ARG ZSCALER_CERT

RUN apk add -U --no-cache ca-certificates git make build-base zip jq

RUN echo "$ZSCALER_CERT" > /usr/local/share/ca-certificates/zscaler.crt

ENV REQUESTS_CA_BUNDLE=/etc/ssl/certs/ca-certificates.crt

RUN update-ca-certificates --verbose

ENV AWSCLI_VERSION "1.18.8"

RUN apk --no-cache update && apk --no-cache upgrade

RUN apk add --update \
    python \
    python-dev \
    py-pip \
    build-base \
    bash \
    && pip install awscli==$AWSCLI_VERSION --upgrade --user \
    && apk --purge -v del py-pip \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["bash"]
