package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func registerFanRoutes(router chi.Router, deps Deps) {
	router.Post("/fans", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name   string `json:"name"`
			RowID  string `json:"row_id"`
			HoodID string `json:"hood_id"`
			ZoneID string `json:"zone_id"`
			Role   string `json:"role"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Fans.Register(request.Name, request.RowID, request.HoodID, request.ZoneID, request.Role)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/fans/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Fans.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Get("/fans", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Fans.List(queryParam(r, "row_id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/fans/{id}/start", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Fans.Start(chi.URLParam(r, "id")); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
	})
	router.Post("/fans/{id}/stop", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Fans.Stop(chi.URLParam(r, "id")); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
	})
	router.Post("/fans/{id}/failover", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Fans.Failover(chi.URLParam(r, "id")); err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "failover-done"})
	})
	router.Post("/rows/{rowID}/start", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Fans.StartRow(chi.URLParam(r, "rowID")); err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "row-started"})
	})
	router.Post("/fans/{id}/vibration-trip", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Fans.TripVibration(chi.URLParam(r, "id")); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "tripped"})
	})
	router.Post("/fans/{id}/vibration-clear", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Audit.Reset("fan", chi.URLParam(r, "id"), "vibration reset"); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "cleared"})
	})
	router.Get("/fans/{id}/vibration", func(w http.ResponseWriter, r *http.Request) {
		latched, err := deps.Fans.VibrationLatched(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"latched": latched})
	})
	router.Get("/fans/{id}/state", func(w http.ResponseWriter, r *http.Request) {
		state, err := deps.Fans.StateOf(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"state": state})
	})
	router.Post("/doors/{rowID}/open", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Doors.Open(chi.URLParam(r, "rowID")); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "open"})
	})
	router.Post("/doors/{rowID}/close", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Doors.Close(chi.URLParam(r, "rowID")); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "closed"})
	})
	router.Get("/doors/{rowID}", func(w http.ResponseWriter, r *http.Request) {
		open, err := deps.Doors.IsOpen(chi.URLParam(r, "rowID"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"open": open})
	})
	router.Get("/quota/{labID}/usage", func(w http.ResponseWriter, r *http.Request) {
		used, err := deps.Quota.Usage(chi.URLParam(r, "labID"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]float64{"used_hours": used})
	})
	router.Post("/quotas", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			LabID    string  `json:"lab_id"`
			MaxHours float64 `json:"max_hours"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Quota.Register(request.LabID, request.MaxHours)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/quotas", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Quota.List(queryParam(r, "lab_id"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/quotas/{labID}/reserve", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Hours float64 `json:"hours"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Quota.Reserve(chi.URLParam(r, "labID"), request.Hours); err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "reserved"})
	})
	router.Post("/quotas/{labID}/release", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Hours float64 `json:"hours"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Quota.Release(chi.URLParam(r, "labID"), request.Hours); err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "released"})
	})
	router.Post("/quotas/{labID}/reset", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Quota.Reset(chi.URLParam(r, "labID")); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "reset"})
	})
	router.Get("/fans/row/{rowID}/summary", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Fans.List(chi.URLParam(r, "rowID"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		running := 0
		for _, item := range items {
			if item.Running() {
				running++
			}
		}
		writeJSON(w, http.StatusOK, map[string]int{"total": len(items), "running": running})
	})
	router.Get("/fans/{id}/speed", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Fans.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"speed": item.Speed})
	})
}
