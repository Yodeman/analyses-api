sqlc:
	sqlc generate

postgres:
	docker run -p 5000:5432 --name postgres16 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=$DBASE_PASSWORD -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root analyses-api

dropdb:
	docker exec -it postgres16 dropdb analyses-api

migrateup:
	migrate -path dbase/migration -database "postgres://root:${DBASE_PASSWORD}@localhost:5000/analyses-api?sslmode=disable" -verbose up 

migratedown:
	migrate -path dbase/migration -database "postgres://root:${DBASE_PASSWORD}@localhost:5000/analyses-api?sslmode=disable" -verbose down

mock:
	mockgen -package mockdb -destination dbase/mock/store.go github.com/yodeman/analyses-api/dbase/sqlc Querier

test:
	go test -v -cover ./...

