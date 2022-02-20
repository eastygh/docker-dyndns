FROM golang:1.17-alpine as builder 
RUN mkdir /build 
COPY src/* /build/ 
WORKDIR /build
RUN go get -d -v ./
RUN go build -o dyndns-api . 

FROM alpine

RUN apk add --no-cache bind9 dnsutils

COPY --from=builder /build/dyndns-api /app/
COPY ./setup.sh /app/setup.sh
RUN chmod +x /app/setup.sh
RUN chmod +x /app/dyndns-api

WORKDIR /app

EXPOSE 53 8080
ENTRYPOINT ["/app/setup.sh"]