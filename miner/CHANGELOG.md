# CHANGELOG

Inspired from [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## Unreleased
### Added
- ...

### Changed
- ...

### Removed
- ...

### Dependencies
- ...

## [v0.5.0]
### Added
- Custom `entrypoint.sh` is copied to the image as `/usr/local/bin/custom-entrypoint.sh` to allow
  for extra args and `dumb-init`

### Dependencies
- Latest Ethereum miner [`https://github.com/no-fee-ethereum-mining/nsfminer@v1.3.5`](https://github.com/no-fee-ethereum-mining/nsfminer/releases/tag/v1.3.5)


## [v0.4.0]
### Changed
- Decreasing the final image size using the `base` image rather than
  the `devel` image required for building

### Dependencies
- Base image `docker.io/nvidia/cuda:11.2-base-ubuntu18.04`


## [v0.3.0]
### Changed
- Using up-to-date [`nsfminer`](https://github.com/no-fee-ethereum-mining/nsfminer) instead of [`ethminer`](https://github.com/ethereum-mining/ethminer)

### Dependencies
- Base image `docker.io/nvidia/cuda:11.2-devel-ubuntu18.04`
- New Ethereum miner [`https://github.com/no-fee-ethereum-mining/nsfminer@v1.2.4`](https://github.com/no-fee-ethereum-mining/nsfminer/releases/tag/v1.2.4)


## [v0.2.1]
### Changed
- API server listens on all interfaces in read-only mode

## [v0.2.0]
### Added
- `API_PORT` env var for configuring the API server port

### Changed
- Enabling API server

### Dependencies
- Bump `https://github.com/ethereum-mining/ethminer` to `v0.19.0`

## [v0.1.0]
### Added
- Miner pre-release
- Miner configured with wallet and mining pool
- Initial workflow implementation

### Dependencies
- Base image `docker.io/nvidia/cuda:11.1-devel-ubuntu18.04`
- Ethereum miner `https://github.com/ethereum-mining/ethminer@v0.18.0`
