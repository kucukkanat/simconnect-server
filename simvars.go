package main

import sim "github.com/micmonay/simconnect"

var simvars = map[string]func(args ...interface{}) sim.SimVar{
	"FUEL TOTAL QUANTITY":                 sim.SimVarFuelTotalQuantity,
	"GENERAL ENG RPM":                     sim.SimVarGeneralEngRpm,
	"RECIP ENG FUEL AVAILABLE":            sim.SimVarRecipEngFuelAvailable,
	"GENERAL ENG THROTTLE LEVER POSITION": sim.SimVarGeneralEngThrottleLeverPosition,
}
