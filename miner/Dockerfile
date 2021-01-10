FROM docker.io/nvidia/cuda:11.1-devel-ubuntu18.04 as base
MAINTAINER haydenfuss@gmail.com

RUN apt-get update

FROM base as build

RUN mkdir /ethminer
WORKDIR /ethminer

RUN apt-get install -y curl \
    && curl -sL -o ethminer.tgz https://github.com/ethereum-mining/ethminer/releases/download/v0.18.0/ethminer-0.18.0-cuda-9-linux-x86_64.tar.gz \
    && tar -xzvf ethminer.tgz \
    && mv bin/ethminer /usr/local/bin/ \
    && rm -rf *

FROM base as final

COPY --from=build /usr/local/bin/ethminer /usr/local/bin/ethminer

ENV GPU_FORCE_64BIT_PTR=0
ENV GPU_MAX_HEAP_SIZE=100
ENV GPU_USE_SYNC_OBJECTS=1
ENV GPU_MAX_ALLOC_PERCENT=100
ENV GPU_SINGLE_ALLOC_PERCENT=100

ENTRYPOINT ["/usr/local/bin/ethminer"]
CMD ["--help"]
