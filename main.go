package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// handlers

type wsSimvarsResponse struct {
	Ok     bool
	Msg    string
	Simvar string
	Index  int
	Unit   string
	Value  string
}

func main() {
	// display banner
	fmt.Printf(banner, version)
	// start http server
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "../clients/testdev")
	e.GET("/ws/simvarsread", wsSimvarsRead)
	e.Logger.Fatal(e.Start(":1323"))
}
