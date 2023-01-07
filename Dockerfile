FROM golang:alpine as app-builder
WORKDIR /go/src/app
COPY . .
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o iotd ./cmd/IOTD/main.go

FROM scratch
COPY --from=app-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=app-builder /go/src/app/iotd /iotd
ENTRYPOINT ["/iotd"]