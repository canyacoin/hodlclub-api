FROM ubuntu:18.04

LABEL maintainer "Cam Stuart <cam.asoftware@gmail.com>"

RUN apt-get update && \
    apt-get install -y ca-certificates

COPY web-app-api.linux-x64 /web-app-api

EXPOSE 8080
ENV GIN_MODE release
ENV ETH_NODE_URL "https://inherently-fast-salmon.quiknode.io/568eb2c7-8118-47e6-a518-7549987fef57/z8XJtqvC36I59F-6Zt5egg==/"

ENTRYPOINT ["/web-app-api"]
