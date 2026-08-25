package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"labvent/internal/lab"
)

func registerGasRoutes(router chi.Router, deps Deps) {
	router.Post("/sensors", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name      string  `json:"name"`
			RoomID    string  `json:"room_id"`
			Threshold float64 `json:"threshold"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Gas.Register(request.Name, request.RoomID, request.Threshold)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Mappings.RegisterSensor(item.ID, queryParam(r, "zone")); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		zoneID := queryParam(r, "zone")
		deps.Router.Cache(item.ID, zoneID)
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/sensors/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Gas.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Get("/sensors", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Gas.List(queryParam(r, "room_id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/sensors/{id}/samples", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			PPM float64 `json:"ppm"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Samples.Record(chi.URLParam(r, "id"), request.PPM); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if err := deps.Alarms.Sample(chi.URLParam(r, "id"), request.PPM); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "sampled"})
	})
	router.Get("/sensors/{id}/sample", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Samples.Latest(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Get("/sensors/{id}/history", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Samples.History(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/sensors/{id}/alarm-confirm", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Alarms.Confirm(chi.URLParam(r, "id")); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "confirmed"})
	})
	router.Get("/sensors/{id}/alarm", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.Alarms.State(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, state)
	})
	router.Get("/sensors/{id}/route", func(w http.ResponseWriter, r *http.Request) {
		zoneID, err := deps.Alarms.Route(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"zone_id": zoneID})
	})
	router.Get("/sensors/{id}/cached-zone", func(w http.ResponseWriter, r *http.Request) {
		zoneID := deps.Router.CachedZone(chi.URLParam(r, "id"))
		writeJSON(w, http.StatusOK, map[string]string{"sensor_id": chi.URLParam(r, "id"), "cached_zone_id": zoneID})
	})
	router.Post("/sensors/{id}/zone", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			ZoneID string `json:"zone_id"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Mappings.RegisterSensor(chi.URLParam(r, "id"), request.ZoneID); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		deps.Router.Cache(chi.URLParam(r, "id"), request.ZoneID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "mapped"})
	})
	router.Post("/partition", func(w http.ResponseWriter, r *http.Request) {
		var request lab.PartitionPlan
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Partitions.Apply(request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "partitioned"})
	})
}
