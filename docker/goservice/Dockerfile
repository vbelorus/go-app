FROM golang:1.18 AS builder
COPY . /app
WORKDIR /app
ENV GO111MODULE=auto
#ENV GOPATH=/app
RUN go get github.com/ClickHouse/clickhouse-go/v2
RUN go install github.com/ClickHouse/clickhouse-go/v2
RUN cd src && CGO_ENABLED=0 GOOS=linux go build -o ../bin/appserver

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/appserver /
CMD ["/appserver"]