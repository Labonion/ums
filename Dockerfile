FROM golang:1.21.6-alpine AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/src/app/main .

COPY .env .

EXPOSE 8080

CMD ["./main"]
