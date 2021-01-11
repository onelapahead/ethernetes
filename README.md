# ethernetes
Learning more about crypto, trading, and deep learning with GPUs and K8s.


## Getting Started

### Requirements

- docker
- [CUDA drivers](https://askubuntu.com/questions/1099015/how-to-install-latest-version-of-cuda-on-ubuntu-18-04)
- [nvidia-docker2](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#setting-up-nvidia-container-toolkit)

### Mining


```bash
docker login ghcr.io
nvidia-docker run --restart=always --detach=true --gpus=0 --name=ethminer ghcr.io/hfuss/miner
```

### Logs

```bash
docker logs ethminer
```

### GPU Stats

```bash
nvidia-smi -q -d TEMPERATURE -i 0 -l 10
```
