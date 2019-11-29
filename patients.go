package main

import (	
	"net/http"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func InitPatientsTable(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		db.Exec(" DROP TABLE patients")
		if _, err := db.Exec(
			" CREATE TABLE IF NOT EXISTS patients ( " +
			" id smallint, "+
			" weight decimal(3,3), "+
			" age smallint" +
			" )"); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating database table: %q", err))
			return
		}
		rows, err := db.Query("SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_NAME = 'patients'")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading patients columns names: %q", err))
			return
		}

		for rows.Next() { 
			var cname string
			if err := rows.Scan(&cname); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning information_schema.COLUMNS: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", cname))
		}
		var count int
		c.String(http.StatusOK, fmt.Sprintf("SUCESSO: %s\n", rows.Scan(&count)))
		defer rows.Close()
	}
}
