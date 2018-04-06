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
    CH_SOLUTIONS_CSV_URL="https://raw.githubusercontent.com/containerum/solution-list/master/containerum-solutions.csv"
EXPOSE 6666
ENTRYPOINT ["/solutions"]
