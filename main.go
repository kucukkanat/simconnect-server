package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"

	sim "github.com/micmonay/simconnect"
)

var cSimVar <-chan []sim.SimVar
var sc *sim.EasySimConnect

// handlers

// ws simvars
func wsSimvars(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Write
			select {
			case simvar := <-cSimVar:
				log.Printf("simvar: %v", simvar)
				err := websocket.Message.Send(ws, fmt.Sprintf("%v", simvar))
				if err != nil {
					c.Logger().Error(err)
				}
			default:
				//time.Sleep(1 * time.Second)
				//log.Printf("ras\n")
			}

			/*
				// Read
				msg := ""
				err = websocket.Message.Receive(ws, &msg)
				if err != nil {
					c.Logger().Error(err)
				}
				fmt.Printf("%s\n", msg)

			*/
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	var err error
	cSimVar = make(chan []sim.SimVar)

	sc, err = sim.NewEasySimConnect()
	if err != nil {
		panic(err)
	}

	sc.SetLoggerLevel(sim.LogInfo) // It is better if you the set before connect
	c, err := sc.Connect("MyApp")
	if err != nil {
		panic(err)
	}
	<-c

	/* {
		c, err := sc.Connect("simserv")
		if err != nil {
			if err.Error() == "No connected" {
				log.Print("not connected\n")
				time.Sleep(1 * time.Second)
				continue
			}
			panic(err)
		}
		log.Print("connected\n")
		<-c // Wait connection confirmation
		sc.ShowText("SimServ connected", 1, sim.SIMCONNECT_TEXT_TYPE_PRINT_GREEN)
		break
	}*/

	cSimVar, err = sc.ConnectToSimVar(
		sim.SimVarPlaneAltitude(),
		sim.SimVarPlaneLatitude(sim.UnitDegrees), // You can force the units
		sim.SimVarPlaneLongitude(),
		sim.SimVarIndicatedAltitude(),
		sim.SimVarGeneralEngRpm(1),
		sim.SimVarAutopilotMaster(),
	)
	if err != nil {
		panic(err)
	}

	cSimStatus := sc.ConnectSysEventSim()
	//wait sim start
	for {
		if <-cSimStatus {
			break
		}
	}

	/*for {
		select {
		case result := <-cSimVar:
			for _, simVar := range result {
				var f float64
				var err error
				if strings.Contains(string(simVar.Unit), "String") {
					log.Printf("%s : %#v\n", simVar.Name, simVar.GetString())
				} else if simVar.Unit == "SIMCONNECT_DATA_LATLONALT" {
					data, _ := simVar.GetDataLatLonAlt()
					log.Printf("%s : %#v\n", simVar.Name, data)
				} else if simVar.Unit == "SIMCONNECT_DATA_XYZ" {
					data, _ := simVar.GetDataXYZ()
					log.Printf("%s : %#v\n", simVar.Name, data)
				} else if simVar.Unit == "SIMCONNECT_DATA_WAYPOINT" {
					data, _ := simVar.GetDataWaypoint()
					log.Printf("%s : %#v\n", simVar.Name, data)
				} else {
					f, err = simVar.GetFloat64()
					log.Println(simVar.Name, fmt.Sprintf("%f", f))
				}
				if err != nil {
					log.Println("return error :", err)
				}
			}
		}
	}*/

	// start echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "../clients/testdev")
	e.GET("/simvars", wsSimvars)
	e.Logger.Fatal(e.Start(":1323"))
}
