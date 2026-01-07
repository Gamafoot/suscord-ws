FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src .

RUN CGO_ENABLED=0 go build -o /out/ws ./cmd/main


FROM alpine:3.19

WORKDIR /app

RUN addgroup -S nonroot \
 && adduser  -S nonroot -G nonroot

COPY config /app/config
COPY --from=builder /out/ws /app/ws

USER nonroot

ENV CONFIG_PATH=/app/config/config.yaml

EXPOSE 8001

ENTRYPOINT ["/app/ws"]
