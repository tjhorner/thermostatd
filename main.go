package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pimvanhespen/go-pi-lcd1602/synchronized"

	lcd1602 "github.com/pimvanhespen/go-pi-lcd1602"

	"github.com/chbmuc/lirc"
	"github.com/julienschmidt/httprouter"
	"github.com/tjhorner/thermostatd/api"
	"github.com/tjhorner/thermostatd/thermostat"
	"github.com/tjhorner/thermostatd/util"
)

func main() {
	lcdi := lcd1602.New(
		21,
		20,
		[]int{19, 13, 26, 6},
		16,
	)

	lcd := synchronized.NewSynchronizedLCD(lcdi)
	lcd.Initialize()
	defer lcd.Close()

	lircd, err := lirc.Init("/var/run/lirc/lircd")
	if err != nil {
		panic(err)
	}

	therm := thermostat.New(lircd, lcd)
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
