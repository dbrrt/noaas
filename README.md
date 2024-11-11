# Nomad-as a Service

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

### Start E2E

```bash
make ci
```

## Comments

- busybox has been used instead of nginx, as it's simpler to configure, can receive dynamic port through command line, doesn't require a custom configuration file.
- only one package has been defined for the sake of simplicity, in a larger codebase, I'd split across multiple package to foster better isolation
- remote code execution is dangerous, I've skipped the sanity checks, but that's something required in a secured context.