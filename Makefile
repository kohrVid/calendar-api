install: db-create db-migrate

db-create:
	go run db/operations/main.go dbCreate
	go run db/migrations/main.go init

db-migrate:
	go run db/migrations/main.go up

db-migrate-down:
	go run db/migrations/main.go down


db-drop:
	go run db/operations/main.go dbDrop

serve:
	go run main.go

test:
	gocov test -count=1 ./... | gocov report

test-hot-reload:
	#gocov test -count=1 ./... | gocov report

.PHONY: install db-create db-migrate db-drop serve test test-hot-reload
