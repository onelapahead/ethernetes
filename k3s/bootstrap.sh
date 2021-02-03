#!/bin/bash

sudo -s

curl -sfL https://get.k3s.io | sh -

snap install helm --classic

alias k=kubectl

k create ns e8s-system
helm upgrade --install \
    argo-cd argo/argo-cd \
    -n e8s-system \
    -f argocd.yaml

