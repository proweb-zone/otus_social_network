package api

import (
	"ms_baskets/internal/app/controller"
	"ms_baskets/pkg/logger"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type SetupBasketsRoute struct {
	logger  *logger.Logger
	gormIns *gorm.DB
}

func SetupBasketsRoutes(logger *logger.Logger, gormIns *gorm.DB) *SetupBasketsRoute {
	return &SetupBasketsRoute{logger: logger, gormIns: gormIns}
}

func (route *SetupBasketsRoute) Init(r chi.Router) {
	basketsController := controller.GetController(route.logger, route.gormIns)

	r.Get("/", basketsController.GetList)
	r.Get("/{id}", basketsController.GetItem)
	r.Put("/create", basketsController.Create)
	r.Post("/update/{id}", basketsController.Update)
	r.Delete("/delete/{id}", basketsController.Delete)
}
