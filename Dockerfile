FROM golang:1.17-alpine

RUN mkdir /app
COPY . /app

RUN go build -o stockcollect .



FROM alpine:3.14

RUN apk update && apk add tzdata && \
	cp /usr/share/zoneinfo/America/New_York /etc/localtime && \
	echo "America/New_York" > /etc/timezone && \
	apk del tzdata && rm -rf /var/cache/apk/*

COPY ./crontab /etc/crontabs/root

RUN mkdir /app
COPY --from=0 /app/stockcollect /app/

CMD chown root:root /etc/crontabs/root && /usr/sbin/crond -f
