FROM golang:1.21.3 AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o alias-gen services/alias-gen/cmd/alias-gen/main.go

COPY build/wait-for-it.sh /src/wait-for-it.sh
RUN chmod +x /src/wait-for-it.sh

FROM debian:buster-slim

RUN apt-get update && \
    apt-get install -y netcat && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/alias-gen /alias-gen
COPY --from=builder /src/wait-for-it.sh /wait-for-it.sh

COPY services/alias-gen/config/ /config/

ENTRYPOINT ["/wait-for-it.sh", "postgres", "5432", "PostgreSQL", "--"]
CMD ["/alias-gen"]
