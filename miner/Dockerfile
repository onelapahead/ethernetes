FROM docker.io/nvidia/cuda:11.1-devel-ubuntu18.04 as base
LABEL maintainer=haydenfuss@gmail.com
LABEL org.opencontainers.image.source=https://github.com/hfuss/ethernetes

RUN apt-get update

FROM base as ethminer

RUN mkdir /ethminer
WORKDIR /ethminer

RUN apt-get install -y curl \
    && curl -sL -o ethminer.tgz https://github.com/ethereum-mining/ethminer/releases/download/v0.19.0-alpha.0/ethminer-0.19.0-alpha.0-cuda-8-linux-x86_64.tar.gz \
    && tar -xzvf ethminer.tgz \
    && mv bin/ethminer /usr/local/bin/ \
    && mkdir /usr/local/bin/kernels \
    && mv bin/kernels/* /usr/local/bin/kernels/ \
    && rm -rf *

FROM base as final

COPY --from=ethminer /usr/local/bin/ethminer /usr/local/bin/ethminer

# TODO see if these env vars do anything still
ENV GPU_FORCE_64BIT_PTR=0 \
    GPU_MAX_HEAP_SIZE=100 \
    GPU_USE_SYNC_OBJECTS=1 \
    GPU_MAX_ALLOC_PERCENT=100 \
    GPU_SINGLE_ALLOC_PERCENT=100 \
    WALLET_ADDRESS=0xf0bEA86827AE84B7a712a4Bc716a15C465be3878 \
    MINE_POOL=us1.ethermine.org \
    STRATUM_TLS_PORT=5555 \
    API_PORT=3333 \
    WORKER_ID=rdu-01a

ENTRYPOINT ["sh", "-c"]
CMD ["ethminer -U -P stratum+tls12://${WALLET_ADDRESS}.${WORKER_ID}@${MINE_POOL}:${STRATUM_TLS_PORT} --report-hashrate --api-bind 0.0.0.0:${API_PORT}"]
EXPOSE ${API_PORT}/tcp