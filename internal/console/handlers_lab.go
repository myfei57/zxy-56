package console

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"labvent/internal/lab"
)

func registerLabRoutes(router chi.Router, deps Deps) {
	router.Post("/namespaces", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name string `json:"name"`
			Code string `json:"code"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Namespaces.Register(request.Name, request.Code)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/namespaces/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Namespaces.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Get("/namespaces", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Namespaces.List()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/namespaces/{id}/zones", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name string  `json:"name"`
			Area float64 `json:"area"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Zones.Register(chi.URLParam(r, "id"), request.Name, request.Area)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/namespaces/{id}/zones", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Zones.List(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Get("/namespaces/{id}/summary", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Namespaces.Summary(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Post("/labs", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			NamespaceID string `json:"namespace_id"`
			Name        string `json:"name"`
			Code        string `json:"code"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Labs.Register(request.NamespaceID, request.Name, request.Code)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/labs/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Labs.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Get("/labs", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Labs.List(queryParam(r, "namespace_id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/rooms", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			LabID     string `json:"lab_id"`
			Name      string `json:"name"`
			Section   string `json:"section"`
			Cleanroom bool   `json:"cleanroom"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Rooms.Register(request.LabID, request.Name, request.Section, request.Cleanroom)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/rooms/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Rooms.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Get("/rooms", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Rooms.List(queryParam(r, "lab_id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/rooms/{id}/pressure", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Pressure float64 `json:"pressure"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Rooms.UpdatePressure(chi.URLParam(r, "id"), request.Pressure)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Post("/rooms/{id}/baseline", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Value float64 `json:"value"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Baseline.Set(chi.URLParam(r, "id"), request.Value); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"room_id": chi.URLParam(r, "id")})
	})
	router.Get("/baselines/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Baseline.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Post("/partitions", func(w http.ResponseWriter, r *http.Request) {
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
	router.Get("/partitions/{roomID}", func(w http.ResponseWriter, r *http.Request) {
		count, err := deps.Partitions.Count(chi.URLParam(r, "roomID"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"moved_sensors": count})
	})
	router.Get("/mapping/{sensorID}", func(w http.ResponseWriter, r *http.Request) {
		zoneID, err := deps.Mappings.CurrentZone(chi.URLParam(r, "sensorID"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"sensor_id": chi.URLParam(r, "sensorID"), "zone_id": zoneID})
	})
	router.Get("/interlock/{roomID}", func(w http.ResponseWriter, r *http.Request) {
		measured, err := strconv.ParseFloat(queryParam(r, "pressure"), 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Interlock.Check(chi.URLParam(r, "roomID"), measured); err != nil {
			writeError(w, http.StatusLocked, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "released"})
	})
}
