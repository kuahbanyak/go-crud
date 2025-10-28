package constants
const (
	DefaultDBTimeout   = 30 // seconds
	MaxDBConnections   = 25
	MaxIdleConnections = 5
	DefaultJWTExpiration = 24 // hours
	JWTIssuer            = "go-crud-api"
	DefaultPort         = ":8080"
	DefaultReadTimeout  = 15 // seconds
	DefaultWriteTimeout = 15 // seconds
	DefaultIdleTimeout  = 60 // seconds
	DefaultCacheExpiration = 300
	CacheKeyPrefix         = "go-crud:"
	MinPasswordLength    = 8
	MaxNameLength        = 255
	MaxDescriptionLength = 1000
	DefaultPageSize = 10
	MaxPageSize     = 100
	RoleAdmin    = "admin"
	RoleUser     = "user"
	RoleMechanic = "mechanic"
	CategoryElectronics = "electronics"
	CategoryClothing    = "clothing"
	CategoryBooks       = "books"
	CategoryHome        = "home"
	CategorySports      = "sports"
)

