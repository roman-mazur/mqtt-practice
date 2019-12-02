FROM golang:1.12 as build

RUN go get github.com/eclipse/paho.mqtt.golang

WORKDIR /go/src/inceptor
COPY . .

RUN go install -ldflags '-extldflags "-fno-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build' ./

FROM alpine:latest as ca-certs
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=ca-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/inceptor /inceptor
ENTRYPOINT ["/inceptor"]
