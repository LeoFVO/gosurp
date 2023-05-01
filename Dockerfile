FROM golang:latest AS builder
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod /src/
RUN go mod download
COPY . .
RUN go build -a -o <tool_name> -trimpath

FROM alpine:latest

RUN apk add --no-cache ca-certificates \
    && rm -rf /var/cache/*

RUN mkdir -p /app \
    && adduser -D <tool_name> \
    && chown -R <tool_name>:<tool_name> /app

USER <tool_name>
WORKDIR /app

COPY --from=builder /src/<tool_name> .

ENTRYPOINT [ "./<tool_name>" ]