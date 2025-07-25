package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
	"os"

	"github.com/PoulDev/lgBlog/internal/blog/model"
	"github.com/microcosm-cc/bluemonday"
)

var DB *sql.DB
var bmp *bluemonday.Policy

func dbExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}

func initSchema(db *sql.DB) {
	schema, err := os.ReadFile("./scripts/schema.sql")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	password, err := RandomPassword()
	if err != nil {
		log.Fatal(err)
	}
	
	_, err = CreateAccount(model.Author{Name: "admin", Picture: "https://i.imgur.com/yOKOBno.png"}, password, model.RoleAdmin)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Login using username: admin, password:", password)

	_, err = NewPost("Titolone", "contenuto contenuto <h1>contenuuuuto</h1>", "descrizione bellona", []int64{1})
	if (err != nil) {
		log.Fatal(err)
	}

}

func LoadDB(filepath string) *sql.DB {
	bmp = bluemonday.UGCPolicy()

	exists := dbExists(filepath)

    db, err := sql.Open("sqlite3", filepath)
    if err != nil {
        log.Fatal(err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

	DB = db

	if !exists {
		log.Println("The database is new, initializing schema...")
		initSchema(db)
	} else {
		log.Println("The database already exists, skipping initialization")
	}

    return db
}
