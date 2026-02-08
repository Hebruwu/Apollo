db-up:
	docker compose up --detach db
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
