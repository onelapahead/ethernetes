#!/bin/bash

sudo -s

curl -sfL https://get.k3s.io | sh -
snap install helm --classic

cat <<EOF > /etc/modprobe.d/blacklist-nouveau.conf
blacklist nouveau
options nouveau modeset=0
EOF
update-initramfs -u

kubectl create ns e8s-system
helm upgrade --install \
    argo-cd argo/argo-cd \
    -n e8s-system \
    -f argocd.yaml

