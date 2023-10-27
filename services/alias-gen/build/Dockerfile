FROM golang:1.21.3 AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o alias-gen services/alias-gen/cmd/alias-gen/main.go

FROM gcr.io/distroless/base-debian10

COPY --from=builder /src/alias-gen /alias-gen

COPY services/alias-gen/config/ /config/

CMD ["/alias-gen"]
