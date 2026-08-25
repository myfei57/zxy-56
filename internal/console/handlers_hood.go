package console

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"labvent/internal/sash"
	"labvent/internal/valve"
)

func registerHoodRoutes(router chi.Router, deps Deps) {
	router.Post("/hoods", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			LabID  string  `json:"lab_id"`
			RoomID string  `json:"room_id"`
			Name   string  `json:"name"`
			Width  float64 `json:"width"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Hoods.Register(request.LabID, request.RoomID, request.Name, request.Width)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/hoods/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Hoods.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Get("/hoods", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Hoods.List(queryParam(r, "room_id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/hoods/{id}/airflow", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Flow int `json:"flow"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Airflow.Persist(chi.URLParam(r, "id"), request.Flow); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "persisted"})
	})
	router.Get("/hoods/{id}/airflow", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Airflow.Latest(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Get("/hoods/{id}/makeup", func(w http.ResponseWriter, r *http.Request) {
		flow, err := deps.Airflow.Makeup(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"makeup_flow": flow})
	})
	router.Post("/hoods/{id}/velocity", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Threshold float64 `json:"threshold"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Velocity.SetFaceThreshold(chi.URLParam(r, "id"), request.Threshold); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	})
	router.Get("/hoods/{id}/velocity", func(w http.ResponseWriter, r *http.Request) {
		measured, err := strconv.ParseFloat(queryParam(r, "measured"), 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Velocity.Verdict(chi.URLParam(r, "id"), measured); err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	router.Post("/hoods/{id}/pressure", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Value float64 `json:"value"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Pressure.Observe(chi.URLParam(r, "id"), request.Value); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "observed"})
	})
	router.Get("/hoods/{id}/pressure-healthy", func(w http.ResponseWriter, r *http.Request) {
		healthy, err := deps.Pressure.Healthy(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"healthy": healthy})
	})
	router.Post("/hoods/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Status string `json:"status"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Hoods.SetStatus(chi.URLParam(r, "id"), request.Status)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Post("/hoods/{id}/valves", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name string `json:"name"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Valves.Register(chi.URLParam(r, "id"), request.Name)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/hoods/{id}/valves", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Valves.List(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/valves/{id}/stroke", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Position int `json:"position"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Valves.SetStroke(chi.URLParam(r, "id"), request.Position); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "stroked"})
	})
	router.Post("/hoods/{id}/valve-schedule", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			SupplyName  string `json:"supply_name"`
			ExhaustName string `json:"exhaust_name"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		schedule := valve.BuildSchedule(chi.URLParam(r, "id"), request.SupplyName, request.ExhaustName)
		if err := deps.Valves.RunSchedule(schedule); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, schedule)
	})
	router.Post("/sashes", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			HoodID  string `json:"hood_id"`
			FullOpen int    `json:"full_open"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Sashes.Register(request.HoodID, request.FullOpen)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/sashes/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Sashes.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Post("/sashes/{id}/move", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Height int `json:"height"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Moves.Move(chi.URLParam(r, "id"), request.Height); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "moved"})
	})
	router.Post("/sashes/{id}/height", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Height int `json:"height"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Sashes.SetHeight(chi.URLParam(r, "id"), request.Height); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "height-set"})
	})
	router.Post("/sashes/{id}/target", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Threshold float64 `json:"threshold"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Sashes.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		target := sash.HeightForTarget(request.Threshold, item.FullOpen)
		if err := deps.Sashes.SetHeight(item.ID, target); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"height": target})
	})
}
