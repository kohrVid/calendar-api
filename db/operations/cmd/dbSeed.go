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
	"github.com/kohrVid/calendar-api/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// dbSeedCmd represents the dbSeed command
var dbSeedCmd = &cobra.Command{
	Use:   "dbSeed",
	Short: "See the calendar API database",
	Long:  `This command can be used to seed a database created for the calendar API based on the environment that it is run in`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig()
		u1, ok := conf["user1"]
		if !ok {
			log.Fatal("No user to seed")
		}

		u2, ok := conf["user2"]
		if !ok {
			log.Fatal("No user to seed")
		}

		user1 := u1.(map[string]interface{})
		user2 := u2.(map[string]interface{})

		db := db.DBConnect(conf)
		defer db.Close()

		seedDB := fmt.Sprintf(`
		  INSERT INTO candidates (first_name, last_name, email)
		    VALUES('%v', '%v', '%v');
		  INSERT INTO candidates (first_name, last_name, email)
		    VALUES('%v', '%v', '%v');
		`,
			user1["first_name"],
			user1["last_name"],
			user1["email"],
			user2["first_name"],
			user2["last_name"],
			user2["email"],
		)

		_, err := db.Exec(seedDB)
		if err != nil {
			log.Fatal("err")
		}

		fmt.Printf(
			"Seeded %v database\n",
			conf["database_name"].(string),
		)
	},
}

func init() {
	rootCmd.AddCommand(dbSeedCmd)
}
