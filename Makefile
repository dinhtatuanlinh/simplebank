postgres:
	docker run --name postgres17 -p 5444:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret --restart unless-stopped -d postgres:12-alpine
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
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
sqlc:
	sqlc generate
test:
	go test -v -cover -coverprofile=.coverage.html ./...

mock:
	mockgen -package mockdb -destination=db/mock/store.go simplebank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go simplebank/worker TaskDistributor
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb \
	--go-grpc_opt=paths=source_relative proto/*.proto

redis:
	docker run --name redis -p 6379:6379 --restart unless-stopped  -d redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migratedown new_migration sqlc test mock proto redis