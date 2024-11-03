postgres:
	docker run --name postgres17 -p 5444:5432 -e POSTGRES_USER=root -e POSTGRES_USER=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres17 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:letmein@localhost:5433/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:letmein@localhost:5433/simple_bank?sslmode=disable" -verbose up
sqlc:
	sqlc generate
test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/dinhtatuanlinh/simplebank Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test mock