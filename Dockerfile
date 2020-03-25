FROM golang:1.13-alpine
RUN apk add -q make git && which make git
RUN go get -u \
        golang.org/x/tools/gopls \
        golang.org/x/tools/cmd/guru \
        golang.org/x/tools/cmd/gorename \
        golang.org/x/tools/cmd/goimports \
        golang.org/x/lint/golint
RUN go get -u \
        github.com/go-delve/delve/cmd/dlv \
        github.com/stamblerre/gocode \
        github.com/rogpeppe/godef \
        github.com/mdempsky/gocode \
        github.com/uudashr/gopkgs/cmd/gopkgs \
        github.com/ramya-rao-a/go-outline \
        github.com/acroca/go-symbols
