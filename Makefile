db-up:
	docker-compose up -d db

db-down:
	docker-compose down -v db

db-reset:
	docker-compose down -v db
	docker-compose up -d db

format:
	$(MAKE) -C backend format

lint:
	$(MAKE) -C backend lint
