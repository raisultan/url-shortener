FROM golang:1.21.3 AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ./url-shortener ./services/main/cmd/url-shortener

COPY services/main/build/wait-for-it.sh /src/wait-for-it.sh
RUN chmod +x /src/wait-for-it.sh

FROM debian:buster-slim

RUN apt-get update && \
    apt-get install -y netcat && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/url-shortener /url-shortener
COPY --from=builder /src/wait-for-it.sh /wait-for-it.sh

COPY services/main/config/ /config/

ENTRYPOINT ["/wait-for-it.sh", "clickhouse", "--"]
CMD ["/url-shortener"]
