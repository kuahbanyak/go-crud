package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"os"

	"github.com/kuahbanyak/go-crud/internal/auth"
	"github.com/kuahbanyak/go-crud/internal/booking"
	"github.com/kuahbanyak/go-crud/internal/dashboard"
	"github.com/kuahbanyak/go-crud/internal/inventory"
	"github.com/kuahbanyak/go-crud/internal/invoice"
	"github.com/kuahbanyak/go-crud/internal/message"
	"github.com/kuahbanyak/go-crud/internal/notification"
	"github.com/kuahbanyak/go-crud/internal/scheduling"
	"github.com/kuahbanyak/go-crud/internal/servicehistory"
	"github.com/kuahbanyak/go-crud/internal/servicepackage"
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

	// Initialize notification hub
	hub := notification.NewHub()
	go hub.Run()

	// Repos & services
	uRepo := user.NewRepo(db)
	vRepo := vehicle.NewRepo(db)
	bRepo := booking.NewRepo(db)
	sRepo := servicehistory.NewRepo(db)
	pRepo := inventory.NewRepo(db)
	iRepo := invoice.NewRepo(db)
	mRepo := message.NewRepo(db)
	schRepo := scheduling.NewRepo(db)
	spRepo := servicepackage.NewRepo(db)
	dashRepo := dashboard.NewRepo(db)

	// handlers
	authH := auth.NewHandler(uRepo)
	userH := user.NewHandler(uRepo)
	vehicleH := vehicle.NewHandler(vRepo)
	bookingH := booking.NewHandler(bRepo, vRepo)
	serviceH := servicehistory.NewHandler(sRepo)
	inventoryH := inventory.NewHandler(pRepo)
	invoiceH := invoice.NewHandler(iRepo)
	messageH := message.NewHandler(mRepo, hub)
	schedulingH := scheduling.NewHandler(schRepo, hub)
	servicePackageH := servicepackage.NewHandler(spRepo)
	dashboardH := dashboard.NewHandler(dashRepo)

	// public
	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	// protected
	authMw := middleware.JWTAuthMiddleware(jwtSecret)
	api := r.Group("/api/v1", authMw)

	// WebSocket endpoint for real-time notifications
	api.GET("/ws", hub.HandleWebSocket)

	// existing endpoints
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

		// Custom invoice body templates
		api.POST("/invoices/templates", invoiceH.CreateCustomBody)
		api.GET("/invoices/templates", invoiceH.GetCustomBodies)
		api.GET("/invoices/templates/:id", invoiceH.GetCustomBody)
		api.PUT("/invoices/templates/:id", invoiceH.UpdateCustomBody)
		api.DELETE("/invoices/templates/:id", invoiceH.DeleteCustomBody)
		api.PUT("/invoices/templates/:id/set-default", invoiceH.SetDefaultCustomBody)
	}

	// Enhanced Feature 1: Real-time Notifications & Communication
	messages := api.Group("/messages")
	{
		messages.POST("", messageH.Create)
		messages.GET("/booking/:booking_id", messageH.GetByBooking)
		messages.GET("/conversation/:booking_id/:user_id", messageH.GetConversation)
		messages.PUT("/:message_id/read", messageH.MarkAsRead)
		messages.GET("/unread-count", messageH.GetUnreadCount)
	}

	// Enhanced Feature 2: Advanced Scheduling & Calendar Management
	scheduling := api.Group("/scheduling")
	{
		// Mechanic availability
		scheduling.POST("/availability", schedulingH.CreateAvailability)
		scheduling.GET("/availability/mechanic/:mechanic_id", schedulingH.GetMechanicAvailability)

		// Service types
		scheduling.POST("/service-types", schedulingH.CreateServiceType)
		scheduling.GET("/service-types", schedulingH.GetServiceTypes)

		// Maintenance reminders
		scheduling.POST("/reminders", schedulingH.CreateReminder)
		scheduling.GET("/reminders/vehicle/:vehicle_id", schedulingH.GetVehicleReminders)
		scheduling.GET("/reminders/due", schedulingH.GetDueReminders)

		// Waitlist
		scheduling.POST("/waitlist", schedulingH.AddToWaitlist)
	}

	// Enhanced Feature 3: Service Categories & Packages
	packages := api.Group("/packages")
	{
		packages.POST("", servicePackageH.CreatePackage)
		packages.GET("", servicePackageH.GetPackages)
		packages.GET("/:id", servicePackageH.GetPackageByID)

		packages.POST("/categories", servicePackageH.CreateCategory)
		packages.GET("/categories", servicePackageH.GetCategories)

		packages.POST("/history", servicePackageH.CreateServiceHistory)
		packages.GET("/history/vehicle/:vehicle_id", servicePackageH.GetVehicleHistory)
	}

	// Enhanced Feature 4: Customer Portal Enhancements
	dashboard := api.Group("/dashboard")
	{
		dashboard.GET("/customer", dashboardH.GetCustomerDashboard)
		dashboard.GET("/vehicle/:vehicle_id", dashboardH.GetVehicleDashboard)
		dashboard.PUT("/vehicle-health", dashboardH.UpdateVehicleHealth)
		dashboard.POST("/recommendations", dashboardH.CreateRecommendation)
		dashboard.PUT("/budget", dashboardH.UpdateBudget)
	}

	return r
}
