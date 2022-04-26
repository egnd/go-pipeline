ARG VARIANT="1.18-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN DEBIAN_FRONTEND=noninteractive apt-get update && DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
    make grep graphviz

ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off
ENV GOOS=linux
ENV GOARCH=amd64