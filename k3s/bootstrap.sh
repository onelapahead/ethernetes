#!/bin/bash

set -e

# ArgoCD
kubectl create ns e8s-system || true
helm repo add argo https://argoproj.github.io/argo-helm

set +e
if ! helm get all argocd -n e8s-system > /dev/null; then
  helm install \
    argocd argo/argo-cd \
    --set server.extraArgs={--insecure} \
    -n e8s-system \
    -f argocd.yaml \
    --wait

  echo "sleeping for awhile to let apps sync and come up"
  sleep 300
fi

# Secrets
set +e
if ! kubectl get secret -n datadog datadog-creds 2> /dev/null; then
  set -e
  kubectl create secret generic datadog-creds --from-literal api-key="$(lpass show --password datadog-api-key)" --namespace="datadog"
fi
set -e

set +e
if ! kubectl get -n cert-manager secret clouddns-dns01-solver-svc-acct 2> /dev/null; then
  set -e

  gcloud iam service-accounts create dns01-solver --display-name "dns01-solver"
  gcloud projects add-iam-policy-binding ethernetes \
   --member serviceAccount:dns01-solver@ethernetes.iam.gserviceaccount.com \
   --role roles/dns.admin

  mkdir -p ${HOME}/.gcloud/clouddns/dns01/
  gcloud iam service-accounts keys create ${HOME}/.gcloud/clouddns/dns01/key.json \
   --iam-account dns01-solver@ethernetes.iam.gserviceaccount.com
  kubectl create secret generic clouddns-dns01-solver-svc-acct -n cert-manager \
   --from-file=${HOME}/.gcloud/clouddns/dns01/key.json

fi
set -e

helm upgrade --install \
    argocd argo/argo-cd \
    -n e8s-system \
    -f argocd.yaml \
    -f argocd-ingress-tls.yaml \
    --atomic \
    --wait

