# Build image container
FROM golang:1.13.4 AS build
ARG K8S_VERSION=v1.15.2
ARG HELM_VERSION=2.9.0
WORKDIR /build
COPY . /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubot .

RUN curl https://storage.googleapis.com/kubernetes-helm/helm-v${HELM_VERSION}-linux-amd64.tar.gz -o /tmp/helm.tar.gz \
	&& tar -zxvf /tmp/helm.tar.gz -C /tmp \
	&& mv /tmp/linux-amd64/helm /usr/local/bin/helm \
	&& rm -rf /tmp/linux-amd64 \
	&& rm /tmp/helm.tar.gz

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl \
	&& chmod +x ./kubectl \
	&& mv ./kubectl /usr/local/bin/kubectl


# CA Certificates
FROM alpine:latest AS alpine-base
RUN apk --update --no-cache add ca-certificates

# Production image container
FROM alpine:latest
LABEL Description="Deployment bot"

ENV KUBOT_CONFIG=/conf/kubot.yml

RUN apk --update --no-cache add git
COPY --from=alpine-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/kubot /opt/kubot/
COPY --from=build /usr/local/bin/helm /usr/local/bin
COPY --from=build /usr/local/bin/kubectl /usr/local/bin

VOLUME ["/opt/kubot/log"]

CMD ["/opt/kubot/kubot"]
