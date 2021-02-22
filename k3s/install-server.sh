#!/bin/bash

hostname="${1:-ethernetes.brxblx.io}"
ssh -t hayden@${hostname} "sudo cp /home/hayden/.ssh/authorized_keys /root/.ssh/"

ssh root@${hostname} <<EOS
set -e

mkdir -p /etc/rancher/k3s/
mkdir -p /var/lib/rancher/k3s/server/manifests

# ensure gpu-operator will work
if [[ ! -e "/etc/modprobe.d/blacklist-nouveau.conf" ]]; then
  cat <<EOF > /etc/modprobe.d/blacklist-nouveau.conf
blacklist nouveau
options nouveau modeset=0
EOF
  update-initramfs -u
fi

apt install docker.io
EOS

scp config.yaml root@${hostname}:/etc/rancher/k3s/config.yaml
scp install/*.yaml root@${hostname}:/var/lib/rancher/k3s/server/manifests/

mkdir -p ${HOME}/.k3s/
touch ${HOME}/.k3s/config.yaml
k3sup install \
  --skip-install \
  --host ${hostname} \
  --user root \
  --ssh-key ${HOME}/.ssh/id_ed25519 \
  --merge \
  --local-path ${HOME}/.k3s/config.yaml

sleep 15

export KUBECONFIG=${HOME}/.k3s/config.yaml
argoServer=$(kubectl get po -n e8s-system -l app.kubernetes.io/component=server -o jsonpath='{ ..metadata.name }')

kubectl exec \
  -i --tty \
  -n e8s-system ${argoServer} -- /bin/sh -c 'argocd login --insecure 127.0.0.1:8080 --username admin --password ${HOSTNAME} && argocd app sync bootstrap'
