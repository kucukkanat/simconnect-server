package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
)

// wsSimvarsRead websocket to read simvars
func wsSimvarsRead(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		var err error
		var msg string

		// Read
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Printf("error - hdl wsSimvars - websocket.Message.Receive - %s", err)
			_ = websocket.Message.Send(ws, fmt.Sprintf("ERROR: %v", err))
			_ = ws.Close()
			return
		}
		//fmt.Printf("%s\n", msg)

		reqSimvars := strings.Split(msg, ";")

		// connect to simvar
		sc, cSimVars, err := scConnectToSimVars(reqSimvars)
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
					//value, err := svar.GetDegrees()
					//log.Printf("name:%s index:%d unit:%s value: %d", svar.Name, svar.Index, svar.Unit, value)
					//log.Printf("%v", string(svar.GetDatumType()))
					//log.Printf("Datatype %v", string(svar.GetDatumType()))

					r := wsSimvarsResponse{
						Ok:     true,
						Simvar: svar.Name,
						Index:  svar.Index,
						Unit:   fmt.Sprintf("%s", svar.Unit),
					}

					switch svar.Unit {
					case "String8", "String64", "String", "SIMCONNECT_DATA_LATLONALT", "SIMCONNECT_DATA_XYZ", "SIMCONNECT_DATA_WAYPOINT":
						r.Value = svar.GetString()
					default:
						fValue, err := svar.GetFloat64()
						if err != nil {
							msg := fmt.Sprintf("err - cannot parse svar value to float - %s", err)
							log.Print(msg)
							_ = websocket.Message.Send(ws, wsSimvarsResponse{Ok: false, Msg: msg})
							return
						}
						r.Value = fmt.Sprintf("%f", fValue)
					}

					response, err := json.Marshal(r)
					if err != nil {
						msg := fmt.Sprintf("err - marshal svar failed - %s", err)
						log.Print(msg)
						_ = websocket.Message.Send(ws, wsSimvarsResponse{Ok: false, Msg: msg})
						return
					}

					err = websocket.Message.Send(ws, string(response))
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
