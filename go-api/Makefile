.PHONY: build test-e2e test-parity test-php test-all clean

# Build the Go API Docker image
build:
	sudo docker build -f docker/Dockerfile -t go-api:latest .

# Build the test runner image
build-test:
	sudo docker build -f docker/Dockerfile.test -t go-api-test:latest .

# Go E2E: Go handlers + real MySQL in Docker
# Requires MySQL running on port 3307 (host network)
test-e2e:
	sudo docker compose -f docker-compose.e2e.yml up -d mysql
	@echo "Waiting for MySQL to be healthy..."
	@for i in $$(seq 1 30); do \
		if sudo docker compose -f docker-compose.e2e.yml ps mysql 2>/dev/null | grep -q healthy; then break; fi; \
		sleep 2; \
	done
	sudo docker compose -f docker-compose.e2e.yml run --rm go-e2e-tests
	sudo docker compose -f docker-compose.e2e.yml down -v

# PHP tests: run existing PHP functional + unit tests in Docker (uses SQLite, no MySQL needed)
test-php:
	cd .. && sudo docker build -t builder:latest -f docker/build/Dockerfile .
	sudo docker build -t php-test:latest -f ../docker/test/Dockerfile ..
	sudo docker run --rm -w /var/www/task-list php-test:latest \
		php vendor/bin/simple-phpunit --testsuite FunctionalTests
	sudo docker run --rm -w /var/www/task-list php-test:latest \
		php vendor/bin/simple-phpunit --testsuite UnitTests

# Parity: both PHP and Go APIs share same MySQL, compare responses
# Generates REPORT.md with detailed request/response comparison data.
# See docker-compose.test.yml for the full service definitions.
test-parity:
	sudo docker compose -f docker-compose.test.yml up -d --build
	@echo "Waiting for all services to be ready..."
	@sleep 20
	@mkdir -p $(PWD)/reports
	sudo docker compose -f docker-compose.test.yml run --rm \
		-v $(PWD)/reports:/reports \
		-e PARITY_REPORT_PATH=/reports/REPORT.md \
		test-runner
	@if [ -f $(PWD)/reports/REPORT.md ]; then \
		cp $(PWD)/reports/REPORT.md $(PWD)/REPORT.md; \
		echo ""; echo "=== Parity report written to REPORT.md ==="; \
	fi
	sudo docker compose -f docker-compose.test.yml down -v

# Run all test suites
test-all: test-e2e test-php test-parity

# Clean up all Docker resources
clean:
	sudo docker compose -f docker-compose.yml down -v 2>/dev/null || true
	sudo docker compose -f docker-compose.e2e.yml down -v 2>/dev/null || true
	sudo docker compose -f docker-compose.test.yml down -v 2>/dev/null || true
	sudo docker rm -f parity-mysql parity-go-api parity-php-fpm parity-nginx 2>/dev/null || true
