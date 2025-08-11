# Car Service Backend (Go) - Scaffold

Minimal scaffold implementing:
- User management (register/login, JWT, roles)
- Vehicle management (CRUD)
- Booking system (schedule, status)
- Service history (records + receipt upload)
- Inventory (parts, stock)
- Invoices (generate PDF stub)
- Notifications (email cron)
- Reporting endpoints (basic aggregates)

Instructions:
1. Copy this repo to your GOPATH or module-enabled workspace.
2. Set environment variables (see .env.example).
3. `go mod tidy`
4. Run PostgreSQL and set DATABASE_DSN accordingly.
5. `go run ./cmd/server`

This scaffold is intentionally compact; replace placeholders and improve security for production.
PORT=8080
DB_DSN=sqlserver://devsql:Kuahpisah1@devdbsql.database.windows.net:1433?database=sqldev
JWT_SECRET=replace_with_secret
STORAGE_PATH=./uploads
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASS=
