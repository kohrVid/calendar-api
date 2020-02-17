/*
Copyright © 2020 Jessica Été <kohrVid@zoho.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/kohrVid/calendar-api/config"
	log "github.com/sirupsen/logrus"

	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

// dbCreateCmd represents the dbCreate command
var dbCreateCmd = &cobra.Command{
	Use:   "dbCreate",
	Short: "Create a new database for the calendar API",
	Long:  `This command can be used to create a new database for the calendar API based on the environment that it is run in`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig()
		databaseUser := conf["database_user"].(string)
		databaseName := conf["database_name"].(string)

		createRole := fmt.Sprintf("CREATE ROLE %v", databaseUser)
		alterRole := fmt.Sprintf(
			"ALTER ROLE %v WITH SUPERUSER LOGIN CREATEDB;",
			databaseUser,
		)

		createDB := fmt.Sprintf(
			"CREATE DATABASE %v WITH OWNER %v ENCODING 'UTF8';",
			databaseName,
			databaseUser,
		)

		db := dbHelpers.PostgresDB()
		defer db.Close()

		_, err := db.Exec(createRole)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec(alterRole)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(createDB)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%v database created\n", databaseName)
	},
}

func init() {
	rootCmd.AddCommand(dbCreateCmd)
}
