FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o philter-api-ref .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/philter-api-ref .

EXPOSE 8080

ENTRYPOINT ["./philter-api-ref"]
