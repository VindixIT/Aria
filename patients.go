package main

import (
	"bytes"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func Write() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		buffer.WriteString("Patients!\n")
		c.String(http.StatusOK, buffer.String())
	}
}