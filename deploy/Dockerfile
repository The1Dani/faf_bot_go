FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY ./bot /app
COPY ./sql /app/sql

RUN apk add --no-cache ca-certificates

CMD [ "/app/bot" ]