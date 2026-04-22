FROM golang:1.24.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/hela-bank-sc main.go

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/hela-bank-sc /app/hela-bank-sc
COPY --from=builder /app/abi /app/abi

EXPOSE 8080

CMD ["./hela-bank-sc"]
