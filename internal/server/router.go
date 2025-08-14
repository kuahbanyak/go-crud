package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"os"

	"github.com/kuahbanyak/go-crud/internal/auth"
	"github.com/kuahbanyak/go-crud/internal/booking"
	"github.com/kuahbanyak/go-crud/internal/inventory"
	"github.com/kuahbanyak/go-crud/internal/invoice"
	"github.com/kuahbanyak/go-crud/internal/servicehistory"
	"github.com/kuahbanyak/go-crud/internal/user"
	"github.com/kuahbanyak/go-crud/internal/vehicle"
	"github.com/kuahbanyak/go-crud/pkg/middleware"
)

func NewServer(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	// Register CORS middleware globally
	r.Use(middleware.CORSMiddleware())
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("dev_secret")
	}

	// Repos & services
	uRepo := user.NewRepo(db)
	vRepo := vehicle.NewRepo(db)
	bRepo := booking.NewRepo(db)
	sRepo := servicehistory.NewRepo(db)
	pRepo := inventory.NewRepo(db)
	iRepo := invoice.NewRepo(db)

	// handlers
	authH := auth.NewHandler(uRepo)
	userH := user.NewHandler(uRepo)
	vehicleH := vehicle.NewHandler(vRepo)
	bookingH := booking.NewHandler(bRepo, vRepo)
	serviceH := servicehistory.NewHandler(sRepo)
	inventoryH := inventory.NewHandler(pRepo)
	invoiceH := invoice.NewHandler(iRepo)

	// public
	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	// protected
	authMw := middleware.JWTAuthMiddleware(jwtSecret)
	api := r.Group("/api/v1", authMw)
	{
		api.GET("/me", userH.Me)
		api.PUT("/me", userH.UpdateProfile)

		api.GET("/vehicles", vehicleH.List)
		api.POST("/vehicles", vehicleH.Create)
		api.GET("/vehicles/:id", vehicleH.Get)
		api.PUT("/vehicles/:id", vehicleH.Update)
		api.DELETE("/vehicles/:id", vehicleH.Delete)

		api.POST("/bookings", bookingH.Create)
		api.GET("/bookings", bookingH.List)
		api.GET("/bookings/:id", bookingH.GetId)
		api.PUT("/bookings/:id/status", bookingH.UpdateStatus)

		api.POST("/service-records", serviceH.Create)
		api.GET("/service-records", serviceH.List)

		api.GET("/parts", inventoryH.List)
		api.POST("/parts", inventoryH.Create)
		api.PUT("/parts/:id", inventoryH.Update)

		api.POST("/invoices/generate", invoiceH.Generate)
		api.GET("/reports/summary", invoiceH.Summary)
	}

	return r
}
