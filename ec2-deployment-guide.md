# EC2 Deployment Guide for Docker Containers

## 1. Launch EC2 Instance

### Security Group Settings:
- **SSH (22)**: Your IP address
- **HTTP (80)**: 0.0.0.0/0 (for nginx)
- **Custom TCP (8080)**: 0.0.0.0/0 (optional - for direct http-server access)
- **Custom TCP (8081)**: 0.0.0.0/0 (optional - for direct ws-server access)

## 2. Connect to EC2 and Install Docker

```bash
# Connect to EC2
ssh -i your-key.pem ec2-user@your-ec2-ip

# Update system
sudo yum update -y  # For Amazon Linux
# OR
sudo apt update && sudo apt upgrade -y  # For Ubuntu

# Install Docker
sudo yum install docker -y  # Amazon Linux
# OR  
sudo apt install docker.io -y  # Ubuntu

# Start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group (to run without sudo)
sudo usermod -a -G docker ec2-user  # Amazon Linux
# OR
sudo usermod -a -G docker ubuntu  # Ubuntus

# Log out and log back in for group changes to take effect
exit
# ssh back in

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

## 3. Create Production Docker Compose File

Create a new `docker-compose.prod.yml` on your EC2:

```yaml
version: '3.8'

services:
  http-server:
    image: likhith2005/test-multistage-golang:http-server
    ports:
      - "8080:8080"
    restart: unless-stopped

  ws-server:
    image: likhith2005/test-multistage-golang:ws-server
    ports:
      - "8081:8081"
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - http-server
      - ws-server
    restart: unless-stopped
```

## 4. Copy nginx.conf to EC2

You'll need to copy your nginx configuration to the EC2 instance:

```bash
# On your local machine, copy nginx.conf to EC2
scp -i your-key.pem nginx.conf ec2-user@your-ec2-ip:~/nginx.conf
```

## 5. Deploy on EC2

```bash
# On EC2 instance
# Create project directory
mkdir ~/my-app && cd ~/my-app

# Create the docker-compose.prod.yml file (copy content from above)
nano docker-compose.prod.yml

# Move nginx.conf to current directory
mv ~/nginx.conf ./

# Pull and run the containers
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d

# Check if containers are running
docker ps

# Check logs if needed
docker-compose -f docker-compose.prod.yml logs
```

## 6. Access Your Application

- **Via nginx (port 80)**: `http://your-ec2-public-ip`
- **Direct http-server (port 8080)**: `http://your-ec2-public-ip:8080`
- **Direct ws-server (port 8081)**: `http://your-ec2-public-ip:8081`

## 7. Useful Management Commands

```bash
# Stop all services
docker-compose -f docker-compose.prod.yml down

# Restart services
docker-compose -f docker-compose.prod.yml restart

# Update images and restart
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Check resource usage
docker stats
```

## 8. Optional: Set up Auto-restart on Boot

```bash
# Create systemd service for auto-start
sudo nano /etc/systemd/system/my-app.service
```

Content for the service file:
```ini
[Unit]
Description=My Docker App
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/ec2-user/my-app
ExecStart=/usr/local/bin/docker-compose -f docker-compose.prod.yml up -d
ExecStop=/usr/local/bin/docker-compose -f docker-compose.prod.yml down

[Install]
WantedBy=multi-user.target
```

Enable the service:
```bash
sudo systemctl enable my-app.service
sudo systemctl start my-app.service
```

## 9. Security Best Practices

1. **Use environment variables** for sensitive data
2. **Set up SSL/TLS** with Let's Encrypt for HTTPS
3. **Configure firewall** (ufw on Ubuntu, firewall-cmd on Amazon Linux)
4. **Regular updates**: Keep your system and Docker images updated
5. **Monitor logs**: Set up log aggregation for production

## Troubleshooting

- **Port conflicts**: Make sure ports aren't already in use
- **Permission issues**: Ensure user is in docker group
- **Network issues**: Check security groups and firewall settings
- **Container issues**: Use `docker logs container-name` to debug 