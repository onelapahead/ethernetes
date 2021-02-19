#!/bin/bash

# negative port number indicates API is in read-only mode preventing miner from being tampered with
# ensure SIGINTs get mapped to SIGTERMs
exec dumb-init --rewrite 15:2 nsfminer \
  -U \
  -P stratum+tls12://${WALLET_ADDRESS}.${HOSTNAME}@${MINE_POOL}:${STRATUM_TLS_PORT} \
  --api-port -${API_PORT} $@