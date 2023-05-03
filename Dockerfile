FROM golang:latest AS builder
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod /src/
RUN go mod download
COPY . .
RUN go build -a -o gosurp -trimpath

FROM alpine:latest

RUN apk add --no-cache ca-certificates \
    && rm -rf /var/cache/*

RUN mkdir -p /app \
    && adduser -D gosurp \
    && chown -R gosurp:gosurp /app

USER gosurp
WORKDIR /app

COPY --from=builder /src/gosurp .

ENTRYPOINT [ "./gosurp" ]