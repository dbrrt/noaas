# Nomad-as a Service

[![noaas-go-ci](https://github.com/dbrrt/noaas/actions/workflows/ci.yml/badge.svg)](https://github.com/dbrrt/noaas/actions/workflows/ci.yml)

## Requirements

- Nomad agent available https://developer.hashicorp.com/nomad/docs/install
- Golang 1.22+ recommended (tested on 1.22) https://go.dev/doc/install
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

### Start NoaaS API local

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

- busybox has been used instead of nginx, as it's simpler to configure, can receive dynamic port through command line, doesn't require a custom configuration file, then I've been able to focus on code quality over nginx configuration.
- Remote code execution is dangerous, I've skipped the sanity checks, but that's something required in a secured context.
- Created a custom action to test my services with a Nomad cluster in GitHub actions
- Replaced script query parameter type from boolean to string for parsing conveniency by routing validation