package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	sim "github.com/micmonay/simconnect"
)

// scConnect connect
func scConnect() (sc *sim.EasySimConnect, err error) {
	sc, err = sim.NewEasySimConnect()
	if err != nil {
		return nil, err
	}

	//sc.SetLoggerLevel(sim.LogInfo) // It is better if you the set before connect
	var c <-chan bool
	for {
		c, err = sc.Connect("simserv")
		if err != nil {
			if err.Error() == "No connected" {
				log.Println("can't connect to msfs. Is it running ?")
				time.Sleep(1 * time.Second)
				continue
			}
			return nil, err
		}
		break
	}
	// wait
	<-c
	return sc, nil
}

// connect to simvars
func scConnectToSimVars(reqSimvars []string) (sc *sim.EasySimConnect, cSimVars <-chan []sim.SimVar, err error) {
	simvarsToConnectTo := make([]sim.SimVar, len(reqSimvars))

	// parse simvars
	// "SIMVAR NAME:INDEX:ARG1,ARG2,...;SIMVAR NAME:INDEX:ARG1,ARG2,..."

	for i, svar := range reqSimvars {
		var args []string
		var index int64
		parts := strings.Split(svar, ":")
		l := len(parts)
		simvar := parts[0]
		if l > 1 {
			// parse index
			if parts[1] != "" {
				index, err = strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
					return nil, nil, fmt.Errorf("scConnectToSimVars failed - bad index %s - %s", parts[1], err)
				}
			}
		}
		if l > 2 {
			args = strings.Split(args[2], ",")
		}
		// check if exist
		if _, exists := simvars[simvar]; !exists {
			return nil, nil, fmt.Errorf("scConnectToSimVars failed - %s is not a valid simvar", simvar)
		}
		v := simvars[simvar](args)
		v.Index = int(index)
		simvarsToConnectTo[i] = v
	}

	sc, err = scConnect()
	if err != nil {
		return
	}
	cSimVars, err = sc.ConnectToSimVar(simvarsToConnectTo...)
	if err != nil {
		return
	}

	cSimStatus := sc.ConnectSysEventSim()
	//wait sim start
	for {
		if <-cSimStatus {
			break
		}
	}
	return
}
