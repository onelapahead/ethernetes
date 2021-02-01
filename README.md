# ethernetes
Learning more about crypto, trading, and deep learning with GPUs and K8s.


## Getting Started

### Requirements

- Ubuntu 18.04
- docker
- [CUDA drivers](https://askubuntu.com/questions/1099015/how-to-install-latest-version-of-cuda-on-ubuntu-18-04)
- [nvidia-docker2](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#setting-up-nvidia-container-toolkit)

### Docker

```bash
sudo docker login ghcr.io
sudo systemctl daemon-reload
sudo systemctl restart docker
```

### Mining

```bash
sudo docker network create eth
sudo nvidia-docker run --network eth -e WORKER_ID=$(hostname) -p 127.0.0.1:3333:3333/tcp --restart=always --detach=true --gpus=0 --name=ethminer ghcr.io/hfuss/miner:latest
```

To test the API server:

```bash
sudo docker run --network eth -it ghcr.io/hfuss/ethminer-exporter:latest client --hostname ethminer getstatdetail
{
  "connection": {
    "connected": true,
    "switches": 1,
    "uri": "stratum+tls12://0xf0bEA86827AE84B7a712a4Bc716a15C465be3878.*****@us1.ethermine.org:5555"
  },
  "devices": [
    {
      "_index": 0,
      "_mode": "CUDA",
      "hardware": {
        "name": "GeForce RTX 3060 Ti 7.79 GB",
        "pci": "09:00.0",
        "sensors": [
          0,
          0,
          0
        ],
        "type": "GPU"
      },
      "mining": {
        "hashrate": "0x031316f0",
        "pause_reason": null,
        "paused": false,
        "shares": [
          1519,
          0,
          0,
          3
        ]
      }
    },
    {
      "_index": 1,
      "_mode": "CUDA",
      "hardware": {
        "name": "GeForce GTX 1070 Ti 7.79 GB",
        "pci": "0a:00.0",
        "sensors": [
          0,
          0,
          0
        ],
        "type": "GPU"
      },
      "mining": {
        "hashrate": "0x019aacb8",
        "pause_reason": null,
        "paused": false,
        "shares": [
          767,
          0,
          0,
          644
        ]
      }
    }
  ],
  "host": {
    "name": "****",
    "runtime": 118882,
    "version": "nsfminer-1.2.4"
  },
  "mining": {
    "difficulty": 3999938964,
    "epoch": 392,
    "epoch_changes": 2,
    "hashrate": "0x04adc3a8",
    "shares": [
      2286,
      0,
      0,
      3
    ]
  },
  "monitors": null
}
```

### Real-Time Logs

```bash
sudo docker logs ethminer --follow --since 10s
```

<!-- TODO kibana searches via ECK -->

### Monitoring via DataDog

```bash
sudo nvidia-docker run -d --gpus=all \
  --restart always \
  --name datadog-agent \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /proc/:/host/proc/:ro \
  -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
  -v /opt/datadog-agent-conf.d:/conf.d:ro \
  -v /opt/datadog-agent-checks.d:/checks.d:ro \
  -e DD_API_KEY=${DD_API_KEY} \
  -e DD_SITE=datadoghq.com \
  ghcr.io/hfuss/datadog-agent:latest
```

<p text="align">
  <img src="docs/img/dashboard.png" width="70%" />
</p>