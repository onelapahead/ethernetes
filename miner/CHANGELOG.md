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

## [v0.2.1]
### Changed
- API server listens on all interfaces in read-only mode

## [v0.2.0]
### Added
- `API_PORT` env var for configuring the API server port

### Changed
- Enabling API server

### Dependecies
- Bump `https://github.com/ethereum-mining/ethminer` to `v0.19.0`

## [v0.1.0]
### Added
- Miner pre-release
- Miner configured with wallet and mining pool
- Initial workflow implementation

### Dependencies
- Base image `docker.io/nvidia/cuda:11.1-devel-ubuntu18.04`
- Ethereum miner `https://github.com/ethereum-mining/ethminer@v0.18.0`
