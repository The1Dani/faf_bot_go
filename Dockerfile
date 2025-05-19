# base go image
FROM golang:latest AS builder

WORKDIR /

RUN mkdir /app

COPY . .

RUN CGO_ENABLED=0 go build -o /app/bot /cmd/bot

RUN chmod +x /app/bot

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/bot /app
COPY --from=builder /sql /app/sql

RUN apk add --no-cache ca-certificates

CMD [ "/app/bot" ]
