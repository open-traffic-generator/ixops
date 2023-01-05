#!/bin/sh

CURRENT_USER=$(whoami)
export PATH=$PATH:/usr/local/go/bin:/home/${CURRENT_USER}/go/bin

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

mk_kind_config() {
    yml="kind: Cluster
        apiVersion: kind.x-k8s.io/v1alpha4
        networking:
          # WARNING: It is _strongly_ recommended that you keep this the default
          # (127.0.0.1) for security reasons. However it is possible to change this.
          # Change to 0.0.0.0 to access kind cluster from outside
          apiServerAddress: 127.0.0.1
          # By default the API server listens on a random open port.
          # You may choose a specific port but probably don't need to in most cases.
          # Using a random port makes it easier to spin up multiple clusters.
          apiServerPort: 6443
        nodes:
          # configure single-node cluster
          - role: control-plane
          # replicate following for multi-node cluster with intended number of worker nodes
          # - role: worker
        "
    echo "$yml" | sed "s/^        //g" | tee ${IXOPS_HOME}/kind.yaml > /dev/null

    for i in $(seq 2 ${KIND_NODE_COUNT})
    do
        echo "  - role: worker" >> ${IXOPS_HOME}/kind.yaml
    done
}

mk_metallb_config() {
    prefix=$(docker network inspect -f '{{.IPAM.Config}}' kind | grep -Eo "[0-9]+\.[0-9]+\.[0-9]+" | tail -n 1)

    yml="apiVersion: metallb.io/v1beta1
        kind: IPAddressPool
        metadata:
          name: kne-pool
          namespace: metallb-system
        spec:
          addresses:
            - ${prefix}.100 - ${prefix}.250

        ---
        apiVersion: metallb.io/v1beta1
        kind: L2Advertisement
        metadata:
          name: kne-l2-adv
          namespace: metallb-system
        spec:
          ipAddressPools:
            - kne-pool
    "
    echo "$yml" | sed "s/^        //g" | tee ${IXOPS_HOME}/metallb.yaml > /dev/null
}

wait_for_pods() {
    for n in $(kubectl get namespaces -o 'jsonpath={.items[*].metadata.name}')
    do
        if [ ! -z "$1" ] && [ "$1" != "$n" ]
        then
            continue
        fi
        for p in $(kubectl get pods -n ${n} -o 'jsonpath={.items[*].metadata.name}')
        do
            if [ ! -z "$2" ] && [ "$2" != "$p" ]
            then
                continue
            fi
            inf "Waiting for pod/${p} in namespace ${n} (timeout=${TIMEOUT_SECONDS}s)..."
            kubectl wait -n ${n} pod/${p} --for condition=ready --timeout=${TIMEOUT_SECONDS}s
        done
    done
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

common_setup() {
    setup_ixops_home \
    && apt_install_pkgs \
    && get_go
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

get_docker() {
    which docker > /dev/null 2>&1 && return
    inf "Installing docker ..."
    sudo apt-get remove docker docker-engine docker.io containerd runc 2> /dev/null

    curl -kfsSL https://download.docker.com/linux/ubuntu/gpg \
        | sudo gpg --batch --yes --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
        | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    sudo apt-get update \
    && sudo apt-get install -y docker-ce docker-ce-cli containerd.io
}

sudo_docker() {
    groups | grep docker > /dev/null 2>&1 && return
    sudo groupadd docker
    sudo usermod -aG docker $CURRENT_USER

    sudo docker version
    inf "Please logout, login again and re-execute previous command" "\n"
    exit 0
}

setup_docker() {
    get_docker \
    && sudo_docker
}

get_kind() {
    which kind > /dev/null 2>&1 && return
    inf "Installing kind ${KIND_VERSION} ..."
    go install sigs.k8s.io/kind@${KIND_VERSION}
}

kind_cluster_exists() {
    kind get clusters | grep kind > /dev/null 2>&1
}

kind_create_cluster() {
    kind_cluster_exists && return
    inf "Creating kind cluster ..."
    mk_kind_config \
    && kind create cluster --config=${IXOPS_HOME}/kind.yaml --wait ${TIMEOUT_SECONDS}s
}

kind_get_kubectl() {
    inf "Copying kubectl from kind cluster to host ..."
    rm -rf kubectl
    docker cp kind-control-plane:/usr/bin/kubectl ./ \
    && sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl \
    && sudo cp -r $HOME/.kube /root/ \
    && rm -rf kubectl
}

kind_get_metallb() {
    inf "Setting up metallb ..."

    kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/${METALLB_VERSION}/config/manifests/metallb-native.yaml \
    && wait_for_pods metallb-system \
    && mk_metallb_config \
    && inf "Applying metallb config map for exposing internal services via public IP addresses ..." \
    && cat ${IXOPS_HOME}/metallb.yaml \
    && kubectl apply -f ${IXOPS_HOME}/metallb.yaml
}

setup_kind() {
    inf "Setting up kind cluster ..."
    setup_docker \
    && get_kind \
    && kind_create_cluster \
    && kind_get_kubectl \
    && kind_get_metallb
}

setup() {
    [ -z "${1}" ] && platform=docker || platform=${1}

    inf "Initiating Setup" "\n"
    common_setup || exit 1

    case $1 in
        docker  )
            setup_kind
        ;;
        kind    )
            setup_kind
        ;;
        gcp    )
            setup_kind
        ;;
        *   )
            err "unsupported image type: ${1}"
        ;;
    esac
}

teardown() {
    [ -z "${1}" ] && platform=docker || platform=${1}

    inf "Initiating Teardown" "\n"
}

newtopo() {
    [ -z "${1}" ] && topo=otg-b2b || topo=${1}

    inf "Initiating Topology Creation" "\n"
}

rmtopo() {
    [ -z "${1}" ] && topo=otg-b2b || topo=${1}

    inf "Initiating Topology Deletion" "\n"
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
