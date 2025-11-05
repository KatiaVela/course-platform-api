package reviews

import (
	"base/app/models"
	"base/core/module"
	"base/core/router"

	"gorm.io/gorm"
)

type Module struct {
	module.DefaultModule
	DB         *gorm.DB
	Service    *ReviewService
	Controller *ReviewController
}

// Init creates and initializes the Review module with all dependencies
func Init(deps module.Dependencies) module.Module {
	// Initialize service and controller
	service := NewReviewService(deps.DB, deps.Emitter, deps.Storage, deps.Logger)
	controller := NewReviewController(service, deps.Storage)

	// Create module
	mod := &Module{
		DB:         deps.DB,
		Service:    service,
		Controller: controller,
	}

	return mod
}

// Routes registers the module routes
func (m *Module) Routes(router *router.RouterGroup) {
	m.Controller.Routes(router)
}

func (m *Module) Init() error {
	return nil
}

func (m *Module) Migrate() error {
	return m.DB.AutoMigrate(&models.Review{})
}

func (m *Module) GetModels() []any {
	return []any{
		&models.Review{},
	}
}
