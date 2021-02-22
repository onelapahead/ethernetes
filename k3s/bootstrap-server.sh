#!/bin/bash

set -e

hostname="${1:-ethernetes.brxblx.io}"

# Secrets
set +e
if ! kubectl get secret -n datadog datadog-creds 2> /dev/null; then
  set -e
  echo "Create DataDog secret..."
  kubectl create secret generic datadog-creds --from-literal api-key="$(lpass show --password datadog-api-key)" --namespace="datadog"
fi
set -e

set +e
if ! kubectl get -n cert-manager secret clouddns-dns01-solver-svc-acct 2> /dev/null; then
  set -e

  echo "Create Google Cloud service account and secret..."
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

echo "Updating ArgoCD..."
scp config.yaml root@${hostname}:/etc/rancher/k3s/config.yaml
scp bootstrap/helmchart-argocd.yaml root@${hostname}:/var/lib/rancher/k3s/server/manifests/helmchart-argocd.yaml
