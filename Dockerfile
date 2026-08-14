FROM golang:1.25.12 as builder

RUN mkdir /build
COPY . /build
RUN cd /build && make build

FROM alpine:3.22.1

RUN apk update --no-cache && apk upgrade --no-cache && apk add ca-certificates

COPY --from=builder /build/eventrouter /app/eventrouter

USER nobody:nobody

CMD ["/bin/sh", "-c", "/app/eventrouter -v 3 -logtostderr"]
