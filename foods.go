package main

import (	
	"net/http"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func InitFoodsTable(db *sql.DB, c *gin.Context) {
/*		if _, err := db.Exec(" DROP TABLE foods"); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error droping database table: %q\n", err))
		}*/
		if _, err := db.Exec(
			" CREATE TABLE IF NOT EXISTS foods ( " +
			" id SERIAL PRIMARY KEY, "+
			" name varchar(255) NOT NULL " +
			" )"); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating database table: %q\n", err))
			return
		}
		rows, err := db.Query("SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_NAME = 'foods'")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading patients columns names: %q\n", err))
			return
		}

		for rows.Next() { 
			var cname string
			if err := rows.Scan(&cname); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning information_schema.COLUMNS: %q\n", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", cname))
		}
		var count int
		c.String(http.StatusOK, fmt.Sprintf("Success: %s\n", rows.Scan(&count)))
		defer rows.Close() 
}
