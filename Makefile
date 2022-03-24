PASSWORD=password
USERNAME=root
DATABASE=garden
CURRENT=$(PWD)

start_pg_server:
	docker run --name postgres -p 5432:5432 \
	-e POSTGRES_USER=$(USERNAME) \
	-e POSTGRES_PASSWORD=$(PASSWORD) \
	-e POSTGRES_DB=$(DATABASE) \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
	-v $(CURRENT)/temporary/data:/var/lib/postgresql/data/pgdata \
	-d postgres:12-alpine

stop_pg_server:
	docker stop postgres
	docker rm postgres

create_migrate:
	echo migrate create -ext sql -dir db/migration -seq desc_table