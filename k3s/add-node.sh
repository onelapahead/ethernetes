#!/bin/bash

hostname="${1:-brxblx-01b}"
shift
serverHostname="${1:-ethernetes.brxblx.io}"
shift
extraArgs="$@"
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
scp bootstrap/*.yaml root@${hostname}:/var/lib/rancher/k3s/server/manifests/

k3sup join \
  --host ${hostname} \
  --server-host ${serverHostname} \
  --user root \
  --ssh-key ${HOME}/.ssh/id_ed25519 \
  ${extraArgs}
