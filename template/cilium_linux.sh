#!/bin/bash

set -e

download(){
  curl -LO https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz
  sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
  rm -rf cilium-linux-amd64.tar.gz
}

installCilium(){
  cilium install
}

installCheck(){
  kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/master/examples/kubernetes/connectivity-check/connectivity-check.yaml
}


enableHubble() {
  # Enabling Hubble requires the TCP port 4245 to be open on all nodes running Cilium. T
  # his is required for Relay to operate correctly.
  cilium hubble enable
}

downloadHubbleClient(){
  export HUBBLE_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/hubble/master/stable.txt)
  curl -LO "https://github.com/cilium/hubble/releases/download/$HUBBLE_VERSION/hubble-linux-amd64.tar.gz"
  curl -LO "https://github.com/cilium/hubble/releases/download/$HUBBLE_VERSION/hubble-linux-amd64.tar.gz.sha256sum"
  sha256sum --check hubble-linux-amd64.tar.gz.sha256sum
  tar zxf hubble-linux-amd64.tar.gz

}

main(){
  download
  installCilium
#  enableHubble
}

main