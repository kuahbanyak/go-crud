package constants

const (
	DefaultDBTimeout   = 30 // seconds
	MaxDBConnections   = 25
	MaxIdleConnections = 5

	// JWT constants
	DefaultJWTExpiration = 24 // hours
	JWTIssuer            = "go-crud-api"

	// HTTP constants
	DefaultPort         = ":8080"
	DefaultReadTimeout  = 15 // seconds
	DefaultWriteTimeout = 15 // seconds
	DefaultIdleTimeout  = 60 // seconds

	// Cache constants
	DefaultCacheExpiration = 300 // seconds (5 minutes)
	CacheKeyPrefix         = "go-crud:"

	// Validation constants
	MinPasswordLength    = 8
	MaxNameLength        = 255
	MaxDescriptionLength = 1000

	// Pagination constants
	DefaultPageSize = 10
	MaxPageSize     = 100

	// User roles
	RoleAdmin   = "admin"
	RoleUser    = "user"
	RoleManager = "manager"

	// Product categories
	CategoryElectronics = "electronics"
	CategoryClothing    = "clothing"
	CategoryBooks       = "books"
	CategoryHome        = "home"
	CategorySports      = "sports"
)
