package main

import (
	"net/http"
	"log"
	"database/sql"
	"os"
)


func Index(w http.ResponseWriter, r *http.Request){
	log.Println("Sucesso")
}

func InitDB(w http.ResponseWriter, r *http.Request){
	log.Println("InitDB Sucesso")
}



func main(){
	//http.HandleFunc("/", Index)
	
		
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	http.HandleFunc("/InitDB", InitDB)

	http.ListenAndServe(":5000", nil)
	
	defer db.Close()

	
}



