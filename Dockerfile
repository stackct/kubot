# Build image container
FROM golang:1.13.4 AS build
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubot .

# CA Certificates
FROM alpine:latest AS certs
RUN apk --update --no-cache add ca-certificates

# Production image container
FROM scratch
LABEL Description="Deployment bot"

ENV KUBOT_CONFIG=/config.yml
ENV KUBOT_SLACK_TOKEN=

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/kubot /

CMD ["/kubot"]
