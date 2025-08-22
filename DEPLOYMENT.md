# Deployment Guide for Export API

## Option 1: Deploy to Render (Recommended)

### Prerequisites
- GitHub repository with your code
- Render account (free tier available)

### Steps

1. **Fork/Clone this repository to your GitHub account**

2. **Connect to Render:**
   - Go to [render.com](https://render.com)
   - Sign up/Login with your GitHub account
   - Click "New +" → "Web Service"

3. **Connect your repository:**
   - Select your GitHub repository
   - Choose the branch (main/master)

4. **Configure the service:**
   - **Name**: `export-api` (or any name you prefer)
   - **Environment**: `Go`
   - **Region**: Choose closest to your users
   - **Branch**: `main` (or your default branch)
   - **Build Command**: `go build -o export-api main.go`
   - **Start Command**: `./export-api`

5. **Set Environment Variables:**
   ```
   DB_HOST=myshaa.com
   DB_PORT=3306
   DB_USER=myshaa_kabu
   DB_PASSWORD=T-Cyj;f5g1y6
   DB_NAME=myshaa_kabu
   PORT=8080
   ```

6. **Deploy:**
   - Click "Create Web Service"
   - Render will automatically build and deploy your application

### Alternative: Use render.yaml (Blueprint)

1. **Push your code with the `render.yaml` file**
2. **In Render dashboard:**
   - Click "New +" → "Blueprint"
   - Select your repository
   - Render will automatically create the service using the configuration

## Option 2: Local Development & Testing

### Prerequisites
- Go 1.21 or higher
- MySQL database access

### Steps

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd export-api
   ```

2. **Set environment variables:**
   ```bash
   export DB_HOST=myshaa.com
   export DB_PORT=3306
   export DB_USER=myshaa_kabu
   export DB_PASSWORD=T-Cyj;f5g1y6
   export DB_NAME=myshaa_kabu
   export PORT=8080
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

5. **Test the API:**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # List tables
   curl http://localhost:8080/tables
   
   # Export data (example)
   curl "http://localhost:8080/export?table=GTPL_108_gT_40E_P_S7_200_Germany&all=true" \
        -o export.xlsx
   ```

## Option 3: Docker Deployment

### Prerequisites
- Docker installed
- Docker Hub account (optional)

### Steps

1. **Build the Docker image:**
   ```bash
   docker build -t export-api .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8080:8080 \
     -e DB_HOST=myshaa.com \
     -e DB_PORT=3306 \
     -e DB_USER=myshaa_kabu \
     -e DB_PASSWORD=T-Cyj;f5g1y6 \
     -e DB_NAME=myshaa_kabu \
     export-api
   ```

3. **Test the containerized API:**
   ```bash
   curl http://localhost:8080/health
   ```

## Option 4: VPS/Cloud Server Deployment

### Prerequisites
- VPS or cloud server (AWS EC2, DigitalOcean, etc.)
- SSH access
- Go 1.21+ installed

### Steps

1. **SSH into your server:**
   ```bash
   ssh user@your-server-ip
   ```

2. **Clone and setup:**
   ```bash
   git clone <your-repo-url>
   cd export-api
   go mod download
   go build -o export-api main.go
   ```

3. **Create systemd service:**
   ```bash
   sudo nano /etc/systemd/system/export-api.service
   ```

   Add this content:
   ```ini
   [Unit]
   Description=Export API
   After=network.target
   
   [Service]
   Type=simple
   User=ubuntu
   WorkingDirectory=/home/ubuntu/export-api
   Environment=DB_HOST=myshaa.com
   Environment=DB_PORT=3306
   Environment=DB_USER=myshaa_kabu
   Environment=DB_PASSWORD=T-Cyj;f5g1y6
   Environment=DB_NAME=myshaa_kabu
   Environment=PORT=8080
   ExecStart=/home/ubuntu/export-api/export-api
   Restart=always
   
   [Install]
   WantedBy=multi-user.target
   ```

4. **Start the service:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable export-api
   sudo systemctl start export-api
   sudo systemctl status export-api
   ```

5. **Configure firewall (if needed):**
   ```bash
   sudo ufw allow 8080
   ```

## Troubleshooting

### Common Issues

1. **Database Connection Failed:**
   - Verify database credentials
   - Check if database is accessible from your deployment location
   - Ensure firewall allows connections on port 3306

2. **Port Already in Use:**
   - Change PORT environment variable
   - Kill existing process: `lsof -ti:8080 | xargs kill -9`

3. **Permission Denied:**
   - Ensure the binary is executable: `chmod +x export-api`
   - Check file ownership and permissions

4. **Memory Issues:**
   - The application is designed to handle large datasets
   - Monitor system resources during large exports
   - Consider increasing server memory if needed

### Health Checks

- **Health endpoint**: `GET /health`
- **Tables endpoint**: `GET /tables`
- **Status endpoint**: `GET /status`

### Monitoring

- Check application logs: `journalctl -u export-api -f`
- Monitor database connections
- Watch for memory usage during large exports

## Security Considerations

1. **Database Security:**
   - Use strong passwords
   - Limit database user permissions
   - Consider using SSL connections

2. **API Security:**
   - The API currently allows all origins (CORS: *)
   - Consider restricting CORS in production
   - Add authentication if needed

3. **Environment Variables:**
   - Never commit sensitive data to version control
   - Use environment variables for all secrets
   - Consider using a secrets management service

## Performance Optimization

1. **Database:**
   - Ensure proper indexes on `created_at` and `id` columns
   - Monitor query performance
   - Consider read replicas for large datasets

2. **Application:**
   - The chunked processing is already optimized
   - Monitor memory usage during exports
   - Consider horizontal scaling for high traffic

## Support

If you encounter issues:
1. Check the application logs
2. Verify database connectivity
3. Test endpoints individually
4. Check system resources
5. Review this deployment guide

For additional help, create an issue in the repository or contact the development team.
