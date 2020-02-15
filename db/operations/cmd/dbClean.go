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
	"github.com/kohrVid/calendar-api/db/operations/dbHelpers"
	"github.com/spf13/cobra"
)

// dbCleanCmd represents the dbClean command
var dbCleanCmd = &cobra.Command{
	Use:   "dbClean",
	Short: "Clean the calendar API database",
	Long:  `This command can be used to delete all rows in the database created for the calendar API based on the environment that it is run in`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig()
		databaseName := conf["database_name"].(string)
		dbHelpers.Clean(conf)

		fmt.Printf("%v database cleaned\n", databaseName)
	},
}

func init() {
	rootCmd.AddCommand(dbCleanCmd)
}
