package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"userapi/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

var err error

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

func createUUID(c *gin.Context) (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(c *gin.Context, plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
