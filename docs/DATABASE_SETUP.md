# Production Database Setup Guide

## Database Options for Production

### Option 1: Cloud Database Services (Recommended)

#### Azure SQL Database
```bash
# 1. Create Azure SQL Database
az sql server create --name your-sql-server --resource-group your-rg --location eastus --admin-user sqladmin
az sql db create --server your-sql-server --resource-group your-rg --name gocrud-prod --service-objective S1

# 2. Configure firewall (allow your server IP)
az sql server firewall-rule create --server your-sql-server --resource-group your-rg --name AllowMyIP --start-ip-address YOUR_SERVER_IP --end-ip-address YOUR_SERVER_IP

# 3. Connection string format:
# sqlserver://username:password@server.database.windows.net:1433?database=gocrud-prod&encrypt=true&trustServerCertificate=false
```

#### AWS RDS SQL Server
```bash
# 1. Create RDS instance
aws rds create-db-instance \
    --db-instance-identifier gocrud-prod \
    --db-instance-class db.t3.micro \
    --engine sqlserver-ex \
    --master-username admin \
    --master-user-password YourStrongPassword \
    --allocated-storage 20

# 2. Configure security group to allow port 1433
# 3. Connection string format:
# sqlserver://admin:password@your-rds-endpoint.region.rds.amazonaws.com:1433?database=gocrud-prod
```

### Option 2: Self-Hosted Database

#### Docker SQL Server (Production)
```yaml
# Add to docker-compose.prod.yml if you want self-hosted DB
  sqlserver-prod:
    image: mcr.microsoft.com/mssql/server:2022-latest
    environment:
      - ACCEPT_EULA=Y
      - SA_PASSWORD=${DB_PASSWORD}
      - MSSQL_PID=Express
    ports:
      - "1433:1433"
    volumes:
      - sqlserver_prod_data:/var/opt/mssql
      - ./scripts/prod-init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - app-network
    restart: unless-stopped
```

#### Traditional SQL Server Installation
```sql
-- Run this on your SQL Server instance
CREATE DATABASE gocrud_prod;
GO

USE gocrud_prod;
GO

-- Create application user (don't use sa in production)
CREATE LOGIN gocrud_user WITH PASSWORD = 'YourStrongPassword123!';
CREATE USER gocrud_user FOR LOGIN gocrud_user;
ALTER ROLE db_owner ADD MEMBER gocrud_user;
GO
```

## Production Database Configuration

### Environment Variables
```bash
# Azure SQL Database
DB_HOST=your-server.database.windows.net
DB_PORT=1433
DB_USER=sqladmin
DB_PASSWORD=YourStrongPassword123!
DB_DATABASE=gocrud-prod

# AWS RDS
DB_HOST=your-rds-endpoint.region.rds.amazonaws.com
DB_PORT=1433
DB_USER=admin
DB_PASSWORD=YourStrongPassword123!
DB_DATABASE=gocrud-prod

# Self-hosted
DB_HOST=your-sql-server-ip
DB_PORT=1433
DB_USER=gocrud_user
DB_PASSWORD=YourStrongPassword123!
DB_DATABASE=gocrud_prod
```

### Connection Pooling & Performance
```yaml
# Add to config.prod.yaml
database:
  host: "${DB_HOST}"
  port: "${DB_PORT}"
  user: "${DB_USER}"
  password: "${DB_PASSWORD}"
  database: "${DB_DATABASE}"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: "1h"
  ssl_mode: "require"
```

### Backup Strategy
```bash
# Automated backups for Azure SQL
az sql db export --admin-password YourPassword --admin-user sqladmin --storage-key YourStorageKey --storage-key-type StorageAccessKey --storage-uri https://yourstorageaccount.blob.core.windows.net/backups/backup.bacpac --name gocrud-prod --resource-group your-rg --server your-sql-server

# For self-hosted SQL Server
docker exec sqlserver-prod /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P YourPassword -Q "BACKUP DATABASE [gocrud_prod] TO DISK = N'/var/opt/mssql/backup/gocrud_prod.bak'"
```
