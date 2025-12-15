package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	handlers "github.com/kuahbanyak/go-crud/internal/adapters/handlers/http"
	"github.com/kuahbanyak/go-crud/internal/adapters/handlers/http/middleware"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/config"
	"github.com/kuahbanyak/go-crud/internal/shared/constants"
)

type HTTPServer struct {
	server                 *http.Server
	router                 *mux.Router
	userHandler            *handlers.UserHandler
	productHandler         *handlers.ProductHandler
	waitingListHandler     *handlers.WaitingListHandler
	settingHandler         *handlers.SettingHandler
	vehicleHandler         *handlers.VehicleHandler
	maintenanceItemHandler *handlers.MaintenanceItemHandler
	healthHandler          *handlers.HealthHandler
	versionHandler         *handlers.VersionHandler
	invoiceHandler         *handlers.InvoiceHandler
	analyticsHandler       *handlers.AnalyticsHandler
	roleHandler            *handlers.RoleHandler
}

func NewHTTPServer(
	cfg *config.Config,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
	waitingListHandler *handlers.WaitingListHandler,
	settingHandler *handlers.SettingHandler,
	vehicleHandler *handlers.VehicleHandler,
	maintenanceItemHandler *handlers.MaintenanceItemHandler,
	healthHandler *handlers.HealthHandler,
	versionHandler *handlers.VersionHandler,
	invoiceHandler *handlers.InvoiceHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	roleHandler *handlers.RoleHandler,
) *HTTPServer {
	router := mux.NewRouter()

	// Global middleware (order matters!)
	router.Use(middleware.RequestID)           // 1. Add request ID first for tracing
	router.Use(middleware.CORS)                // 2. Handle CORS
	router.Use(middleware.ValidateRequestSize) // 3. Limit request size
	router.Use(middleware.Logging)             // 4. Log requests (will include request ID)

	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	router.Use(middleware.RateLimit(rateLimiter))

	// Handle preflight OPTIONS requests for all routes
	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(constants.DefaultReadTimeout) * time.Second,
		WriteTimeout: time.Duration(constants.DefaultWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(constants.DefaultIdleTimeout) * time.Second,
	}

	httpServer := &HTTPServer{
		server:                 server,
		router:                 router,
		userHandler:            userHandler,
		productHandler:         productHandler,
		waitingListHandler:     waitingListHandler,
		settingHandler:         settingHandler,
		vehicleHandler:         vehicleHandler,
		maintenanceItemHandler: maintenanceItemHandler,
		healthHandler:          healthHandler,
		versionHandler:         versionHandler,
		invoiceHandler:         invoiceHandler,
		analyticsHandler:       analyticsHandler,
		roleHandler:            roleHandler,
	}

	httpServer.setupRoutes()

	return httpServer
}

func (s *HTTPServer) setupRoutes() {
	// Health check endpoints
	s.router.HandleFunc("/health", s.healthHandler.HealthCheck).Methods("GET")
	s.router.HandleFunc("/health/live", s.healthHandler.LivenessCheck).Methods("GET")
	s.router.HandleFunc("/health/ready", s.healthHandler.ReadinessCheck).Methods("GET")

	// Detailed health check (admin only)
	healthDetailRoute := s.router.HandleFunc("/health/detail", s.healthHandler.DetailedHealthCheck).Methods("GET").Subrouter()
	healthDetailRoute.Use(middleware.Auth)
	healthDetailRoute.Use(middleware.RequireRole(constants.RoleAdmin))

	// Version info endpoint
	s.router.HandleFunc("/api/version", s.versionHandler.GetVersion).Methods("GET")

	api := s.router.PathPrefix("/api/v1").Subrouter()

	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", s.userHandler.Register).Methods("POST")
	authRoutes.HandleFunc("/login", s.userHandler.Login).Methods("POST")
	authRoutes.HandleFunc("/refresh", s.userHandler.RefreshToken).Methods("POST")

	productRoutes := api.PathPrefix("/products").Subrouter()
	productRoutes.HandleFunc("", s.productHandler.GetAllProducts).Methods("GET")
	productRoutes.HandleFunc("/{id:[0-9]+}", s.productHandler.GetProduct).Methods("GET")

	userRoutes := api.PathPrefix("/users").Subrouter()
	userRoutes.Use(middleware.Auth)
	userRoutes.HandleFunc("/profile", s.userHandler.GetProfile).Methods("GET")
	userRoutes.HandleFunc("/profile", s.userHandler.UpdateProfile).Methods("PUT")

	adminUserRoutes := userRoutes.NewRoute().Subrouter()
	adminUserRoutes.Use(middleware.Auth)
	adminUserRoutes.Use(middleware.RequireRole(constants.RoleAdmin))
	adminUserRoutes.HandleFunc("", s.userHandler.GetUsers).Methods("GET")
	adminUserRoutes.HandleFunc("/{id}", s.userHandler.GetUser).Methods("GET")
	adminUserRoutes.HandleFunc("/{id}", s.userHandler.UpdateUser).Methods("PUT")
	adminUserRoutes.HandleFunc("/{id}", s.userHandler.DeleteUser).Methods("DELETE")

	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.Auth)
	adminRoutes.Use(middleware.RequireRole(constants.RoleAdmin))
	adminProductRoutes := adminRoutes.PathPrefix("/products").Subrouter()
	adminProductRoutes.HandleFunc("", s.productHandler.CreateProduct).Methods("POST")
	adminProductRoutes.HandleFunc("/{id}", s.productHandler.UpdateProduct).Methods("PUT")
	adminProductRoutes.HandleFunc("/{id}/stock", s.productHandler.UpdateProductStock).Methods("PATCH")
	adminProductRoutes.HandleFunc("/{id}", s.productHandler.DeleteProduct).Methods("DELETE")

	// Waiting List Routes (Customer)
	waitingListRoutes := api.PathPrefix("/waiting-list").Subrouter()
	waitingListRoutes.Use(middleware.Auth)
	waitingListRoutes.HandleFunc("/take", s.waitingListHandler.TakeQueueNumber).Methods("POST")
	waitingListRoutes.HandleFunc("/my-queue", s.waitingListHandler.GetMyQueue).Methods("GET")
	waitingListRoutes.HandleFunc("/today", s.waitingListHandler.GetTodayQueue).Methods("GET")
	waitingListRoutes.HandleFunc("/date", s.waitingListHandler.GetQueueByDate).Methods("GET")
	waitingListRoutes.HandleFunc("/number/{number}", s.waitingListHandler.GetQueueByNumber).Methods("GET")
	waitingListRoutes.HandleFunc("/availability", s.waitingListHandler.CheckAvailability).Methods("GET")
	waitingListRoutes.HandleFunc("/{id}/cancel", s.waitingListHandler.CancelQueue).Methods("PUT")
	waitingListRoutes.HandleFunc("/{id}/progress", s.waitingListHandler.GetServiceProgress).Methods("GET")

	// Waiting List Routes (Admin only - manage queue operations)
	adminWaitingListRoutes := adminRoutes.PathPrefix("/waiting-list").Subrouter()
	adminWaitingListRoutes.HandleFunc("/{id}/call", s.waitingListHandler.CallCustomer).Methods("PUT")
	adminWaitingListRoutes.HandleFunc("/{id}/start", s.waitingListHandler.StartService).Methods("PUT")
	adminWaitingListRoutes.HandleFunc("/{id}/complete", s.waitingListHandler.CompleteService).Methods("PUT")
	adminWaitingListRoutes.HandleFunc("/{id}/no-show", s.waitingListHandler.MarkNoShow).Methods("PUT")

	// Maintenance Items Routes (Customer)
	maintenanceRoutes := api.PathPrefix("/maintenance").Subrouter()
	maintenanceRoutes.Use(middleware.Auth)
	maintenanceRoutes.HandleFunc("/waiting-list/{waiting_list_id}/items", s.maintenanceItemHandler.CreateInitialItems).Methods("POST")
	maintenanceRoutes.HandleFunc("/waiting-list/{waiting_list_id}/items", s.maintenanceItemHandler.GetItemsByWaitingList).Methods("GET")
	maintenanceRoutes.HandleFunc("/waiting-list/{waiting_list_id}/inspection-summary", s.maintenanceItemHandler.GetInspectionSummary).Methods("GET")
	maintenanceRoutes.HandleFunc("/items/approve", s.maintenanceItemHandler.ApproveItems).Methods("POST")

	// Maintenance Items Routes (Admin/Mechanic)
	adminMaintenanceRoutes := adminRoutes.PathPrefix("/maintenance").Subrouter()
	adminMaintenanceRoutes.HandleFunc("/items/discovered", s.maintenanceItemHandler.AddDiscoveredItem).Methods("POST")
	adminMaintenanceRoutes.HandleFunc("/items/{id}", s.maintenanceItemHandler.UpdateItem).Methods("PUT")
	adminMaintenanceRoutes.HandleFunc("/items/{id}/complete", s.maintenanceItemHandler.CompleteItem).Methods("PUT")
	adminMaintenanceRoutes.HandleFunc("/items/{id}", s.maintenanceItemHandler.DeleteItem).Methods("DELETE")

	// Vehicle Routes (User can manage their own vehicles)
	vehicleRoutes := api.PathPrefix("/vehicles").Subrouter()
	vehicleRoutes.Use(middleware.Auth)
	vehicleRoutes.HandleFunc("", s.vehicleHandler.CreateVehicle).Methods("POST")
	vehicleRoutes.HandleFunc("", s.vehicleHandler.GetMyVehicles).Methods("GET")
	vehicleRoutes.HandleFunc("/{id}", s.vehicleHandler.GetVehicle).Methods("GET")
	vehicleRoutes.HandleFunc("/{id}", s.vehicleHandler.UpdateVehicle).Methods("PUT")
	vehicleRoutes.HandleFunc("/{id}", s.vehicleHandler.DeleteVehicle).Methods("DELETE")

	// Vehicle Routes (Admin - Get all vehicles)
	adminVehicleRoutes := adminRoutes.PathPrefix("/vehicles").Subrouter()
	adminVehicleRoutes.HandleFunc("", s.vehicleHandler.GetAllVehicles).Methods("GET")

	// Settings Routes (Public - for customers to see shop info)
	settingsPublicRoutes := api.PathPrefix("/settings").Subrouter()
	settingsPublicRoutes.Use(middleware.Auth)
	settingsPublicRoutes.HandleFunc("/public", s.settingHandler.GetPublicSettings).Methods("GET")

	// Settings Routes (Admin only)
	settingsAdminRoutes := adminRoutes.PathPrefix("/settings").Subrouter()
	settingsAdminRoutes.HandleFunc("", s.settingHandler.GetAllSettings).Methods("GET")
	settingsAdminRoutes.HandleFunc("", s.settingHandler.CreateSetting).Methods("POST")
	settingsAdminRoutes.HandleFunc("/category/{category}", s.settingHandler.GetSettingsByCategory).Methods("GET")
	settingsAdminRoutes.HandleFunc("/key/{key}", s.settingHandler.GetSetting).Methods("GET")
	settingsAdminRoutes.HandleFunc("/key/{key}", s.settingHandler.UpdateSetting).Methods("PUT")
	settingsAdminRoutes.HandleFunc("/{id}", s.settingHandler.DeleteSetting).Methods("DELETE")

	// Invoice Routes (Admin)
	invoiceAdminRoutes := adminRoutes.PathPrefix("/invoices").Subrouter()
	invoiceAdminRoutes.HandleFunc("", s.invoiceHandler.CreateInvoice).Methods("POST")
	invoiceAdminRoutes.HandleFunc("", s.invoiceHandler.GetInvoices).Methods("GET")
	invoiceAdminRoutes.HandleFunc("/{id}", s.invoiceHandler.GetInvoice).Methods("GET")
	invoiceAdminRoutes.HandleFunc("/{id}", s.invoiceHandler.UpdateInvoice).Methods("PUT")
	invoiceAdminRoutes.HandleFunc("/{id}", s.invoiceHandler.DeleteInvoice).Methods("DELETE")

	// Invoice Routes (Customer)
	invoiceRoutes := api.PathPrefix("/invoices").Subrouter()
	invoiceRoutes.Use(middleware.Auth)
	invoiceRoutes.HandleFunc("/{id}", s.invoiceHandler.GetInvoice).Methods("GET")
	invoiceRoutes.HandleFunc("/{id}/pay", s.invoiceHandler.PayInvoice).Methods("POST")
	invoiceRoutes.HandleFunc("/{id}/download", s.invoiceHandler.DownloadInvoice).Methods("GET")

	// Waiting List Invoice Route
	waitingListRoutes.HandleFunc("/{id}/invoice", s.invoiceHandler.GetInvoiceByWaitingList).Methods("GET")

	// Analytics Routes (Admin only)
	analyticsRoutes := adminRoutes.PathPrefix("/analytics").Subrouter()
	analyticsRoutes.HandleFunc("/overview", s.analyticsHandler.GetOverview).Methods("GET")
	analyticsRoutes.HandleFunc("/revenue-stats", s.analyticsHandler.GetRevenueStats).Methods("GET")
	analyticsRoutes.HandleFunc("/service-stats", s.analyticsHandler.GetServiceStats).Methods("GET")
	analyticsRoutes.HandleFunc("/queue-stats", s.analyticsHandler.GetQueueStats).Methods("GET")
	analyticsRoutes.HandleFunc("/mechanic-performance", s.analyticsHandler.GetMechanicPerformance).Methods("GET")

	// Role Routes (Admin only)
	roleRoutes := adminRoutes.PathPrefix("/roles").Subrouter()
	roleRoutes.HandleFunc("", s.roleHandler.CreateRole).Methods("POST")
	roleRoutes.HandleFunc("", s.roleHandler.GetAllRoles).Methods("GET")
	roleRoutes.HandleFunc("/active", s.roleHandler.GetActiveRoles).Methods("GET")
	roleRoutes.HandleFunc("/{id}", s.roleHandler.GetRole).Methods("GET")
	roleRoutes.HandleFunc("/{id}", s.roleHandler.UpdateRole).Methods("PUT")
	roleRoutes.HandleFunc("/{id}", s.roleHandler.DeleteRole).Methods("DELETE")
	roleRoutes.HandleFunc("/{id}/users", s.roleHandler.GetUsersByRole).Methods("GET")

	// User Role Assignment Routes (Admin only)
	userRoleRoutes := adminRoutes.PathPrefix("/users/{userId}/roles").Subrouter()
	userRoleRoutes.HandleFunc("", s.roleHandler.GetUserRoles).Methods("GET")
	userRoleRoutes.HandleFunc("", s.roleHandler.AssignRoleToUser).Methods("POST")
	userRoleRoutes.HandleFunc("", s.roleHandler.RemoveRoleFromUser).Methods("DELETE")
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
