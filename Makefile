postgres:
	docker run --name postgres17 -p 5444:5432 -e POSTGRES_USER=root -e POSTGRES_USER=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres17 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:letmein@localhost:5433/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:letmein@localhost:5433/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:letmein@localhost:5433/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:letmein@localhost:5433/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover -coverprofile=.coverage.html ./...

mock:
	mockgen -package mockdb -destination=db/mock/store.go simplebank/db/sqlc Store
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb \
	--go-grpc_opt=paths=source_relative proto/*.proto

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test mock proto