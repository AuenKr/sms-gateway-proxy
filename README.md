# SMS Gateway

This project is a small proxy service for sending SMS requests to another SMS provider.

Main project: <https://github.com/capcom6/android-sms-gateway>

API docs: <https://api.sms-gate.app/docs/index.html>

## Use cases

- keep the main SMS gateway username and password hidden from other users or apps
- allow only users with the correct `SERVER_AUTH_KEY` to use this service
- let users send requests like a normal message endpoint without handling encryption themselves
- encrypt the message data in this service before forwarding it to the upstream SMS gateway

## What it does

- receives `POST /messages` requests
- checks that the request is valid
- protects the API with an auth header
- encrypts message data before sending it on
- forwards the request to the real SMS gateway
- provides `GET /health` to confirm the service is running

## Required settings

Create a `.env` file with these values:

- `GATEWAY_BASE_URL`
- `USERNAME` for the upstream SMS gateway
- `PASSWORD` for the upstream SMS gateway
- `PASSPHRASE_KEY` used to encrypt message data before forwarding
- `SERVER_AUTH_KEY` required in the `Authorization` header when calling this service

Optional values:

- `PORT` default: `8080`
- `DEV_MODE` default: `false` (`false` encrypts data, `true` skips encryption)
- `PBKDF2_ITERATIONS` default: `1400000`

## Run

Docker image: `auenkr/sms-gateway-proxy`

### Run with Docker

```bash
docker pull auenkr/sms-gateway-proxy:latest
docker run --env-file .env -p 8080:8080 auenkr/sms-gateway-proxy:latest
```

### Run locally

```bash
go run ./cmd/sms-gateway
```
