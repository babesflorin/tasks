.PHONY: build test-e2e test-parity test-php test-all clean

# Build the Go API Docker image
build:
	sudo docker compose build go-api

# Go E2E: Go handlers + real MySQL in Docker
test-e2e:
	sudo docker compose -f docker-compose.e2e.yml up -d mysql
	@echo "Waiting for MySQL to be healthy..."
	@timeout=60; while [ $$timeout -gt 0 ]; do \
		if sudo docker compose -f docker-compose.e2e.yml ps mysql | grep -q healthy; then break; fi; \
		sleep 2; timeout=$$((timeout - 2)); \
	done
	sudo docker compose -f docker-compose.e2e.yml run --rm go-e2e-tests
	sudo docker compose -f docker-compose.e2e.yml down -v

# PHP tests: run existing PHP functional+unit tests in Docker
test-php:
	cd .. && sudo docker compose up -d --build
	@echo "Waiting for PHP stack to be ready..."
	@sleep 15
	cd .. && sudo docker compose exec -T php-fpm php bin/phpunit --testsuite FunctionalTests
	cd .. && sudo docker compose exec -T php-fpm php bin/phpunit --testsuite UnitTests

# Parity: both PHP and Go APIs share same MySQL, compare responses
test-parity:
	sudo docker compose -f docker-compose.test.yml up -d --build
	@echo "Waiting for all services to be ready..."
	@sleep 20
	sudo docker compose -f docker-compose.test.yml run --rm test-runner
	sudo docker compose -f docker-compose.test.yml down -v

# Run all test suites
test-all: test-e2e test-php test-parity

# Clean up all Docker resources
clean:
	sudo docker compose -f docker-compose.yml down -v 2>/dev/null || true
	sudo docker compose -f docker-compose.e2e.yml down -v 2>/dev/null || true
	sudo docker compose -f docker-compose.test.yml down -v 2>/dev/null || true
