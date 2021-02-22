#!/bin/bash

hostname="${1:-brxblx-01b}"
shift
serverHostname="${1:-ethernetes.brxblx.io}"
shift
extraArgs="$@"
ssh -t hayden@${hostname} "sudo mkdir -p /root/.ssh/"
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

apt install -y docker.io open-iscsi
EOS

if [[ "$extraArgs" == *"--server"* ]]; then
  scp config.yaml root@${hostname}:/etc/rancher/k3s/config.yaml
  scp bootstrap/*.yaml root@${hostname}:/var/lib/rancher/k3s/server/manifests/
fi

k3sup join \
  --host ${hostname} \
  --server-host ${serverHostname} \
  --user root \
  --ssh-key ${HOME}/.ssh/id_ed25519 \
  --k3s-version v1.20.2+k3s1 \
  --k3s-extra-args --docker \
  ${extraArgs}
