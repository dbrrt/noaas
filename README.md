# Nomad-as a Service

## Requirements

- Nomad agent available https://developer.hashicorp.com/nomad/docs/install
- Golang 1.20+ (tested on 1.22)
- Docker (optional, to use Nomad dockerized)

## Get Started

### Start Nomad service

For the sake of simplicity and portability, nomad has been dockerized and configuration added to a compose file.

```bash
make start_nomad_compose
```
Once nomad agent has been started, the nomad UI can be accessed via http://localhost:4646

## Start NoaaS

```bash
make dev
```


### Start E2E

```bash
//
```