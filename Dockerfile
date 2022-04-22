FROM golang:1.16 as builder

RUN mkdir /build
COPY . /build
RUN cd /build && make build

FROM gcr.io/pingcap-public/pingcap/alpine-glibc:alpine-3.14.3

RUN apk update --no-cache && apk add ca-certificates

COPY --from=builder /build/eventrouter /app/eventrouter

USER nobody:nobody

CMD ["/bin/sh", "-c", "/app/eventrouter -v 3 -logtostderr"]
