#!/bin/bash

hostname="${1:-brx-01a}"
numGPUs="${2:-2}"

cat <<EOF
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: miner-${hostname}
spec:
  destination:
    name: in-cluster
    namespace: ethereum
  project: default
  source:
    repoURL: https://github.com/hfuss/ethernetes
    path: charts/miner
    targetRevision: main
    helm:
      values: |
        nodeSelector:
          kubernetes.io/hostname: ${hostname}

        resources:
          limits:
            nvidia.com/gpu: ${numGPUs}
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
EOF
