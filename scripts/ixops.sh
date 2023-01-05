#!/bin/sh

export PATH=$PATH:/usr/local/go/bin:/home/$(logname)/go/bin

# source path for current session
. $HOME/.profile

GO_VERSION=1.19
PROTOC_VERSION=3.20.1
KUBERNETES_VERSION=v1.23.6
METALLB_VERSION=v0.13.6
CONTAINERLAB_VERSION=0.26.2

IXOPS_HOME="${HOME}/.ixops"
TIMEOUT_SECONDS=300

### IXIA-C VARIABLES START ###
IXIA_C_RELEASE="0.0.1-3698"             # set to "local" for private builds
IXIA_C_OPERATOR_RELEASE="0.3.1"
IXIA_C_OPERATOR_YAML="https://github.com/open-traffic-generator/ixia-c-operator/releases/download/v${IXIA_C_OPERATOR_RELEASE}/ixiatg-operator.yaml"
IXIA_C_HOME="${IXOPS_HOME}/ixia-c"
### IXIA-C VARIABLES END ###

### MESHNET VARIABLES START
MESHNET_COMMIT=4674905
MESHNET_REPO=https://github.com/networkop/meshnet-cni
MESHNET_HOME="${IXOPS_HOME}/meshnet-cni"
MESHNET_IMAGE="networkop/meshnet\:v0.3.0"
MESHNET_LINK=VXLAN                      # set to "GRPC" for gRPC link
### MESHNET VARIABLES END

### METRICS SERVER VARIABLES START
METRICS_SERVER_ENABLE=false
METRICS_SERVER_YAML="https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.6.2/components.yaml"
### METRICS SERVER VARIABLES END

### KIND VARIABLES START
KIND_VERSION=v0.17.0
KIND_NODE_COUNT=1
### KIND VARIABLES END


inf() {
    echo "${2}\033[1;32m${1}\033[0m${2}"
}

wrn() {
    echo "${2}\033[1;33m${1}\033[0m${2}"
}

err() {
    echo "\n\033[1;31m${1}\033[0m\n"
    [ ! -z ${2} ] && exit ${2}
}

check_platform() {
    grep "Ubuntu" /etc/os-release > /dev/null 2>&1 || err "This operation is only supported on Ubuntu" 1
}

apt_update() {
    if [ "${APT_UPDATE}" = "true" ]
    then
        sudo apt-get update
        APT_GET_UPDATE=false
    fi
}

apt_install() {
    inf "Installing ${1} ..."
    apt_update \
    && sudo apt-get install -y --no-install-recommends ${1}
}

apt_install_curl() {
    curl --version > /dev/null 2>&1 && return
    apt_install curl
}

apt_install_vim() {
    dpkg -s vim > /dev/null 2>&1 && return
    apt_install vim
}

apt_install_git() {
    git version > /dev/null 2>&1 && return
    apt_install git
}

apt_install_lsb_release() {
    lsb_release -v > /dev/null 2>&1 && return
    apt_install lsb_release
}

apt_install_gnupg() {
    gpg -k > /dev/null 2>&1 && return
    apt_install gnupg
}

apt_install_ca_certs() {
    dpkg -s ca-certificates > /dev/null 2>&1 && return
    apt_install ca-certificates
}

apt_install_pkgs() {
    uname -a | grep -i linux > /dev/null 2>&1 || return 0
    inf "Installing apt packages that are not already installed ..."
    apt_install_curl \
    && apt_install_vim \
    && apt_install_git \
    && apt_install_lsb_release \
    && apt_install_gnupg \
    && apt_install_ca_certs
}

mk_ixops_home() {
    mkdir -p ${IXIA_C_HOME}
}

setup_ixops_home() {
    mk_ixops_home
}

get_go() {
    which go > /dev/null 2>&1 && return
    inf "Installing Go ${GO_VERSION} ..."
    # install golang per https://golang.org/doc/install#tarball
    curl -kL https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz | sudo tar -C /usr/local/ -xzf - \
    && echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> $HOME/.profile \
    && . $HOME/.profile \
    && go version
}

common_setup() {
    setup_ixops_home \
    && apt_install_pkgs \
    && get_go
}

setup() {
    [ -z "${1}" ] && platform=docker || platform=${1}

    wrn "+++++ setup +++++" "\n"
    common_setup
}

teardown() {
    [ -z "${1}" ] && platform=docker || platform=${1}

    wrn "+++++ teardown +++++" "\n"
}

newtopo() {
    [ -z "${1}" ] && topo=otg-b2b || topo=${1}

    wrn "+++++ newtopo +++++" "\n"
}

rmtopo() {
    [ -z "${1}" ] && topo=otg-b2b || topo=${1}

    wrn "+++++ rmtopo +++++" "\n"
}

help() {
    inf "Welcome to ixops - The easiest way to manage emulated network topologies involving Ixia-C" "\n"
    wrn "Usage ./ixops.sh [subcommand]: "
    echo "setup [docker|kind|gcp] [kne|clab|k8s]    -   Setup prerequisites for a given platform (default=docker)"
    echo "teardown [docker|kind|gcp]                -   Teardown given platform (default=docker)"
    echo "newtopo [otg-b2b|otg-dut-otg|<topo-file>] -   Create emulated topology (default=otg-b2b)"
    echo "rmtopo [otg-b2b|otg-dut-otg|<topo-file>]  -   Delete emulated topology (default=otg-b2b)"
    echo "\n"
    
    wrn "Notes:"
    echo "  - To execute functions directly from the script (not listed above), execute: ./ixops.sh [function-name]"
    echo "  - Always use the latest version of this script: curl -kLO https://raw.githubusercontent.com/open-traffic-generator/ixops/main/scripts/ixops.sh"
    echo "\n"
}

case $1 in
    ""  )
        err "usage: $0 [name of any function in script]" 1
    ;;
    *   )
        # shift positional arguments so that arg 2 becomes arg 1, etc.
        cmd=${1}
        shift 1
        ${cmd} ${@} || err "failed executing './ixops.sh ${cmd} ${@}'"
    ;;
esac
