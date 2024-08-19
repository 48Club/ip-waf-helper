FROM golang:alpine

COPY . /app

WORKDIR /app

RUN go build -o ip-waf-helper .

FROM alpine

COPY --from=0 /app/ip-waf-helper /app/ip-waf-helper

RUN apk --no-cache add tzdata

ENV MARIADB_ROOT_PASSWORD=\
    TZ=\
    GIN_TOKEN=

ENTRYPOINT ["/app/ip-waf-helper"]