package main

import (
	"net/http"
	"log"
	"database/sql"
	"os"
)
func main(){
	http.HandleFunc("/", Index)
	http.ListenAndServe(":5000", nil)
}

func Index(w http.ResponseWriter, r *http.Request){
	log.Println("Sucesso")
}

db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
if err != nil {
	log.Fatalf("Error opening database: %q", err)
}