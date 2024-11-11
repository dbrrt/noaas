tidy:
	go mod tidy

dev:
	go run .

start_nomad:
	nomad agent -dev -bind 0.0.0.0 - network-interface='{{ GetDefaultInterfaces | attr "name" }}' 

start_nomad_compose:
	docker compose up --build

