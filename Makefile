.PHONY: help build

help:
	@echo "make help                        Show this help message"
	@echo "make build-server                Build binary server stack build"
	@echo "make run-dev                     Run local dev environment"
	@echo "make run-tests                   Run test suite for web server functionality"
	@echo "make run-server                  Run web server"
	@echo "make get-user                    Get user assets"

r:
	./scripts/get-user.sh

build-server:
	./scripts/install-go.sh && ./scripts/build-server.sh

run-dev:
	./scripts/install-go.sh && ./scripts/run-dev.sh

run-tests:
	./scripts/install-go.sh && ./scripts/run-tests.sh

run-server:
	./scripts/install-go.sh && ./scripts/run-server.sh

get-user:
	./scripts/get-user.sh
