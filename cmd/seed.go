package cmd

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"watchmen/repository"

	// nolint:staticcheck
	"golang.org/x/crypto/ssh/terminal"

	"watchmen/config"
	"watchmen/database"
)

var seedDatabaseCMD = &cobra.Command{
	Use:   "seed",
	Short: "seed database with data",
	Run: func(cmd *cobra.Command, args []string) {
		seedBaseAPIDB()
	},
}

func init() {
	seedDatabaseCMD.PersistentFlags().BoolVar(&skipSeededDB, "skip-seeded-db", false, "Skip seeding process if any record exists in the db")
}

func seedBaseAPIDB() {
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

	log.Info("Truncating `users`")
	if _, err := db.Exec("TRUNCATE TABLE users;"); err != nil {
		log.Fatalf("error in truncating `users`: `%s`", err)
	}

	log.Info("Creating users")
	users := createMockUsers(10)
	for _, p := range users {
		_, err := db.Exec("INSERT INTO users"+
			" (email, password, fullname, cellphone, created_at)"+
			" VALUES (?, ?, ?, ?, ?)",
			p.Email,
			p.Password,
			p.FullName,
			p.Cellphone,
			p.CreatedAt,
		)
		if err != nil {
			log.Fatal("inserting users failed: ", err)
		}
	}
}

func createMockUsers(n int) []*repository.User {
	// Original: 12345678
	defaultPassword := "$2y$10$LWjMToGy6JNut4FzoGpYC.g3ofR4qF5Ktr.JeluEhj7ykiWw5VOZO"
	var users []*repository.User

	for i := 1; i < n+1; i++ {
		p := &repository.User{
			Email:     fmt.Sprintf("watchman.%02d@mock.com", i),
			Password:  defaultPassword,
			FullName:  fmt.Sprintf("watch%02d man%02d", i, i),
			Cellphone: fmt.Sprintf("+9891200000%02d", i),
			CreatedAt: time.Now(),
		}

		users = append(users, p)
	}

	return users
}
