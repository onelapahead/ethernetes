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

## [v0.1.0]
### Added
- DataDog agent pre-release
- Custom DataDog agent configured with NVML check for GPU monitoring, forked from [`ngi644/datadog_nvml`](https://github.com/ngi644/datadog_nvml/)
- Initial workflow implementation

### Dependencies
- Base image `docker.io/datadog/agent:7.24.0`
- NVML check `https://github.com/ngi644/datadog_nvml/tree/476a4bea631710def768336ebda6f1f59e46e81d`
- Python library `nvidia-ml-py==v7.352.0`
