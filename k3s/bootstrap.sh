#!/bin/bash

kubectl create ns e8s-system || true
helm repo add argo https://argoproj.github.io/argo-helm
helm upgrade --install \
    argo-cd argo/argo-cd \
    -n e8s-system \
    -f argocd.yaml \
    --atomic \
    --wait
kubectl create ns longhorn-system || true
kubectl create ns cert-manager || true
kubectl create ns kaleido || true
kubectl create ns datadog || true