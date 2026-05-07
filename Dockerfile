FROM golang:1.25-alpine AS builder


WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/bin/sms-gateway ./cmd/sms-gateway


FROM alpine:3.22

RUN addgroup -S appGroup && adduser -S -G appGroup app

WORKDIR /app

COPY --from=builder /src/bin/sms-gateway ./sms-gateway

USER app

EXPOSE 8080

CMD ["./sms-gateway"]
