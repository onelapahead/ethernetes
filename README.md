# ethernetes

Learning more about blockchain and deep learning with GPUs and K8s.

## Getting Started

### Requirements

- Ubuntu 18.04
- docker 20+
- helm 3.5+
- kubectl 1.20+
- argocd 1.8+

### Install

Initializing the first node and ArgoCD consists of the following:

```bash
ssh $(whoami)@brx-01a 'sudo etherenetes/k3s/install.sh'
mkdir -p ${HOME}/.k3s/
ssh -t $(whoami)@brx-01a 'sudo cat /etc/rancher/k3s/k3s.yaml'
```

Copy the output of the final command to a file you can use as your `KUBECONFIG`,
then from there you can add secrets and manage ArgoCD:

```bash
pushd k3s
./bootstrap.sh
popd
```

### Deploying a Miner

First, prepare the miner config for a particular host or node:

```bash
numGPUs=2
hostname=brx-01a
```

#### via ArgoCD

To deploy a new miner to the existing ethernetes cluster you
can do so using ArgoCD and GitOps:

```bash
git checkout main
git pull --rebase origin main
git checkout -b miner-${hostname}

./hack/generate-miner-manifest.sh ${hostname} ${numGPUs} > gitops/deploys/application-${hostname}.yaml
git add gitops/deploys/application-${hostname}.yaml
git commit -m "Deploying a New Miner to ${hostname}"
gh pr create --web --base main
```

Once the PR is merged, the miner will be deployed via [ArgoCD](https://cd.brxblx.io/applications/deploys).

#### via Helm

You can deploy to a particular host on your own cluster using Helm:

```bash
kubectl create ns ethereum
cat <<EOF > miner.yaml

miningPools:
  - us1.ethermine.org
  - us2.ethermine.org
  - eu1.ethermine.org

nodeSelector:
  kubernetes.io/hostname: ${hostname}
  
resources:
  limits:
    nvidia.com/gpu: ${numGPUs}

EOF

helm upgrade --install ethereum-miner charts/miner \
  --wait \
  -n ethereum
  -f miner.yaml

helm test --logs -n ethereum ethereum-miner
```

### Adding a New Node

TODO ...

### Managing Apps via ArgoCD

Application manifests for ArgoCD live underneath the [`gitops/`](gitops/) folder
of this repo. You can access ArgoCD via the CLI:

```bash
argocd login --grpc-web cd.brxblx.io:443
argocd app list
```

[`gitops/bootstrap/`](gitops/bootstrap/) describes the cluster namespaces, controllers, 
and operators needed for ingress, storage and logs, TLS, and leveraging GPUs:

```bash
argocd app get bootstrap
```

[`gitops/deploy/`](gitops/deploy/) describes the namespaces and manifests for deploying the
monitoring stack (i.e. DataDog and Elastic), and the deployments of Ethereum miners,
private blockchain nodes, web apps, and more:

```bash
argocd app get deploys
```

Visit [cd.brxblx.io](https://cd.brxblx.io) to explore and manage the apps via
the UI:

<p align="center">
  <img src="docs/img/argocd.png" width="98%" />
</p>

### Monitoring

#### Elastic

You can see logs from all the miners in the existing cluster [here](https://search.brxblx.io/goto/48ff67e4c824ac8c67314bf8e2293212),

<p align="center">
  <img src="docs/img/logs-search.png" width="98%" />
</p>

Explore the cluster's logs at [search.brxblx.io](https://search.brxblx.io).

#### DataDog

Using [DataDog](https://app.datadoghq.com/dashboard/hes-3t9-pq3/ethereum-miners?from_ts=1613723498904&live=true&to_ts=1613737898904),
it's easy to visualize the health of the miners with respect to the GPUs and system
resources:

<p align="center">
  <img src="docs/img/datadog-dashboard.png" width="98%" />
</p>
