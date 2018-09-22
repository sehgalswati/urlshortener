FROM busybox:ubuntu-14.04

ADD . /app/urlshortener

WORKDIR /app/urlshortener

ENV http_proxy=""
ENV https_proxy=""

CMD ["./urlshortener"]

EXPOSE 8080
