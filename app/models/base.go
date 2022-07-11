package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	_ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

var err error
var tracer = otel.Tracer("UserAPI-models")

/*
func init() {
	fmt.Println("initializing...")
	Db, err = sql.Open("sqlite3", config.Config.DbName)
	if err != nil {
		log.Fatalln(err)
	}

	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING NOT NULL UNIQUE,
		name STRING,
		email STRING,
		password STRING,
		created_at DATETIME)`, "users")

	Db.Exec(cmdU)

	log.Println("initializing...DONE!!!!")
}
*/

func init() {
	fmt.Println("Now migration...")
	Db, err = sql.Open("postgres", "host=postgresql.prod.svc.cluster.local port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Println(err)
	}

	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id serial PRIMARY KEY,
		uuid text NOT NULL UNIQUE,
		name text,
		email text,
		password text,
		created_at timestamp)`, "users")

	Db.Exec(cmdU)
	fmt.Println("Now migration...DONE!!")

	log.Println("initializing...DONE!!!!")
}

func createUUID(c *gin.Context) (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(c *gin.Context, plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
