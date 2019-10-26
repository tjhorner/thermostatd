package main

import (
	"net/http"

	"github.com/tjhorner/thermostatd/thermostat"

	"github.com/chbmuc/lirc"

	"github.com/julienschmidt/httprouter"
	"github.com/tjhorner/thermostatd/api"
)

func main() {
	lircd, err := lirc.Init("/var/run/lirc/lircd")
	if err != nil {
		panic(err)
	}

	therm := thermostat.New(lircd)
	therm.Reset()

	router := httprouter.New()
	api.RouteAll(router, therm)
	http.ListenAndServe(":8080", router)
}
