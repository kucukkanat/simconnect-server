package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
)

// wsSimvarsWite websocket to write simvars
// client send an array of simvarsWriteCmd

type simvarsWriteCmd struct {
	SimvarName  string
	SimvarIndex int
	Value       string
}

func wsSimvarsWrite(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		//var err error

		// init sc
		sc, err := scConnect()
		if err != nil {
			_ = websocket.Message.Send(ws, fmt.Sprintf("ERROR: %v", err))
			log.Printf("error - hdl wsSimvarsWrite - scConnectToSimVars - %s", err)
			return
		}
		defer func() {
			_ = ws.Close()
			sc.Close()
		}()

		// read from client
		var msg string
		for {
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Printf("error - hdl wsSimvarsWrite - websocket.Message.Receive - %s", err)
				_ = websocket.Message.Send(ws, fmt.Sprintf("ERROR: %v", err))
				return
			}

			if msg == "quit" {
				_ = websocket.Message.Send(ws, "bye")
				return
			}

			// unmarshal msg
			var commands []simvarsWriteCmd
			if err = json.Unmarshal([]byte(msg), &commands); err != nil {
				log.Printf("error - hdl wsSimvarsWrite - websocket.Message.Receive - %s", err)
				_ = websocket.Message.Send(ws, fmt.Sprintf("ERROR bad write command: %v", err))
				return
			}

			for _, cmd := range commands {
				// get data type
				datatype := "float64"

				// parse value according to data type

				//

				newsimvar, exists := simvars[cmd.SimvarName]
				if !exists {
					log.Printf("error - hdl wsSimvarsWrite - no such simvar - %s", cmd.SimvarName)
					_ = websocket.Message.Send(ws, fmt.Sprintf("ERROR no such simvar: %s", cmd.SimvarName))
					return
				}
				nsv := newsimvar()

				switch datatype {
				case "float64":
					nsv.SetFloat64(6000.0)
				default:
					nsv.S

				}

				sc.SetSimObject(nsv)
				time.Sleep(1000 * time.Millisecond)
			}

		}

	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
