package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tjhorner/thermostatd/thermostat"
)

// API is an API
type API interface {
	prefix() string
	route(router *httprouter.Router)
}

// Route will route a specific API
func Route(api API, router *httprouter.Router) {
	api.route(router)
}

// RouteAll routes all available API versions
func RouteAll(router *httprouter.Router, therm *thermostat.Thermostat) {
	apis := []API{
		&APIv1{therm},
	}

	for _, api := range apis {
		Route(api, router)
	}
}
