FROM docker.io/nvidia/cuda:11.2.0-devel-ubuntu18.04 as base
LABEL maintainer=haydenfuss@gmail.com
LABEL org.opencontainers.image.source=https://github.com/hfuss/ethernetes

RUN apt-get update

FROM base as ethminer

RUN mkdir /ethminer
WORKDIR /ethminer

RUN apt-get install -y curl \
    && curl -sL -o ethminer.tgz https://github.com/no-fee-ethereum-mining/nsfminer/releases/download/v1.2.4/nsfminer_1.2.4-ubuntu_18.04-cuda_11.2-opencl.tgz \
    && tar -xzvf ethminer.tgz \
    && mv nsfminer /usr/local/bin/ \
    && rm -rf *

FROM base as final

COPY --from=ethminer /usr/local/bin/nsfminer /usr/local/bin/nsfminer
RUN ln -sf /usr/local/bin/nsfminer /usr/local/bin/ethminer

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
CMD ["nsfminer -U -P stratum+tls12://${WALLET_ADDRESS}.${WORKER_ID}@${MINE_POOL}:${STRATUM_TLS_PORT} --report-hashrate --api-port -${API_PORT}"]
EXPOSE ${API_PORT}/tcp
