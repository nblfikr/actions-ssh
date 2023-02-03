FROM golang:1.19.4-alpine as builder

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine:3.17

COPY --from=builder /go/bin/app /

ENTRYPOINT ["/app"]