package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/tjhorner/thermostatd/util"

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

	token := os.Getenv("THERMOSTATD_TOKEN")
	if token == "" {
		fmt.Println("WARNING: Running without an auth token. Anyone that is able to connect to this device will be able to change your thermostat.")
		http.ListenAndServe(":8080", router)
	} else {
		http.ListenAndServe(":8080", util.NewTokenAuthMiddleware(token, router))
	}
}
