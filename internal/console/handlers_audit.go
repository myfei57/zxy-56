package console

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func registerAuditRoutes(router chi.Router, deps Deps) {
	router.Get("/audit", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Audit.List(queryParam(r, "subject"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Get("/audit/recent", func(w http.ResponseWriter, r *http.Request) {
		limit := 20
		if raw := queryParam(r, "limit"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			limit = parsed
		}
		items, err := deps.Audit.Recent(limit)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Get("/audit/source/{source}", func(w http.ResponseWriter, r *http.Request) {
		items, err := deps.Audit.BySource(chi.URLParam(r, "source"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
	router.Get("/audit/count/{subject}", func(w http.ResponseWriter, r *http.Request) {
		count, err := deps.Audit.Count(chi.URLParam(r, "subject"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"count": count})
	})
	router.Get("/audit/confirmed/{subject}", func(w http.ResponseWriter, r *http.Request) {
		count, err := deps.Audit.ConfirmedCount(chi.URLParam(r, "subject"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"confirmed": count})
	})
	router.Get("/audit/reset/{subject}", func(w http.ResponseWriter, r *http.Request) {
		count, err := deps.Audit.ResetCount(chi.URLParam(r, "subject"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"resets": count})
	})
}
