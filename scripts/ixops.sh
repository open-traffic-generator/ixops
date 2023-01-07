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

### KNE VARIABLES START
KNE_REPO=https://github.com/openconfig/kne.git
KNE_COMMIT=v0.1.7
KNE_HOME="${IXOPS_HOME}/kne"
KNE_CLI="kne"
KNE_VERBOSITY="debug"
### KNE VARIABLES END

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

wait_for_no_namespace() {
    start=$SECONDS
    inf "Waiting for namespace ${1} to be deleted (timeout=${TIMEOUT_SECONDS}s)..."
    while true
    do
        found=""
        for n in $(kubectl get namespaces -o 'jsonpath={.items[*].metadata.name}')
        do
            if [ "$1" = "$n" ]
            then
                found="$n"
                break
            fi
        done

        if [ -z "$found" ]
        then
            return 0
        fi

        elapsed=$(( SECONDS - start ))
        if [ $elapsed -gt ${TIMEOUT_SECONDS} ]
        then
            err "Namespace ${1} not deleted after ${TIMEOUT_SECONDS}s" 1
        fi
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

apt_install_make() {
    make -v > /dev/null 2>&1 && return
    apt_install make
}

apt_install_pkgs() {
    uname -a | grep -i linux > /dev/null 2>&1 || return 0
    inf "Installing apt packages that are not already installed ..."
    apt_install_curl \
    && apt_install_vim \
    && apt_install_git \
    && apt_install_lsb_release \
    && apt_install_gnupg \
    && apt_install_ca_certs \
    && apt_install_make
}

mk_ixops_home() {
    mkdir -p ${IXIA_C_HOME}
}

setup_ixops_home() {
    mk_ixops_home
}

common_setup() {
    check_platform \
    && setup_ixops_home \
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

get_meshnet() {
    inf "Installing meshnet-cni ${MESHNET_REPO} (${MESHNET_COMMIT}) ..."
    rm -rf ${MESHNET_HOME}
    oldpwd=${PWD}
    cd ${IXOPS_HOME}

    applydir=manifests/base
    if [ "${MESHNET_LINK}" = "GRPC" ]
    then
        applydir=manifests/overlays/grpc-link
    fi

    git clone ${MESHNET_REPO} && cd ${MESHNET_HOME} && git checkout ${MESHNET_COMMIT} \
    && cat manifests/base/daemonset.yaml | sed "s#image: networkop/meshnet:latest#image: ${MESHNET_IMAGE}#g" | tee manifests/base/daemonset.yaml.patched > /dev/null \
    && mv manifests/base/daemonset.yaml.patched manifests/base/daemonset.yaml \
    && kubectl apply -k ${applydir} \
    && wait_for_pods meshnet \
    && cd ${oldpwd}
}

get_metrics_server() {
    [ "${METRICS_SERVER_ENABLE}" != true ] && return 0
    inf "Installing metrics server ${METRICS_SERVER_YAML} ..."
    oldpwd=${PWD}
    cd ${IXOPS_HOME} && rm -rf metrics-server.yaml
    curl -kL -o metrics-server.yaml ${METRICS_SERVER_YAML} \
    && cat metrics-server.yaml | sed 's/- args:/- args:\n        - --kubelet-insecure-tls/g' | tee metrics-server.yaml.patched > /dev/null \
    && mv metrics-server.yaml.patched metrics-server.yaml \
    && kubectl apply -f metrics-server.yaml \
    && cd ${oldpwd}
}

rm_metrics_server() {
    [ -f "${IXOPS_HOME}/metrics-server.yaml" ] || return 0
    inf "Removing metrics server ${METRICS_SERVER_YAML} ..."
    kubectl delete -f "${IXOPS_HOME}/metrics-server.yaml" \
    && rm -rf "${IXOPS_HOME}/metrics-server.yaml"
}

get_ixia_c_operator() {
    inf "Installing ixia-c-operator ${IXIA_C_OPERATOR_YAML} ..."
    kubectl apply -f ${IXIA_C_OPERATOR_YAML} \
    && wait_for_pods ixiatg-op-system
}

rm_ixia_c_operator() {
    inf "Removing ixia-c-operator ${IXIA_C_OPERATOR_YAML} ..."
    kubectl delete -f ${IXIA_C_OPERATOR_YAML} \
    && wait_for_no_namespace ixiatg-op-system
}

get_kne() {
    which ${KNE_CLI} > /dev/null 2>&1 && return
    inf "Installing KNE ${KNE_REPO} (${KNE_COMMIT}) ..."
    rm -rf ${KNE_HOME}
    oldpwd=${PWD}
    mk_ixops_home
    cd ${IXOPS_HOME}
    git clone ${KNE_REPO} && cd ${KNE_HOME} && git checkout ${KNE_COMMIT} && make install && cd ${oldpwd}
}

setup_kind_kne() {
    inf "Setting up prerequisites for KNE ..."
    kind_get_metallb \
    && get_meshnet \
    && get_ixia_c_operator \
    && kubectl get pods -A \
    && get_kne
}

setup_kind_k8s() {
    inf "Setting up prerequisites for K8S ..."
    kind_get_metallb \
    && get_meshnet \
    && kubectl get pods -A
}

setup_kind_cluster() {
    inf "Setting up kind cluster ..."
    setup_docker \
    && get_kind \
    && kind_create_cluster \
    && kind_get_kubectl \
    && get_metrics_server \
    && kubectl get pods -A
}

teardown_kind_cluster() {
    inf "Tearing down kind cluster ..."
    kind delete cluster 2> /dev/null
    rm -rf $HOME/.kube
}

setup_kind() {
    [ -z "${1}" ] && cluster=kne || cluster=${1}

    case $1 in
        kne )
            setup_kind_cluster \
            && setup_kind_kne
        ;;
        k8s )
            setup_kind_cluster \
            && setup_kind_k8s
        ;;
        *   )
            err "unsupported cluster type: ${1}" 1
        ;;
    esac
}

teardown_kind() {
    check_platform \
    && teardown_kind_cluster
}

setup_gcp() {
    err "unimplemented" 1
}

mk_ixia_c_config_map() {
    yml='apiVersion: v1
        kind: ConfigMap
        metadata:
          name: ixiatg-release-config
          namespace: ixiatg-op-system
        data:
          versions: |
            {
            "release": "local",
            "images": [
                    {
                        "name": "controller",
                        "path": "ghcr.io/open-traffic-generator/licensed/ixia-c-controller",
                        "tag": "0.0.1-3698"
                    },
                    {
                        "name": "gnmi-server",
                        "path": "ghcr.io/open-traffic-generator/ixia-c-gnmi-server",
                        "tag": "1.10.5"
                    },
                    {
                        "name": "traffic-engine",
                        "path": "ghcr.io/open-traffic-generator/ixia-c-traffic-engine",
                        "tag": "1.6.0.19"
                    },
                    {
                        "name": "protocol-engine",
                        "path": "ghcr.io/open-traffic-generator/licensed/ixia-c-protocol-engine",
                        "tag": "1.00.0.252"
                    }
                ]
            }
        '
    
    echo "${yml}"
}

kne_topo_file() {
    [ -f "${1}" ] && echo ${1} || echo "${IXIA_C_HOME}/${1}.kne.yaml"
}

mk_kne_topo_otg_b2b() {
    echo "name: otg-b2b
        nodes:
          - name: otg
            vendor: KEYSIGHT
            version: ${IXIA_C_RELEASE}
            services:
              8443:
                name: https
                inside: 8443
              40051:
                name: grpc
                inside: 40051
              50051:
                name: gnmi
                inside: 50051
        links:
          - a_node: otg
            a_int: eth1
            z_node: otg
            z_int: eth2
        "
}

mk_kne_topo() {
    inf "Making KNE topology file for ${1} ..."
    case $1 in
        otg-b2b     )
            yml=$(mk_kne_topo_otg_b2b)
        ;;
        otg-dut-otg )
            err "unimplemented" 1
        ;;
        *   )
            if [ -f "${1}" ]
            then
                wrn "${1} is a user provided file"
                return
            fi
            err "unsupported kne topo type: ${1}" 1
        ;;
    esac

    echo "$yml" | sed "s/^        //g" | tee $(kne_topo_file $1) > /dev/null
}

set_ixia_c_config_map() {
    [ -f "${1}" ] && conf=$(echo ${1}) || conf=$(mk_ixia_c_config_map | sed "s/^        //g")
    echo ${mk_ixia_c_config_map} | kubectl apply -f -
}

get_topo_namespace() {
    grep -E "^name" "${1}" | cut -d\  -f2 | sed -e s/\"//g
}

kne_cli() {
    ${KNE_CLI} -v ${KNE_VERBOSITY} $@
}

ctop_kne() {
    inf "Intiating KNE topology creation ..."
    mk_kne_topo "${1}" || exit 1

    topo=$(kne_topo_file $1)
    namespace=$(get_topo_namespace ${topo})
    inf "Using topology ${topo} with namespace ${namespace}"

    # if release is set to local or path to ixia config map has been provided
    if [ "${IXIA_C_RELEASE}" = "local" ] || [ -f "${2}" ]
    then
        set_ixia_c_config_map ${2} || exit 1
    fi

    kne_cli create ${topo} \
    && wait_for_pods ${namespace}
}

dtop_kne() {
    inf "Intiating KNE topology deletion ..."

    topo=$(kne_topo_file $1)
    namespace=$(get_topo_namespace ${topo})
    inf "Using topology ${topo} with namespace ${namespace}"

    kne_cli delete ${topo} \
    && wait_for_no_namespace ${namespace}
}

setup() {
    [ -z "${1}" ] && platform=docker || platform=${1}

    inf "Initiating Setup" "\n"
    common_setup || exit 1

    case $1 in
        docker  )
            setup_docker
        ;;
        kind    )
            setup_kind ${2}
        ;;
        gcp    )
            setup_gcp
        ;;
        *   )
            err "unsupported image type: ${1}" 1
        ;;
    esac

    inf "Please logout, login again (if any binary got installed) !" "\n"
}

teardown() {
    [ -z "${1}" ] && platform=docker || platform=${1}

    inf "Initiating Teardown" "\n"

    case $1 in
        docker  )
            err "unimplemented" 1
        ;;
        kind    )
            teardown_kind
        ;;
        gcp    )
            err "unimplemented" 1
        ;;
        *   )
            err "unsupported image type: ${1}" 1
        ;;
    esac
}

ctop() {
    [ -z "${1}" ] && topo=otg-b2b || topo=${1}

    inf "Initiating Topology Creation" "\n"
    ctop_kne ${topo}
}

dtop() {
    [ -z "${1}" ] && topo=otg-b2b || topo=${1}

    inf "Initiating Topology Deletion" "\n"
    dtop_kne ${topo}
}

help() {
    inf "Welcome to ixops - The easiest way to manage emulated network topologies involving Ixia-C" "\n"
    wrn "Usage ./ixops.sh [subcommand]: "
    echo "setup [docker|kind|gcp] [kne|clab|k8s]    -   Setup prerequisites for a given platform (default=docker)"
    echo "teardown [docker|kind|gcp]                -   Teardown given platform (default=docker)"
    echo "ctop [otg-b2b|otg-dut-otg|<topo-file>]    -   Create emulated topology (default=otg-b2b)"
    echo "dtop [otg-b2b|otg-dut-otg|<topo-file>]    -   Delete emulated topology (default=otg-b2b)"
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
