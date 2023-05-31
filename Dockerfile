FROM golang:1.20 AS builder

WORKDIR /build

COPY cmd cmd
COPY server server
COPY vendor vendor
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile

RUN CGO_ENABLED=0 go build -mod=vendor ./cmd/mtls-server

FROM cfssl/cfssl as cert-builder

WORKDIR /certs/

COPY certgen certgen

RUN cd certgen && chmod +x genca.sh && chmod +x gencert.sh

RUN  cd certgen && ./genca.sh && ./gencert.sh mtls-example-client

FROM alpine

WORKDIR /site

COPY --from=builder /build/mtls-server /bin/mtls-server
COPY site /site
COPY --from=cert-builder /certs/certgen/mtls-example-client.p12 .
COPY --from=cert-builder /certs/certgen/root.pem .

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["mtls-server"]