FROM golang:1.14-alpine AS builder

ENV CGO_ENABLED=0
WORKDIR /build
COPY proxy.go proxy.go
RUN apk --update add ca-certificates
RUN go build -o proxy .


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/proxy /proxy
EXPOSE 1337
ENTRYPOINT ["/proxy"]
