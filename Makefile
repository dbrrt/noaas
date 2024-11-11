# project binaries
GO_BIN := go
# project sources
ROOT_PRJ := .

.PHONY: tidy
tidy:
	$(GO_BIN) mod tidy

.PHONY: dev
dev:
	$(GO_BIN) run $(ROOT_PRJ)

.PHONY: start_nomad
start_nomad:
	nomad agent -dev -bind 0.0.0.0 - network-interface='{{ GetDefaultInterfaces | attr "name" }}' 

.PHONY: start_nomad_compose
start_nomad_compose:
	docker compose up --build

.PHONY: ci
ci:
	$(GO_BIN) test