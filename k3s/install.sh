#!/bin/bash

set -e

mkdir -p /etc/rancher/k3s/
cat <<EOF > /etc/rancher/k3s/config.yaml
write-kubeconfig-mode: "0644"
tls-san:
  - "ethernetes.local"
docker: true
EOF

apt install docker.io
curl -sfL https://get.k3s.io | sh -
snap install helm --classic

# ensure gpu-operator will work
if [[ ! -e "/etc/modprobe.d/blacklist-nouveau.conf" ]]; then
  cat <<EOF > /etc/modprobe.d/blacklist-nouveau.conf
blacklist nouveau
options nouveau modeset=0
EOF
  update-initramfs -u
fi

if ! cat /root/.bashrc | grep KUBECONFIG; then
  cat <<EOF >> /root/.bashrc
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
EOF
fi

if ! cat /home/hayden/.bash_aliases | grep kubectl; then
  cat <<EOF >> /home/hayden/.bash_aliases
alias k='kubectl'
alias oc='kubectl'
alias h='helm'
EOF
fi

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

[[ -d "ethernetes/" ]] || git clone git@github.com/hfuss/etherenetes

pushd ethernetes/k3s/
  kubectl create ns e8s-system || true
  helm upgrade --install \
      argo-cd argo/argo-cd \
      -n e8s-system \
      -f argocd.yaml \
      --atomic \
      --wait
popd