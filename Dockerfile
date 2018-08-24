FROM golang:1.10-alpine as builder
RUN apk add --update make git
WORKDIR /go/src/git.containerum.net/ch/solutions
COPY . .
RUN VERSION=$(git describe --abbrev=0 --tags) make build-for-docker

FROM alpine:3.7 as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM alpine:3.7

RUN apk --no-cache add tzdata ca-certificates
# app
COPY --from=builder /tmp/solutions /
# migrations
COPY pkg/migrations /migrations
# timezone data
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# tls certificates
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV CH_SOLUTIONS="impl" \
    CH_SOLUTIONS_MIGRATIONS_PATH="migrations" \
    CH_SOLUTIONS_TEXTLOG="true" \
    CH_SOLUTIONS_DEBUG="true" \
    CH_SOLUTIONS_DB="postgres" \
    CH_SOLUTIONS_PG_LOGIN="solutions" \
    CH_SOLUTIONS_PG_PASSWORD="" \
    CH_SOLUTIONS_PG_ADDR="postgres:5432" \
    CH_SOLUTIONS_PG_DBNAME="usermanager" \
    CH_SOLUTIONS_PG_NOSSL=true \
    CH_SOLUTIONS_KUBE_API_URL="http://kube-api:1214" \
    CH_SOLUTIONS_RESOURCE_URL="http://resource-service:1213"

EXPOSE 6767

ENTRYPOINT ["/solutions"]
