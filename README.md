# Go HTTP and WebSocket Server with Nginx

A Go application featuring HTTP and WebSocket servers, with Nginx as a reverse proxy. This project demonstrates how to containerize and deploy a Go application with multiple services.

## Table of Contents
- [Overview](#overview)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Quick Deployment Guide](#quick-deployment-guide)
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [EC2 Deployment](#ec2-deployment)
- [Troubleshooting](#troubleshooting)
- [Maintenance](#maintenance)

## Overview

This project consists of:
- HTTP Server (port 8080)
- WebSocket Server (port 8081)
- Nginx Reverse Proxy (port 80)

## Project Structure

```
ngnix-ws-http/
├── cmd/
│   ├── http-server/
│   │   ├── Dockerfile
│   │   └── http-server.go
│   └── ws-server/
│       ├── Dockerfile
│       └── ws-server.go
├── docker-compose.yml
├── docker-compose.prod.yml
├── nginx.conf
├── go.mod
└── go.sum
```

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later
- Git
- AWS Account (for EC2 deployment)
- Docker Hub account (for pushing images)

## Quick Deployment Guide

### 1. Local Setup (Before EC2)
```bash
# 1. Login to Docker Hub
docker login

# 2. Build and push images
docker buildx build --platform linux/amd64 -t your-dockerhub-username/test-multistage-golang:http-server -f cmd/http-server/Dockerfile . --push
docker buildx build --platform linux/amd64 -t your-dockerhub-username/test-multistage-golang:ws-server -f cmd/ws-server/Dockerfile . --push

# 3. Verify images are pushed
docker images | grep your-dockerhub-username
```

### 2. EC2 Setup
```bash
# 1. Connect to EC2
ssh ubuntu@<ec2-public-ip>

# 2. Install dependencies
sudo apt-get update
sudo apt-get install -y docker.io docker-compose git

# 3. Add user to docker group
sudo usermod -aG docker ubuntu
# Log out and log back in for changes to take effect
exit
ssh ubuntu@<ec2-public-ip>

# 4. Clone and deploy
git clone <repository-url> ~/ngnix-ws-http
cd ~/ngnix-ws-http

# 5. Start services
docker-compose -f docker-compose.prod.yml up -d

# 6. Verify deployment
docker ps
curl http://localhost/
curl http://localhost/api/
curl http://localhost/ws-ping
```

### 3. Security Group Configuration
- Inbound Rules:
  - SSH (22): 0.0.0.0/0
  - HTTP (80): 0.0.0.0/0
  - Custom TCP (8080): 0.0.0.0/0
  - Custom TCP (8081): 0.0.0.0/0

### 4. Verification Checklist
- [ ] Docker images are pushed to Docker Hub
- [ ] Security group ports are open
- [ ] All containers are running (`docker ps`)
- [ ] HTTP server responds (`curl http://<ec2-ip>/`)
- [ ] API endpoint works (`curl http://<ec2-ip>/api/`)
- [ ] WebSocket ping responds (`curl http://<ec2-ip>/ws-ping`)
- [ ] No errors in logs (`docker-compose -f docker-compose.prod.yml logs`)

## Local Development

1. Clone the repository:
```bash
git clone <repository-url>
cd ngnix-ws-http
```

2. Install dependencies:
```bash
go mod download
```

3. Run locally:
```bash
# Run HTTP server
go run cmd/http-server/http-server.go

# Run WebSocket server
go run cmd/ws-server/ws-server.go
```

## Docker Deployment

1. Build and push Docker images:
```bash
# HTTP Server
docker buildx build --platform linux/amd64 -t your-dockerhub-username/test-multistage-golang:http-server -f cmd/http-server/Dockerfile . --push

# WebSocket Server
docker buildx build --platform linux/amd64 -t your-dockerhub-username/test-multistage-golang:ws-server -f cmd/ws-server/Dockerfile . --push
```

2. Start containers:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

## EC2 Deployment

1. Install required software on EC2:
```bash
# Update system
sudo apt-get update
sudo apt-get upgrade -y

# Install Docker
sudo apt-get install -y docker.io

# Install Docker Compose
sudo apt-get install -y docker-compose

# Add user to docker group
sudo usermod -aG docker ubuntu

# Install Git
sudo apt-get install -y git
```

2. Deploy application:
```bash
# Clone repository
git clone <repository-url> ~/ngnix-ws-http
cd ~/ngnix-ws-http

# Start containers
docker-compose -f docker-compose.prod.yml up -d
```

3. Configure Security Group:
- Open port 80 (HTTP)
- Open port 8080 (HTTP Server)
- Open port 8081 (WebSocket Server)
- Open port 22 (SSH)

## Testing the Deployment

1. Test HTTP Server:
```bash
curl http://<ec2-public-ip>/
```

2. Test API Endpoint:
```bash
curl http://<ec2-public-ip>/api/
```

3. Test WebSocket Ping:
```bash
curl http://<ec2-public-ip>/ws-ping
```

## Troubleshooting

### Common Issues

1. **Platform Mismatch**
   - Error: `platform (linux/arm64/v8) does not match the detected host platform (linux/amd64/v3)`
   - Solution: Always build with `--platform linux/amd64` flag

2. **Permission Issues**
   - Error: `Got permission denied while trying to connect to the Docker daemon`
   - Solution: Add user to docker group and relogin

3. **Port Conflicts**
   - Error: `Error starting userland proxy: listen tcp 0.0.0.0:80: bind: address already in use`
   - Solution: Stop existing services or change port mappings

### Checking Logs

```bash
# All containers
docker-compose -f docker-compose.prod.yml logs

# Specific container
docker-compose -f docker-compose.prod.yml logs http-server
docker-compose -f docker-compose.prod.yml logs ws-server
docker-compose -f docker-compose.prod.yml logs nginx
```

## Maintenance

### Update Application

```bash
# Pull latest changes
git pull

# Rebuild and restart
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d
```

### Cleanup

```bash
# Remove unused containers
docker container prune

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune
```

### Monitoring

```bash
# Check container status
docker ps

# Monitor resource usage
docker stats

# View logs in real-time
docker-compose -f docker-compose.prod.yml logs -f
```

## Best Practices

1. **Security**
   - Keep system and packages updated
   - Use non-root user in containers
   - Implement proper security groups
   - Use environment variables for sensitive data

2. **Performance**
   - Use multi-stage builds
   - Implement proper caching
   - Monitor resource usage
   - Use appropriate restart policies

3. **Maintenance**
   - Regular updates
   - Log monitoring
   - Backup configuration
   - Document changes

## License

This project is licensed under the MIT License - see the LICENSE file for details.
