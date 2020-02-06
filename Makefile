define CREATEDB
CREATE ROLE calendar_api;\
ALTER ROLE calendar_api WITH SUPERUSER LOGIN CREATEDB;\
CREATE DATABASE calendar_api WITH OWNER calendar_api ENCODING 'UTF8';\

endef

define DROPDB
DROP ROLE calendar_api;\
DROP DATABASE calendar_api;\

endef

install: db-create db-migrate

db-create:
	echo "$(CREATEDB)" | psql -U postgres
	go run db/migrations/main.go init

db-migrate:
	go run db/migrations/main.go up

db-migrate-down:
	go run db/migrations/main.go down


db-drop:
	echo "$(DROPDB)" | psql -U postgres

serve:
	go run main.go

test:
	gocov test -count=1 ./... | gocov report

test-hot-reload:
	#gocov test -count=1 ./... | gocov report

.PHONY: install db-create db-migrate db-drop serve test test-hot-reload
