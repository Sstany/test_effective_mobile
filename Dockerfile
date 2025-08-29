FROM golang:1.25-alpine AS builder


WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o core /app/cmd/core/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/core /app/core
COPY .env .env

EXPOSE 8080

CMD [ "./core" ]
