# Production Deployment Guide

## 🚀 Deployment Options

### Option 1: Docker Deployment (Recommended)
**Pros:** Easy scaling, consistent environment, easy rollbacks
**Best for:** Cloud servers, containerized environments

### Option 2: Traditional Binary Deployment
**Pros:** Lower resource usage, direct system integration
**Best for:** VPS, dedicated servers, existing infrastructure

---

## 🐳 Docker Production Deployment

### Prerequisites
- Docker & Docker Compose installed
- Domain name configured
- SSL certificates (optional but recommended)

### Quick Start
1. **Prepare environment:**
   ```bash
   # Copy and configure environment
   cp .env.prod.template .env.prod
   # Edit .env.prod with your actual values
   ```

2. **Deploy:**
   ```bash
   # Windows
   deploy-prod.bat
   
   # Linux/Mac
   chmod +x deploy-prod.sh
   ./deploy-prod.sh
   ```

### Manual Docker Deployment
```bash
# Build and deploy
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs -f go-crud-api
```

---

## 🖥️ Non-Docker Production Deployment

### Prerequisites
- Go 1.24+ installed
- SQL Server database available
- Reverse proxy (Nginx recommended)

### Build & Deploy
```bash
# 1. Build for production
CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o go-crud-api ./cmd/api

# 2. Create production directory
mkdir -p /opt/go-crud
cp go-crud-api /opt/go-crud/
cp -r configs /opt/go-crud/

# 3. Create systemd service (Linux)
sudo cp scripts/go-crud.service /etc/systemd/system/
sudo systemctl enable go-crud
sudo systemctl start go-crud
```

### Windows Service
```powershell
# Install as Windows service using NSSM
nssm install go-crud-api "C:\go-crud\go-crud-api.exe"
nssm set go-crud-api AppDirectory "C:\go-crud"
nssm start go-crud-api
```

---

## 🗃️ Database Configuration

### Production Database Options:

1. **Azure SQL Database (Cloud)**
   - Managed service, automatic backups
   - Connection: `your-server.database.windows.net:1433`

2. **AWS RDS SQL Server (Cloud)**
   - Managed service, multi-AZ deployment
   - Connection: `your-rds-endpoint.region.rds.amazonaws.com:1433`

3. **Self-Hosted SQL Server**
   - Full control, custom configuration
   - Can be containerized or traditional installation

### Database Setup Steps:
1. Choose your database option (see DATABASE_SETUP.md)
2. Create database and user
3. Configure connection in .env.prod
4. GORM will auto-migrate tables on first run

---

## 🔒 Security Checklist

- [ ] Use strong JWT secret (32+ characters)
- [ ] Enable HTTPS with valid SSL certificates
- [ ] Configure CORS for your domain only
- [ ] Use database user with minimal permissions
- [ ] Set up firewall rules
- [ ] Enable rate limiting
- [ ] Configure proper logging

---

## 📊 Monitoring & Maintenance

### Health Checks
```bash
# API health
curl https://yourdomain.com/health

# Database connection
docker-compose -f docker-compose.prod.yml exec go-crud-api /app/main --health-check
```

### Logs
```bash
# Docker logs
docker-compose -f docker-compose.prod.yml logs -f

# System logs (non-Docker)
journalctl -u go-crud -f
```

### Backup
- Database: Automated daily backups
- Application: Version control + deployment scripts
- Configs: Secure backup of .env.prod

---

## 🚀 Scaling Options

### Horizontal Scaling (Multiple Instances)
```yaml
# docker-compose.prod.yml
services:
  go-crud-api:
    deploy:
      replicas: 3
  
  nginx:
    # Load balancer configuration
```

### Vertical Scaling (More Resources)
```yaml
deploy:
  resources:
    limits:
      cpus: '2.0'
      memory: 1G
```
