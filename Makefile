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

test:
	go test -v ./... -coverprofile=cover.out
	echo "start render coverage report to coverage.html"
	go tool cover --html=cover.out -o coverage.html
	echo "create coverage report at: coverage.html"

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	proto/*.proto

evans:
	evans -r repl --port 8001

.PHONY: start_pg_server stop_pg_server create_migrate test proto
