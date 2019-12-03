package main

import (
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
    "log"
	"os"
	"text/template"   
)

var tmpl = template.Must(template.ParseGlob("form/*"))

func InitDB(db * sql.DB){
	InitFoodsTable(db)
	log.Println("InitDB Sucesso")
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	return db
}

func main(){
	
	log.Println("Server started on: http://127.0.0.1:5000")
	
	database := dbConn()

	log.Println("database")

	InitDB(database)

	http.HandleFunc("/", ListFoods)
	http.HandleFunc("/new", NewFood)
	http.HandleFunc("/insert", InsertFood)

	http.ListenAndServe(":5000", nil)
	
	defer database.Close()

}



