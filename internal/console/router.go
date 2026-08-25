package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func buildRouter(deps Deps) http.Handler {
	router := chi.NewRouter()
	router.Use(requestLogger)
	router.Use(recoverPanic)
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	registerLabRoutes(router, deps)
	registerHoodRoutes(router, deps)
	registerFanRoutes(router, deps)
	registerGasRoutes(router, deps)
	registerChemRoutes(router, deps)
	registerAuditRoutes(router, deps)
	return router
}
