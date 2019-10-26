package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tjhorner/thermostatd/thermostat"
)

// APIv1 is v1 of the thermostat API
type APIv1 struct {
	therm *thermostat.Thermostat
}

// prefix implements API.Prefix
func (api *APIv1) prefix() string {
	return "/v1/"
}

func (api *APIv1) getPath(endpoint string) string {
	return fmt.Sprintf("%s%s", api.prefix(), endpoint)
}

// route implements API.route
func (api *APIv1) route(router *httprouter.Router) {
	router.GET(api.getPath("state"), api.jsonify(api.getCurrentState))
	router.PATCH(api.getPath("state"), api.jsonify(api.patchState))
	router.PUT(api.getPath("power"), api.jsonify(api.putPower))
	router.PUT(api.getPath("mode"), api.jsonify(api.putMode))
	router.PUT(api.getPath("temperature"), api.jsonify(api.putTemperature))
	router.PUT(api.getPath("fan_speed"), api.jsonify(api.putFanSpeed))
}

func (api *APIv1) jsonify(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		handle(w, r, p)
	}
}

type apiError struct {
	Error string `json:"error"`
}

func (api *APIv1) apiError(err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	enc := json.NewEncoder(w)
	enc.Encode(apiError{err.Error()})
}

func (api *APIv1) getCurrentState(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enc := json.NewEncoder(w)
	enc.Encode(api.therm.State)
}

func (api *APIv1) patchState(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	draft := api.therm.State

	power := r.Form.Get("powered_on")
	if power != "" {
		draft.PoweredOn = power == "true"
	}

	tmp := r.Form.Get("target_temperature")
	if tmp != "" {
		newTemp, err := strconv.Atoi(tmp)
		if err != nil {
			api.apiError(err, w, r)
			return
		}
		draft.TargetTemperature = newTemp
	}

	mode := thermostat.Mode(r.Form.Get("current_mode"))
	if mode != "" {
		if !mode.IsValid() {
			api.apiError(errors.New("invalid mode"), w, r)
			return
		}
		draft.CurrentMode = mode
	}

	speed := thermostat.FanSpeed(r.Form.Get("fan_speed"))
	if speed != "" {
		if !speed.IsValid() {
			api.apiError(errors.New("invalid fan speed"), w, r)
			return
		}
		draft.FanSpeed = speed
	}

	err = api.therm.SetState(draft)
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	api.getCurrentState(w, r, p)
}

func (api *APIv1) putPower(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	on := r.Form.Get("on")
	err = api.therm.SetPower(on == "true")
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	api.getCurrentState(w, r, p)
}

func (api *APIv1) putMode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	mode := thermostat.Mode(r.Form.Get("mode"))
	if !mode.IsValid() {
		api.apiError(errors.New("invalid mode"), w, r)
		return
	}

	err = api.therm.SetMode(mode)
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	api.getCurrentState(w, r, p)
}

func (api *APIv1) putTemperature(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	tmp := r.Form.Get("temperature")
	newTemp, err := strconv.Atoi(tmp)
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	err = api.therm.SetTargetTemperature(newTemp)
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	api.getCurrentState(w, r, p)
}

func (api *APIv1) putFanSpeed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	speed := thermostat.FanSpeed(r.Form.Get("speed"))
	if !speed.IsValid() {
		api.apiError(errors.New("invalid fan speed"), w, r)
		return
	}

	err = api.therm.SetFanSpeed(speed)
	if err != nil {
		api.apiError(err, w, r)
		return
	}

	api.getCurrentState(w, r, p)
}
