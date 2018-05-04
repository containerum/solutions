FROM golang:1.10-alpine as builder
WORKDIR /go/src/git.containerum.net/ch/solutions
COPY . .
WORKDIR cmd/solutions
RUN CGO_ENABLED=0 go build -v -ldflags="-w -s -extldflags '-static'" -tags="jsoniter" -o /bin/solutions

FROM alpine:3.7
RUN apk --no-cache add tzdata ca-certificates
# app
COPY --from=builder /bin/solutions /
# migrations
COPY pkg/migrations /migrations
ENV CH_SOLUTIONS="impl" \
    CH_SOLUTIONS_MIGRATIONS_PATH="migrations" \
    CH_SOLUTIONS_TEXTLOG="true" \
    CH_SOLUTIONS_DEBUG="true" \
    CH_SOLUTIONS_DB="postgres" \
    CH_SOLUTIONS_DB_URL="postgres://usermanager:ae9Oodai3aid@postgres:5432/solutions?sslmode=disable" \
    CH_SOLUTIONS_CSV_URL="https://pastebin.com/raw/JTGHiZWk" \
    CH_SOLUTIONS_KUBE_API_URL="http://kube-api:1214" \
    CH_SOLUTIONS_RESOURCE_URL="http://resource-service:1213" \
    CH_SOLUTIONS_CONVERTER_URL="http://model-converter:6543"
EXPOSE 6767
ENTRYPOINT ["/solutions"]
