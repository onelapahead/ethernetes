#!/bin/bash

set -e

# ArgoCD
kubectl create ns e8s-system || true
helm repo add argo https://argoproj.github.io/argo-helm > /dev/null

set +e
if ! helm get all argocd -n e8s-system > /dev/null; then
  echo "Installing ArgoCD..."
  helm install \
    argocd argo/argo-cd \
    --set server.extraArgs={--insecure} \
    -n e8s-system \
    -f argocd.yaml \
    --wait

  argoServer=$(kubectl get po -n e8s-system -l app.kubernetes.io/component=server -o jsonpath='{ ..metadata.name }')

  kubectl exec \
    -i --tty \
    -n e8s-system ${argoServer} -- /bin/sh -c 'argocd login --insecure 127.0.0.1:8080 --username admin --password ${HOSTNAME} && argocd app sync bootstrap'
  echo "Initial install complete. Copy /etc/rancher/k3s/k3s.yaml to the local machine and re-run bootstrap.sh to add secrets and Ingress secured with TLS."
  exit 0
fi

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

echo "Managing ArgoCD..."
helm upgrade --install \
    argocd argo/argo-cd \
    -n e8s-system \
    -f argocd.yaml \
    -f argocd-ingress-tls.yaml \
    --atomic \
    --wait

