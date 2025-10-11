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
	server             *http.Server
	router             *mux.Router
	userHandler        *handlers.UserHandler
	productHandler     *handlers.ProductHandler
	waitingListHandler *handlers.WaitingListHandler
	settingHandler     *handlers.SettingHandler
	vehicleHandler     *handlers.VehicleHandler
}

func NewHTTPServer(
	cfg *config.Config,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
	waitingListHandler *handlers.WaitingListHandler,
	settingHandler *handlers.SettingHandler,
	vehicleHandler *handlers.VehicleHandler,
) *HTTPServer {
	router := mux.NewRouter()
	router.Use(middleware.CORS)
	router.Use(middleware.Logging)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(constants.DefaultReadTimeout) * time.Second,
		WriteTimeout: time.Duration(constants.DefaultWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(constants.DefaultIdleTimeout) * time.Second,
	}

	httpServer := &HTTPServer{
		server:             server,
		router:             router,
		userHandler:        userHandler,
		productHandler:     productHandler,
		waitingListHandler: waitingListHandler,
		settingHandler:     settingHandler,
		vehicleHandler:     vehicleHandler,
	}

	// Setup routes
	httpServer.setupRoutes()

	return httpServer
}

func (s *HTTPServer) setupRoutes() {
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")
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
}

func (s *HTTPServer) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"ok","message":"Server is running"}`))
	if err != nil {
		return
	}
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
