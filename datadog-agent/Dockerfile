# source: https://github.com/ngi644/datadog_nvml/blob/master/Dockerfile
# forking to lock agent version and reduce external dependecies to only datadog
FROM docker.io/datadog/agent:7.24.0
LABEL maintainer=haydenfuss@gmail.com
LABEL org.opencontainers.image.source=https://github.com/hfuss/ethernetes

ENV NVIDIA_VISIBLE_DEVICES=all

RUN mkdir -p /etc/datadog-agent/checks.d/ && \
    mkdir -p /etc/datadog-agent/conf.d/nvml.d && \
    /opt/datadog-agent/embedded/bin/pip install nvidia-ml-py3==7.352.0

ADD https://raw.githubusercontent.com/ngi644/datadog_nvml/476a4bea631710def768336ebda6f1f59e46e81d/nvml.py /etc/datadog-agent/checks.d/nvml.py
ADD https://raw.githubusercontent.com/ngi644/datadog_nvml/476a4bea631710def768336ebda6f1f59e46e81d/nvml.yaml.default /etc/datadog-agent/conf.d/nvml.d/nvml.yaml

RUN chown -R dd-agent /etc/datadog-agent/checks.d/nvml.py && \
    chown -R dd-agent /etc/datadog-agent/conf.d/nvml.d/nvml.yaml