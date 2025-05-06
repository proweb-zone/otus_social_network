package server

import (
	"ms_baskets/internal/api"
	zerolog "ms_baskets/internal/app/middleware/logger"
	"ms_baskets/pkg/logger"
	"ms_baskets/pkg/request_id"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func ConfigureRouting(router *chi.Mux, gormIns *gorm.DB, logger *logger.Logger) *chi.Mux {
	router.Use(request_id.RequestID)
	router.Use(zerolog.New(logger.GetLogger()))

	router.Route("/api", func(r chi.Router) {
		r.Route("/baskets", api.SetupBasketsRoutes(logger, gormIns).Init)
	})

	return router
}
