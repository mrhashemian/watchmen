package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// nolint:staticcheck
	"golang.org/x/crypto/ssh/terminal"

	"watchmen/config"
	"watchmen/database"
)

var migrateDatabaseCMD = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrateBaseAPIDB()
	},
}

func migrateBaseAPIDB() {
	if !(host == "" && port == "" && db == "" && user == "") {
		fmt.Print("Password: ")
		terminalOutput, err := terminal.ReadPassword(0)
		if err != nil {
			log.Fatalf("there is problem on reading password from terminal: %s", err)
		}
		config.C.BaseAPIDatabase.Password = string(terminalOutput)
	}

	db := database.InitBaseAPIDB()
	defer database.CloseDB(db)

	drop := `
DROP TABLE IF EXISTS users;
`
	_, err := db.MustExec(drop).RowsAffected()
	if err != nil {
		log.Fatal("base-api migration (drop tables) failed: ", err)
	}

	schema := `
CREATE TABLE users (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  email varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  password varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  fullname varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  cellphone varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  created_at timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  updated_at timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (id)
);
`
	_, err = db.MustExec(schema).RowsAffected()
	if err != nil {
		log.Fatal("base-api migration (create tables) failed: ", err)
	}
}
