package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"labvent/internal/chem"
)

func registerChemRoutes(router chi.Router, deps Deps) {
	router.Post("/chemicals", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name     string    `json:"name"`
			Class    string    `json:"class"`
			ExpiryAt time.Time `json:"expiry_at"`
			Quantity float64   `json:"quantity"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Chems.Register(request.Name, request.Class, request.ExpiryAt, request.Quantity)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	})
	router.Get("/chemicals/{id}", func(w http.ResponseWriter, r *http.Request) {
		item, err := deps.Chems.Get(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Get("/chemicals", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Chems.List()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Post("/chemicals/{id}/quantity", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Quantity float64 `json:"quantity"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := deps.Chems.UpdateQuantity(chi.URLParam(r, "id"), request.Quantity)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	})
	router.Post("/chemicals/{id}/revise", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			EffectiveUntil time.Time `json:"effective_until"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Expiry.Revise(chi.URLParam(r, "id"), request.EffectiveUntil); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "revised"})
	})
	router.Get("/chemicals/{id}/expiry", func(w http.ResponseWriter, r *http.Request) {
		date, err := deps.Expiry.EffectiveDate(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]time.Time{"effective_until": date})
	})
	router.Post("/chemicals/{id}/issue", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Quantity float64   `json:"quantity"`
			Now      time.Time `json:"now"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		now := request.Now
		if now.IsZero() {
			now = time.Now()
		}
		if err := deps.Issue.Issue(chi.URLParam(r, "id"), request.Quantity, now); err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "issued"})
	})
	router.Post("/rules", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Rules []chem.Rule `json:"rules"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Classify.SetRules(request.Rules); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "rules-updated"})
	})
	router.Get("/rules", func(w http.ResponseWriter, r *http.Request) {
		rules, err := deps.Classify.Rules()
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, rules)
	})
	router.Post("/chemicals/refresh", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.Classify.RefreshAll(); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "refreshed"})
	})
	router.Get("/cabinets/{chemID}", func(w http.ResponseWriter, r *http.Request) {
		required, err := deps.Cabinets.Required(chi.URLParam(r, "chemID"))
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"cabinet_type": required})
	})
	router.Post("/cabinets/{chemID}/verify", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			CabinetType string `json:"cabinet_type"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := deps.Cabinets.Verify(chi.URLParam(r, "chemID"), request.CabinetType); err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "verified"})
	})
}
