db-up:
	docker compose up --detach --wait db
	$(MAKE) -C backend migrate

db-down:
	docker compose down -v db

db-reset:
	docker compose down -v db
	docker compose up --detach db

format:
	$(MAKE) -C backend format

lint:
	$(MAKE) -C backend lint

integration-test: db-up
	docker compose up --detach --wait backend
	cd integration_tests && uv run main.py; result=$$?; docker compose down -v; exit $$result
