.PHONY: help build

help:
	@echo "make help                        Show this help message"

	@echo "make run-dev                     Run local dev environment"
	@echo "make run-asset                   Run asset microservice"

	@echo "make test-unit                   Run test unit suite for all microservices"
	@echo "make test-integration            Run test integration suite for all microservices (you need to have running local environment for this)"
	@echo "make test-database               Run test database integration suite for all microservices (you need to have running local environment for this)"
	@echo "make test-end2end                Run end to end test integration suite for main microservices scenario (you need to have running local environment for this)"
	@echo "make test-asset                  Run all the test for asset microservice"
	@echo "make test-blacklist              Run all the test for blacklist microservice"
	@echo "make test-user                   Run all the test for user microservice"

	@echo "make create-assets               Hit REST API to create user assets"
	@echo "make create-chart                Hit REST API to create user chart"
	@echo "make create-assets-with-errors   Hit REST API to create user assets with one existing and one non-existing user"
	@echo "make read-assets                 Hit REST API to read user assets"
	@echo "make update-assets               Hit REST API to update user assets"
	@echo "make delete-assets               Hit REST API to delete user assets"

	@echo "make create-user                 Hit REST API to create user"

	@echo "make build-asset                 Build binary asset stack"
	@echo "make db                          Open local mysql database in terminal"
	@echo "make redis                       Open local redis database in terminal"

run-dev:
	./scripts/install-go.sh && ./scripts/run-dev.sh
run-asset:
	./scripts/install-go.sh && ./scripts/run-asset.sh

test-unit:
	./scripts/install-go.sh && ./scripts/test-unit.sh
test-integration:
	./scripts/install-go.sh && ./scripts/test-integration.sh
test-database:
	./scripts/install-go.sh && ./scripts/test-database.sh
test-end2end:
	./scripts/install-go.sh && ./scripts/test-end2end.sh
test-asset:
	./scripts/install-go.sh && ./scripts/test-asset.sh
test-blacklist:
	./scripts/install-go.sh && ./scripts/test-blacklist.sh
test-user:
	./scripts/install-go.sh && ./scripts/test-user.sh

create-assets:
	./scripts/install-go.sh && ./scripts/create-assets.sh
create-chart:
	./scripts/install-go.sh && ./scripts/create-chart.sh
create-assets-with-errors:
	./scripts/install-go.sh && ./scripts/create-assets-with-errors.sh
read-assets:
	./scripts/install-go.sh && ./scripts/read-assets.sh
update-assets:
	./scripts/install-go.sh && ./scripts/update-assets.sh
delete-assets:
	./scripts/install-go.sh && ./scripts/delete-assets.sh

create-user:
	./scripts/install-go.sh && ./scripts/create-user.sh

build-asset:
	./scripts/install-go.sh && ./scripts/build-asset.sh
db:
	./scripts/db.sh
redis:
	./scripts/redis.sh
