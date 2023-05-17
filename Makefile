.PHONY: help build

help:
	@echo "make help                        Show this help message"
	@echo "make run-dev                     Run local dev environment"
	@echo "make run-unit-tests              Run test unit suite for all microservices"
	@echo "make run-integration-tests       Run test integration suite for web server functionality"
	@echo "make run-blacklist-tests         Run all the test for blacklist microservice"
	@echo "make run-server-tests            Run all the test for server microservice"
	@echo "make create-asset                Hit REST API for create user asset"
	@echo "make read-assets                  Hit REST API for read user assets"
	@echo "make update-asset                Hit REST API for update user asset"
	@echo "make delete-asset                Hit REST API for delete user asset"

	@echo "make build-server                Build binary server stack build HINT: we can go without this, for discussion"
	@echo "make run-server                  Run web server HINT: we can go without this, for discussion"

r:
	./scripts/get-user.sh

run-dev:
	./scripts/install-go.sh && ./scripts/run-dev.sh
run-unit-tests:
	./scripts/install-go.sh && ./scripts/run-unit-tests.sh
run-integration-tests:
	./scripts/install-go.sh && ./scripts/run-integration-tests.sh
run-blacklist-tests:
	./scripts/install-go.sh && ./scripts/run-blacklist-tests.sh
run-server-tests:
	./scripts/install-go.sh && ./scripts/run-server-tests.sh
create-asset:
	./scripts/install-go.sh && ./scripts/create-asset.sh
read-assets:
	./scripts/install-go.sh && ./scripts/read-assets.sh
update-asset:
	./scripts/install-go.sh && ./scripts/update-asset.sh
delete-asset:
	./scripts/install-go.sh && ./scripts/delete-asset.sh


build-server:
	./scripts/install-go.sh && ./scripts/build-server.sh
run-server:
	./scripts/install-go.sh && ./scripts/run-server.sh