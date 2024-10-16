FROM golang:1.22.2-alpine3.19 AS builder
LABEL authors="D1mitrii"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/app

FROM scratch AS release

COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app"]