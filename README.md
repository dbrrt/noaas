# Nomad-as a Service

[![noaas-go-ci](https://github.com/dbrrt/noaas/actions/workflows/ci.yml/badge.svg)](https://github.com/dbrrt/noaas/actions/workflows/ci.yml)

## Requirements

- Nomad agent available https://developer.hashicorp.com/nomad/docs/install
- Golang 1.20+ (tested on 1.22)
- Docker (optional, to use Nomad dockerized)

## Get Started

### Start Nomad service

```bash
make start_nomad
```

For the sake of simplicity and portability, nomad has been dockerized and configuration added to a compose file.

```bash
make start_nomad_compose
```
Once nomad agent has been started, the nomad UI can be accessed via http://localhost:4646

## Start NoaaS

```bash
make dev
```


### Start Unit

```bash
make unit
```

### Start E2E

1) Start Nomad in a container

```
make start_nomad_compose
```

2) Execute API tests

```bash
make ci
```

## Notes

- busybox has been used instead of nginx, as it's simpler to configure, can receive dynamic port through command line, doesn't require a custom configuration file.
- only one package has been defined for the sake of simplicity, in a larger codebase, I'd split across multiple package to foster better isolation
- remote code execution is dangerous, I've skipped the sanity checks, but that's something required in a secured context.
- I've created a custom action to test my services with a Nomad cluster in GitHub actions