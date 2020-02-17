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

swagger:
	(cd docs; \
	swagger-merger -i api.yaml -o swagger-final.yaml; \
	swagger serve ./swagger-final.yaml)

test:
	ENV=test make db-clean install db-seed -i
	#ENV=test gocov test -count=1 ./... | gocov report
	ENV=test go test -run -count=1 ./...

test-hot-reload:
	./watch_test.sh

.PHONY: install db-create db-migrate db-seed db-clean db-drop serve swagger test test-hot-reload
