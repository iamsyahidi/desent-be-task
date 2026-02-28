# Deployment Guide

This guide provides instructions for deploying the Go REST API Service to Render (recommended free tier) and other cloud platforms.

## Environment Variables

The application requires the following environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | HTTP server port | `8080` | No |
| `JWT_SECRET` | Secret key for JWT token signing | `default-secret-key` | **Yes (Production)** |
| `JWT_EXPIRATION` | Token expiration duration | `24h` | No |

**Important**: Always set a strong `JWT_SECRET` in production (minimum 32 characters, use a cryptographically secure random string).

## Health Check Configuration

All platforms should configure health checks using:

- **Endpoint**: `GET /ping`
- **Expected Response**: `200 OK` with `{"status":"ok"}`
- **Interval**: 30 seconds
- **Timeout**: 5 seconds
- **Failure Threshold**: 3 consecutive failures

## Recommended: Deploy to Render (Free Tier)

Render offers a generous free tier perfect for this API service. Follow these steps:

### Prerequisites
- GitHub account
- Render account (sign up at https://render.com)
- Your code pushed to a GitHub repository

### Deployment Steps

1. **Sign up/Login to Render**:
   - Go to https://render.com
   - Sign up or login with your GitHub account

2. **Create a new Web Service**:
   - Click "New +" button in the dashboard
   - Select "Web Service"
   - Connect your GitHub repository
   - Grant Render access to your repository

3. **Configure the service**:
   - **Name**: Choose a unique name (e.g., `my-go-api`)
   - **Region**: Choose closest to your users (e.g., Oregon, Frankfurt)
   - **Branch**: `main` (or your default branch)
   - **Root Directory**: Leave empty (unless your code is in a subdirectory)
   - **Environment**: Docker
   - **Instance Type**: Free

4. **Set environment variables**:
   - Click "Advanced" or go to "Environment" tab
   - Add the following variables:
     ```
     JWT_SECRET=your-secure-random-secret-key-min-32-chars
     JWT_EXPIRATION=24h
     ```
   - **Important**: Generate a strong JWT_SECRET (you can use: `openssl rand -base64 32`)

5. **Configure health check** (optional but recommended):
   - Health Check Path: `/ping`
   - This helps Render monitor your service

6. **Deploy**:
   - Click "Create Web Service"
   - Render will automatically build and deploy your Docker container
   - First deployment takes 2-5 minutes

7. **Get your public URL**:
   - Your service will be available at: `https://your-service-name.onrender.com`
   - Copy this URL for testing

### Free Tier Limitations
- Service spins down after 15 minutes of inactivity
- First request after spin-down takes ~30 seconds (cold start)
- 750 hours/month of runtime (sufficient for most use cases)
- Automatic HTTPS included

### Testing Your Deployment

Once deployed, test with:
```bash
curl https://your-service-name.onrender.com/ping
```

Expected response: `{"status":"ok"}`

## Alternative Platforms

### Railway

1. **Create a new project** in Railway
2. **Deploy from GitHub**:
   - Connect your GitHub repository
   - Railway will auto-detect the Dockerfile
3. **Configure environment variables**:
   - Go to your service settings
   - Add variables:
     ```
     JWT_SECRET=your-secure-random-secret-key-here
     JWT_EXPIRATION=24h
     ```
4. **Configure health check** (optional):
   - Railway automatically monitors your service
   - Health check endpoint: `/ping`
5. **Deploy**: Railway will automatically deploy on push

The service will be available at: `https://your-service-name.up.railway.app`

### Fly.io

1. **Install Fly CLI**:
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. **Login to Fly**:
   ```bash
   fly auth login
   ```

3. **Launch your app**:
   ```bash
   fly launch
   ```
   - Choose your app name
   - Select a region
   - Don't deploy yet (we need to set secrets first)

4. **Set secrets**:
   ```bash
   fly secrets set JWT_SECRET=your-secure-random-secret-key-here
   fly secrets set JWT_EXPIRATION=24h
   ```

5. **Configure health checks** (fly.toml):
   ```toml
   [[services]]
     http_checks = []
     internal_port = 8080
     processes = ["app"]
     protocol = "tcp"
     script_checks = []

     [[services.http_checks]]
       interval = "30s"
       grace_period = "5s"
       method = "get"
       path = "/ping"
       protocol = "http"
       timeout = "5s"
       tls_skip_verify = false
   ```

6. **Deploy**:
   ```bash
   fly deploy
   ```

The service will be available at: `https://your-app-name.fly.dev`

### Vercel (Serverless)

**Note**: Vercel requires a different approach as it uses serverless functions. The current implementation is optimized for traditional server deployment. For Vercel, consider using Render, Railway, or Fly.io instead.

## Docker Deployment (Self-Hosted)

### Build the Docker image:
```bash
docker build -t go-rest-api-service .
```

### Run the container:
```bash
docker run -d \
  -p 8080:8080 \
  -e JWT_SECRET=your-secure-random-secret-key-here \
  -e JWT_EXPIRATION=24h \
  --name api-service \
  go-rest-api-service
```

### Check container logs:
```bash
docker logs api-service
```

### Stop the container:
```bash
docker stop api-service
docker rm api-service
```

## Testing Your Deployment

Once deployed to Render, test your API using these curl commands (replace `your-service-name.onrender.com` with your actual Render URL):

### 1. Health Check (Level 1)
```bash
curl https://your-service-name.onrender.com/ping
```
Expected: `{"status":"ok"}`

### 2. Echo Test (Level 2)
```bash
curl -X POST https://your-service-name.onrender.com/echo \
  -H "Content-Type: application/json" \
  -d '{"message":"hello world"}'
```
Expected: `{"message":"hello world"}`

### 3. Create a Book (Level 3)
```bash
curl -X POST https://your-service-name.onrender.com/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language",
    "author": "Alan Donovan",
    "year": 2015
  }'
```
Expected: `201 Created` with book data including generated ID

### 4. Get Book by ID (Level 4)
```bash
# Replace {book-id} with the ID from step 3
curl https://your-service-name.onrender.com/books/{book-id}
```
Expected: `200 OK` with book data

### 5. Generate Auth Token (Level 5)
```bash
curl -X POST https://your-service-name.onrender.com/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass"
  }'
```
Expected: `200 OK` with `{"token":"eyJhbGc..."}`

### 6. List Books with Auth (Level 5)
```bash
# Replace {token} with the token from step 5
curl https://your-service-name.onrender.com/books \
  -H "Authorization: Bearer {token}"
```
Expected: `200 OK` with paginated book list

### 7. Search Books by Author (Level 6)
```bash
curl "https://your-service-name.onrender.com/books?author=Alan%20Donovan" \
  -H "Authorization: Bearer {token}"
```
Expected: `200 OK` with filtered results

### 8. Pagination (Level 6)
```bash
curl "https://your-service-name.onrender.com/books?page=1&limit=5" \
  -H "Authorization: Bearer {token}"
```
Expected: `200 OK` with pagination metadata

### 9. Update Book (Level 4)
```bash
curl -X PUT https://your-service-name.onrender.com/books/{book-id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language (Updated)",
    "author": "Alan Donovan",
    "year": 2016
  }'
```
Expected: `200 OK` with updated book data

### 10. Delete Book (Level 4)
```bash
curl -X DELETE https://your-service-name.onrender.com/books/{book-id}
```
Expected: `204 No Content`

### 11. Error Handling (Level 7)
```bash
# Test 404 - Non-existent book
curl https://your-service-name.onrender.com/books/non-existent-id

# Test 400 - Invalid book data
curl -X POST https://your-service-name.onrender.com/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Missing required fields"}'

# Test 401 - Missing auth token
curl https://your-service-name.onrender.com/books
```

## Monitoring and Logs

### Render
- View logs in the Render dashboard under your service
- Logs are available in real-time

### Railway
- View logs in the Railway dashboard
- Use `railway logs` CLI command

### Fly.io
- View logs with: `fly logs`
- Real-time logs: `fly logs -f`

### Docker
- View logs with: `docker logs api-service`
- Follow logs: `docker logs -f api-service`

## Troubleshooting

### Service won't start
- Check that `JWT_SECRET` is set
- Verify port 8080 is not already in use
- Check logs for error messages

### Health check failing
- Ensure `/ping` endpoint is accessible
- Verify the service is listening on port 8080
- Check firewall rules

### Authentication not working
- Verify `JWT_SECRET` is set correctly
- Check token expiration settings
- Ensure Authorization header format: `Bearer {token}`

### 500 Internal Server Error
- Check application logs for stack traces
- Verify all environment variables are set
- Check for panics in the recovery middleware logs

## Security Recommendations

1. **Always use HTTPS in production** (handled automatically by most platforms)
2. **Set a strong JWT_SECRET** (minimum 32 characters, cryptographically random)
3. **Enable CORS** with appropriate origins if needed
4. **Add rate limiting** for production use
5. **Monitor logs** for suspicious activity
6. **Keep dependencies updated** regularly
7. **Use secrets management** provided by your platform

## Scaling

### Horizontal Scaling
All platforms support horizontal scaling:
- **Render**: Increase instance count in service settings
- **Railway**: Add replicas in service settings
- **Fly.io**: Use `fly scale count N` command

### Vertical Scaling
Increase resources per instance:
- **Render**: Upgrade instance type
- **Railway**: Increase memory/CPU limits
- **Fly.io**: Use `fly scale vm` command

**Note**: The in-memory storage is not shared across instances. For production with multiple instances, consider adding a shared database (PostgreSQL, Redis, etc.).
