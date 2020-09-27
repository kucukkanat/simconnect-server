package main

import (
	"log"
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
func scConnectToSimVars(simvarsStr []string) (sc *sim.EasySimConnect, cSimVars <-chan []sim.SimVar, err error) {
	simvarsToConnectTo := make([]sim.SimVar, len(simvarsStr))
	for i, svar := range simvarsStr {
		// check if exist
		v := simvars[svar]
		v.Index = 1
		simvarsToConnectTo[i] = v
	}

	sc, err = scConnect()
	if err != nil {
		return
	}

	//toto := simvars["GENERAL ENG RPM"]

	cSimVars, err = sc.ConnectToSimVarSlice(simvarsToConnectTo)
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
