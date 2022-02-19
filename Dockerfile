FROM golang:1.17-alpine as builder 
RUN mkdir /build 
COPY src/* /build/ 
WORKDIR /build
RUN go get -d -v ./
RUN go build -o dyndns-api . 

FROM debian:stretch
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -qq && \
    apt-get install -q -y bind9 dnsutils -qq --no-install-recommends && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/dyndns-api /app/
COPY ./setup.sh /app/setup.sh
RUN chmod +x /app/setup.sh

WORKDIR /app

EXPOSE 53 8080
ENTRYPOINT ["/app/setup.sh"]