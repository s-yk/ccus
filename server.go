package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/s-yk/ccus/db"
	"github.com/s-yk/ccus/handler"
)

func main() {
	e := echo.New()

	e.GET("/records", handler.Records)

	if err := db.GetConnection(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Disconnect(); err != nil {
			fmt.Printf("db disconnection error: %s", err)
		}
	}()

	port := os.Getenv("CCUS_PORT")
	if port == "" {
		port = "80"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
