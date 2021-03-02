FROM golang:1.15 as builder
ENV GOPATH /opt/
COPY ./producer.go ./helper.go /opt/producer/
COPY ./metrics /opt/producer/metrics
RUN cd /opt/producer \
    && go mod init producer \
    && CGO_ENABLED=0 go build producer
FROM alpine:3.12
RUN mkdir /producer
WORKDIR /producer
COPY --from=builder /opt/producer/producer /producer/
CMD ["./producer"]