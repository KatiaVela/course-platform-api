package payments

import (
	"base/app/models"
	"base/core/module"
	"base/core/router"

	"gorm.io/gorm"
)

type Module struct {
	module.DefaultModule
	DB         *gorm.DB
	Service    *PaymentService
	Controller *PaymentController
}

// Init creates and initializes the Payment module with all dependencies
func Init(deps module.Dependencies) module.Module {
	// Initialize service and controller
	service := NewPaymentService(deps.DB, deps.Emitter, deps.Storage, deps.Logger)
	controller := NewPaymentController(service, deps.Storage)

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
	return m.DB.AutoMigrate(&models.Payment{})
}

func (m *Module) GetModels() []any {
	return []any{
		&models.Payment{},
	}
}
