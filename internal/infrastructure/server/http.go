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
	server         *http.Server
	router         *mux.Router
	userHandler    *handlers.UserHandler
	productHandler *handlers.ProductHandler
	bookingHandler *handlers.BookingHandler
}

func NewHTTPServer(
	cfg *config.Config,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
	bookingHandler *handlers.BookingHandler,
) *HTTPServer {
	router := mux.NewRouter()

	// Apply global middleware
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
		server:         server,
		router:         router,
		userHandler:    userHandler,
		productHandler: productHandler,
		bookingHandler: bookingHandler,
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

	userRoutes.HandleFunc("", s.userHandler.GetUsers).Methods("GET")
	userRoutes.HandleFunc("/{id}", s.userHandler.GetUser).Methods("GET")
	userRoutes.HandleFunc("/{id}", s.userHandler.UpdateUser).Methods("PUT")
	userRoutes.HandleFunc("/{id}", s.userHandler.DeleteUser).Methods("DELETE")

	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.Auth)
	adminRoutes.Use(middleware.RequireRole(constants.RoleAdmin))
	adminProductRoutes := adminRoutes.PathPrefix("/products").Subrouter()
	adminProductRoutes.HandleFunc("", s.productHandler.CreateProduct).Methods("POST")
	adminProductRoutes.HandleFunc("/{id}", s.productHandler.UpdateProduct).Methods("PUT")
	adminProductRoutes.HandleFunc("/{id}/stock", s.productHandler.UpdateProductStock).Methods("PATCH")
	adminProductRoutes.HandleFunc("/{id}", s.productHandler.DeleteProduct).Methods("DELETE")

	// Booking routes (all protected)
	bookingRoutes := api.PathPrefix("/bookings").Subrouter()
	bookingRoutes.Use(middleware.Auth)
	bookingRoutes.HandleFunc("", s.bookingHandler.CreateBooking).Methods("POST")
	bookingRoutes.HandleFunc("", s.bookingHandler.GetAllBookings).Methods("GET")
	bookingRoutes.HandleFunc("/{id}", s.bookingHandler.GetBooking).Methods("GET")
	bookingRoutes.HandleFunc("/{id}", s.bookingHandler.UpdateBooking).Methods("PUT")
	bookingRoutes.HandleFunc("/{id}", s.bookingHandler.DeleteBooking).Methods("DELETE")
}

func (s *HTTPServer) healthCheck(w http.ResponseWriter, r *http.Request) {
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
