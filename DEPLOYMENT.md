# Docker Deployment Guide

Complete guide for deploying the Groceries backend to production using Docker and Docker Compose.

## 📋 Prerequisites

- VPS with SSH access (IP: 91.99.95.75)
- Docker installed on VPS
- Docker Compose installed on VPS
- SSH key pair for deployment
- GitHub repository with secrets configured

## 🚀 Quick Start

### 1. Set Up VPS

```bash
# SSH into your VPS
ssh user@91.99.95.75

# Create deployment directory
sudo mkdir -p /home/groceries/backend
sudo chown -R $USER:$USER /home/groceries
```

### 2. Configure GitHub Secrets

Add these secrets in GitHub (Settings → Secrets):

- `VPS_HOST`: 91.99.95.75
- `VPS_USER`: your SSH username
- `SSH_PRIVATE_KEY`: your SSH private key

### 3. Deploy

Push to main branch:

```bash
git push origin main
```

The deployment will:
1. Build and test the code
2. Deploy to VPS via SSH
3. Build Docker images
4. Start containers with Docker Compose

## 📂 VPS Directory Structure

```
/home/groceries/
├── backend/
│   ├── main.go
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── ... (all source files)
```

## 🐳 Docker Services

The deployment creates these containers:

1. **groceries-api** - Backend API (port 80)
2. **groceries-mongodb** - MongoDB database
3. **groceries-redis** - Redis cache

Optional tools (disabled by default):
- **groceries-mongo-express** - MongoDB UI (port 8081)
- **groceries-redis-commander** - Redis UI (port 8082)

## 🌐 Access Points

- **API**: http://91.99.95.75/api/v1
- **Health Check**: http://91.99.95.75/api/v1/health
- **Swagger**: http://91.99.95.75/swagger/index.html

## 🛠️ Manual Deployment

If you need to deploy manually:

```bash
# SSH into VPS
ssh user@91.99.95.75

# Navigate to deployment directory
cd /home/groceries/backend

# Stop existing containers
docker-compose down

# Build and start containers
docker-compose up -d --build

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

## 🔍 Troubleshooting

### Check Container Status

```bash
docker-compose ps
docker-compose logs [service-name]
```

### Restart a Service

```bash
docker-compose restart [service-name]
```

### View Real-time Logs

```bash
docker-compose logs -f groceries-api
```

### Connect to Database

```bash
docker exec -it groceries-mongodb mongosh
```

## 🔒 Security Considerations

1. Change default MongoDB credentials
2. Use environment variables for secrets
3. Enable firewall rules on VPS
4. Use HTTPS with reverse proxy (Nginx/Caddy)
5. Regularly update Docker images

## 📊 Monitoring

### Health Check

```bash
curl http://localhost/api/v1/health
```

### Container Health

```bash
docker-compose ps
```

### Resource Usage

```bash
docker stats
```

## 🔄 Updates

To update the application:

1. Push changes to `main` branch
2. GitHub Actions will automatically deploy
3. Or manually run `docker-compose up -d --build`

## 🗑️ Cleanup

To completely remove the deployment:

```bash
cd /home/groceries/backend
docker-compose down -v  # Also removes volumes
rm -rf /home/groceries/backend
```

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
