FROM golang:1.25.5-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN apk --no-cache add ca-certificates

RUN go build -pgo=auto -ldflags="-s -w" -o ./main ./cmd

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["./main"]