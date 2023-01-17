package database

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/qustavo/sqlhooks/v2/hooks/othooks"
	log "github.com/sirupsen/logrus"

	"watchmen/config"
)

// NewMySQLConnection create connection to a MySQL/MariaDB server with passed arguments
// and returns an SQLx DB struct.
func NewMySQLConnection(
	name,
	baseDSN string,
	retry int,
	maxOpenConn int,
	maxIdleConn int,
	retryTimeout time.Duration,
	timeout time.Duration) *sqlx.DB {
	var db *sqlx.DB
	var err error
	var id int
	counter := 0

	sql.Register("tracing-mysql-"+name, sqlhooks.Wrap(&mysql.MySQLDriver{}, othooks.New(opentracing.GlobalTracer())))
	db, err = sqlx.Open("tracing-mysql-"+name, baseDSN)
	if err != nil {
		log.Fatalf("Cannot open database %s: %s", baseDSN, err)
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(timeout)

	counter = 0
	ticker := time.NewTicker(retryTimeout)
	for ; true; <-ticker.C {
		counter++
		err := db.QueryRow("SELECT connection_id()").Scan(&id)
		if err == nil {
			break
		}

		log.Errorf("Cannot connect to database %s: %s", baseDSN, err)
		if counter >= retry {
			log.Fatalf("Cannot connect to database %s after %d retries: %s", baseDSN, counter, err)
		}
	}

	ticker.Stop()

	return db
}

// InitBaseAPIDB create a connection to the BaseAPI DB
func InitBaseAPIDB() *sqlx.DB {
	db := NewMySQLConnection(
		"base-api",
		config.C.BaseAPIDatabase.String(),
		config.C.BaseAPIDatabase.DialRetry,
		config.C.BaseAPIDatabase.MaxConn,
		config.C.BaseAPIDatabase.IdleConn,
		config.C.BaseAPIDatabase.DialTimeout,
		config.C.BaseAPIDatabase.Timeout,
	)

	log.Info("App connected to the base-api database")

	return db
}

// CloseDB closes passed DB
func CloseDB(db *sqlx.DB) {
	err := db.Close()
	if err != nil {
		log.Error(err)
	}
}
