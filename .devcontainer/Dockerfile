# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.187.0/containers/go/.devcontainer/base.Dockerfile

ARG VARIANT="1.16"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive && \
    apt-get -y install --no-install-recommends openssh-client git make grep && \
    which ssh git make grep
