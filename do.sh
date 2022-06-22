#!/bin/sh
# author: biplab.mal@keysight.com

ARGS="${@}"

GO_VERSION=1.18
BIN_FILE=bin/ixops
export PATH=$PATH:/usr/local/go/bin:/home/$(logname)/go/bin

cecho() {
    echo "\n\033[1;32m${1}\033[0m\n"
}

# get installers based on host architecture
if [ "$(arch)" = "aarch64" ] || [ "$(arch)" = "arm64" ]
then
    echo "Host architecture is ARM64"
    GO_TARGZ=go${GO_VERSION}.linux-arm64.tar.gz
elif [ "$(arch)" = "x86_64" ]
then
    echo "Host architecture is x86_64"
    GO_TARGZ=go${GO_VERSION}.linux-amd64.tar.gz
else
    echo "Host architecture $(arch) is not supported"
    exit 1
fi


get_go() {
    go version 2> /dev/null && return
    cecho "Installing Go ..."
    # install golang per https://golang.org/doc/install#tarball
    curl -kL https://dl.google.com/go/${GO_TARGZ} | sudo tar -C /usr/local/ -xzf - \
    && echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> $HOME/.profile \
    && . $HOME/.profile \
    && go version
}

get_go_deps() {
    cecho "Installing go dependencies..."
    go mod download
}

rm_go() {
    go version 2> /dev/null || return
    cecho "Uninstalling Go ..."
    sudo rm -rvf /usr/local/go/
}

unit() {
    cecho "Running unit tests..."
    CGO_ENABLED=0 go test -v $(go list ./...)
}

build() {
    cecho "Building ixops..."
    bin_dir=$(dirname ${BIN_FILE})
    bin_file=$(basename ${BIN_FILE})
    mkdir -p ${bin_dir}

    GO_ENABLED=0 go build -v -o ${bin_dir}/
}

case $1 in
    *   )
        # shift positional arguments so that arg 2 becomes arg 1, etc.
        cmd=${1}
        shift 1
        ${cmd} ${@} || cecho "usage: $0 [name of any function in script]"
    ;;
esac



