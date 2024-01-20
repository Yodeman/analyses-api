DB_URL := postgres://${DB_USER}:${DB_PASSWORD}@${DB_ADDRESS}/${DB_NAME}?sslmode=disable

sqlc:
	sqlc generate

postgres:
	docker run -p 5000:5432 --network analyses-network --name postgres16 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root ${DB_NAME}

dropdb:
	docker exec -it postgres16 dropdb ${DB_NAME}

migrateup:
	migrate -path dbase/migration -database "${DB_URL}" -verbose up 

migratedown:
	migrate -path dbase/migration -database "${DB_URL}" -verbose down

mock:
	mockgen -package mockdb -destination dbase/mock/store.go github.com/yodeman/analyses-api/dbase/sqlc Querier

test:
	go test -v -cover ./...

