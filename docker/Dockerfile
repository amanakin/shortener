FROM golang:1.20-alpine

WORKDIR /shortener

COPY .. /shortener

RUN apk update && apk add --no-cache make
RUN make all

CMD ["bin/shortener", "-c", "etc/shortener.yaml"]