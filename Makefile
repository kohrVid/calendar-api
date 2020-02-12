install: db-create db-migrate

db-create:
	go run db/operations/main.go dbCreate
	go run db/migrations/main.go init

db-migrate:
	go run db/migrations/main.go up

db-migrate-down:
	go run db/migrations/main.go down

db-seed:
	go run db/operations/main.go dbSeed

db-clean:
	go run db/operations/main.go dbClean

db-drop:
	go run db/operations/main.go dbDrop

serve:
	go run main.go

test:
	ENV=test make db-clean install db-seed -i
	ENV=test gocov test -count=1 ./... | gocov report

test-hot-reload:
	./watch_test.sh

.PHONY: install db-create db-migrate db-seed db-clean db-drop serve test test-hot-reload
