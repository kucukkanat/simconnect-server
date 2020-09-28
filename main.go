package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

// handlers

// ws simvars
func wsSimvars(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		var err error

		// Read
		msg := ""
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Printf("error - hdl wsSimvars - websocket.Message.Receive - %s", err)
			_ = websocket.Message.Send(ws, fmt.Sprintf("ERROR: %v", err))
			_ = ws.Close()
			return
		}
		fmt.Printf("%s\n", msg)

		// parse simvars

		// connect to simvar
		sc, cSimVars, err := scConnectToSimVars([]string{"GENERAL ENG THROTTLE LEVER POSITION"})
		if err != nil {
			_ = websocket.Message.Send(ws, fmt.Sprintf("ERROR: %v", err))
			log.Printf("error - scConnectToSimVars - %s", err)
			return
		}

		defer func() {
			_ = ws.Close()
			sc.Close()
		}()

		for {
			// Write
			select {
			case simvar := <-cSimVars:
				for _, svar := range simvar {
					value, err := svar.GetInt()
					log.Printf("svar name:%s value: %d", svar.Name, value)
					err = websocket.Message.Send(ws, fmt.Sprintf("Throttle: %d%%", value))
					if err != nil {
						if !strings.Contains(err.Error(), "was aborted ") {
							c.Logger().Error(err)
						}
						return
					}
				}
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
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
	e.GET("/simvars", wsSimvars)
	e.Logger.Fatal(e.Start(":1323"))
}
