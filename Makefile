MIGRATIONS_DIR=./migrations

tests:
	go test ./...
run:
	chmod +x ./build_and_run_server.sh
	./build_and_run_server.sh
build_for_prod:
	@echo "Setting environment variables and building for production..."
	@GOOS=linux GOARCH=arm64 go build -o milestone_core_prod
deploy_core_image:
	chmod +x ./deploy_core_image.sh
	./deploy_core_image.sh
deploy_core_prod:
	chmod +x ./deployment/deployment_prod_script.sh
	./deployment/deployment_prod_script.sh
migrate:
	goose -dir migrations postgres "postgres://postgres:password@localhost:5432/milestone_db" up
migrate-prod:
	goose -dir migrations postgres "postgres://postgres:y4WHZrDKcZXSxvq8%40@34.231.242.39:5432/milestone_db" up
migrate-down-to:
	goose -dir migrations postgres "postgres://postgres:password@localhost:5432/milestone_db" down-to $(filter-out $@,$(MAKECMDGOALS))
migration:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Usage: make migration <migration_name>"; \
		exit 1; \
	fi
	cd $(MIGRATIONS_DIR) && goose create $(filter-out $@,$(MAKECMDGOALS)) sql
prepare-test-db:
	chmod +x ./scripts/prepare_test_db.sh
	cd ./scripts && ./prepare_test_db.sh
