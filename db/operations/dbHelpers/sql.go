package dbHelpers

import (
	"fmt"
	"log"

	"github.com/fatih/structs"
	"github.com/kohrVid/calendar-api/config"
	"github.com/kohrVid/calendar-api/db"
)

func Clean(conf map[string]interface{}) {
	databaseUser := conf["database_user"].(string)

	truncateTables := fmt.Sprintf(`
CREATE OR REPLACE FUNCTION truncate_tables(username IN VARCHAR) RETURNS void AS $$
DECLARE
    statements CURSOR FOR
        SELECT tablename FROM pg_tables
        WHERE tableowner = username
	  AND schemaname = 'public'
	  AND tablename != 'gopg_migrations';
BEGIN
    FOR stmt IN statements LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' RESTART IDENTITY CASCADE;';
    END LOOP;
END;
$$ LANGUAGE plpgsql;
		`)

	cleanDB := fmt.Sprintf(
		"%v SELECT truncate_tables('%v');",
		truncateTables,
		databaseUser,
	)

	db := db.DBConnect(conf)
	defer db.Close()

	_, err := db.Exec(cleanDB)
	if err != nil {
		log.Fatal(err)
	}
}

func Seed(conf map[string]interface{}) {
	users := config.ToMapList(conf["users"])
	if len(users) < 1 {
		log.Fatal("No user to seed")
	}

	var seedDB string

	for _, user := range users {
		s := fmt.Sprintf(`
			  INSERT INTO candidates (first_name, last_name, email)
			    VALUES('%v', '%v', '%v');
			`,
			user["first_name"],
			user["last_name"],
			user["email"],
		)
		seedDB += s
	}

	db := db.DBConnect(conf)
	defer db.Close()

	_, err := db.Exec(seedDB)
	if err != nil {
		log.Fatal(err)
	}
}

func SetSqlColumns(model *structs.Struct, params *structs.Struct) string {
	sql := "SET "
	for _, k := range model.Names() {
		if !params.Field(k).IsZero() {
			col := model.Field(k)

			sql += fmt.Sprintf(
				"%v = ",
				model.Field(k).Tag("json"),
			)

			switch col.Kind().String() {
			case "string":
				sql += fmt.Sprintf(
					"'%v', ",
					params.Field(k).Value(),
				)

			default:
				sql += fmt.Sprintf(
					"%v, ",
					params.Field(k).Value(),
				)
			}

			/*
			  This function exists because go-pg doesn't
			  support the structs library. If it did, the string
			  manipulation could be replaced with the line below.
			  Currently this line is required to ensure that the
			  controller returns the correct JSON object but isn't
			  used in the ORM command.
			*/
			col.Set(params.Field(k).Value())
		}
	}

	return sql[:len(sql)-2]
}
