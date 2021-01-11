# ethernetes
Learning more about crypto, trading, and deep learning with GPUs and K8s.


## Getting Started

### Requirements

- docker
- CUDA
- nvidia-docker2

```bash
docker build miner/ -t hfuss/miner
nvidia-docker run --restart=always --detach=true --gpus=0 --name=ethminer hfuss/miner 

docker logs ethminer
```
