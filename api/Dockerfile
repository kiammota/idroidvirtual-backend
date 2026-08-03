FROM golang:1.26-alpine AS builder
ENV TZ=America/Sao_Paulo
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api .

EXPOSE 7192

CMD ["./api"]