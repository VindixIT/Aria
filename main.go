package main

import (
	"net/http"
	"log"
	"database/sql"
	"os"
	
	
)
func main(){
	//http.HandleFunc("/", Index)
	
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	/*if _, err := db.Exec("create table if not exists foods as (with food_groups as (select 'frutas' as name union select 'verduras' union select 'legumes' union select 'carnes') select * from food_groups)"); 
		err != nil {  
		log.Println("Error creating database table: %q", err) 
	 }*/
	
	rows, err :=db.Query("select now()")
	for rows.Next(){
		var teste string
		err:=rows.Scan(&teste)
		if err != nil {
			panic(err)
		}
		log.Println(teste)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	//log.Println("select")	
	http.ListenAndServe(":5000", nil)
	defer rows.Close()
	defer db.Close()
}

func Index(w http.ResponseWriter, r *http.Request){
	log.Println("Sucesso")
}

