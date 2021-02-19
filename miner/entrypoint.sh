#!/bin/bash

# DEBUG
set -ex

miningPools=($MINING_POOLS)
poolArgs=""
for pool in "${miningPools[@]}"; do
  poolArgs="${poolArgs} -P stratum+ssl://${WALLET_ADDRESS}.${HOSTNAME}@${pool}:${STRATUM_TLS_PORT}"
done

# negative port number indicates API is in read-only mode preventing miner from being tampered with
# ensure SIGINTs get mapped to SIGTERMs
exec dumb-init --rewrite 15:2 nsfminer \
  -U \
  ${poolArgs} \
  --api-port "-${API_PORT}" $@