# Go REST API Service

A lightweight, production-ready REST API service built with Go and Fiber, featuring JWT authentication, CRUD operations, search, pagination, and comprehensive error handling. Designed for quick deployment to cloud platforms with Docker support.

## Features

- ✅ **8 Progressive Test Levels**: From basic health checks to complete CRUD with authentication
- 🔐 **JWT Authentication**: Secure token-based authentication with configurable expiration
- 📚 **Book Management**: Full CRUD operations for book resources
- 🔍 **Search & Filter**: Case-insensitive author search
- 📄 **Pagination**: Configurable page size and navigation
- 🛡️ **Error Handling**: Consistent error responses with proper HTTP status codes
- 🚀 **Fast & Lightweight**: Built with Fiber framework for high performance
- 🐳 **Docker Ready**: Multi-stage Dockerfile for minimal production images
- ☁️ **Cloud Native**: Easy deployment to Render, Railway, Fly.io
- 🧪 **Well Tested**: Unit tests and property-based tests included

## Quick Start

### Prerequisites

- Go 1.25 or higher
- Docker (optional, for containerized deployment)

### Local Development

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd go-rest-api-service
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set environment variables** (optional):
   ```bash
   export PORT=8080
   export JWT_SECRET=your-secret-key
   export JWT_EXPIRATION=24h
   ```

4. **Run the service**:
   ```bash
   go run main.go
   ```

5. **Test the service**:
   ```bash
   curl http://localhost:8080/ping
   ```
   Expected response: `{"status":"ok"}`

### Using Docker

1. **Build the image**:
   ```bash
   docker build -t go-rest-api-service .
   ```

2. **Run the container**:
   ```bash
   docker run -p 8080:8080 \
     -e JWT_SECRET=your-secret-key \
     go-rest-api-service
   ```

## API Documentation

### Base URL
- Local: `http://localhost:8080`
- Production (Railway): `https://your-service-name.up.railway.app`

### Endpoints

#### Level 1: Health Check

**GET /ping**
- Description: Health check endpoint
- Authentication: None
- Response: `200 OK`
  ```json
  {"status": "ok"}
  ```

#### Level 2: Echo

**POST /echo**
- Description: Echo back the request body
- Authentication: None
- Request Body: Any valid JSON
- Response: `200 OK` with the same JSON body

Example:
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "hello"}'
```

#### Level 5: Authentication

**POST /auth/token**
- Description: Generate JWT authentication token
- Authentication: None
- Request Body:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- Response: `200 OK`
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
- Error: `401 Unauthorized` for invalid credentials

Example:
```bash
curl -X POST http://localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "testpass"}'
```

#### Level 3-4: Book Management

**POST /books**
- Description: Create a new book
- Authentication: None
- Request Body:
  ```json
  {
    "title": "string (required)",
    "author": "string (required)",
    "year": "integer (required, 1000-9999)"
  }
  ```
- Response: `201 Created`
  ```json
  {
    "id": "uuid",
    "title": "The Go Programming Language",
    "author": "Alan Donovan",
    "year": 2015
  }
  ```
- Error: `400 Bad Request` for invalid data

**GET /books/:id**
- Description: Get a book by ID
- Authentication: None
- Response: `200 OK` with book data
- Error: `404 Not Found` if book doesn't exist

**PUT /books/:id**
- Description: Update a book
- Authentication: None
- Request Body: Same as POST /books
- Response: `200 OK` with updated book data
- Error: `404 Not Found` if book doesn't exist

**DELETE /books/:id**
- Description: Delete a book
- Authentication: None
- Response: `204 No Content`
- Error: `404 Not Found` if book doesn't exist

#### Level 5-6: Book List (Protected)

**GET /books**
- Description: List all books with optional filtering and pagination
- Authentication: **Required** (Bearer token)
- Query Parameters:
  - `author` (optional): Filter by author name (case-insensitive)
  - `page` (optional, default: 1): Page number
  - `limit` (optional, default: 10): Items per page
- Response: `200 OK`
  ```json
  {
    "data": [
      {
        "id": "uuid",
        "title": "Book Title",
        "author": "Author Name",
        "year": 2023
      }
    ],
    "page": 1,
    "limit": 10,
    "total_items": 25,
    "total_pages": 3
  }
  ```
- Error: `401 Unauthorized` if token is missing or invalid

Example:
```bash
# List all books
curl http://localhost:8080/books \
  -H "Authorization: Bearer YOUR_TOKEN"

# Search by author
curl "http://localhost:8080/books?author=Alan%20Donovan" \
  -H "Authorization: Bearer YOUR_TOKEN"

# With pagination
curl "http://localhost:8080/books?page=2&limit=5" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Error Responses

All errors follow a consistent format:

```json
{
  "error": "Error message description"
}
```

HTTP Status Codes:
- `400 Bad Request`: Invalid input or validation error
- `401 Unauthorized`: Missing or invalid authentication token
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Unexpected server error

## Testing

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Tests Verbosely

```bash
go test -v ./...
```

### Run Specific Test Package

```bash
# Test handlers
go test ./handlers/...

# Test services
go test ./service/...

# Test repository
go test ./repository/...
```

### Test Structure

The project includes:
- **Unit Tests**: Test specific functions and edge cases
- **Property-Based Tests**: Test universal properties across random inputs (using gopter)
- **Integration Tests**: Test end-to-end API behavior

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | HTTP server port | `8080` | No |
| `JWT_SECRET` | Secret key for JWT signing | `default-secret-key` | **Yes (Production)** |
| `JWT_EXPIRATION` | Token expiration duration | `24h` | No |

### Configuration Examples

**Development**:
```bash
export PORT=8080
export JWT_SECRET=dev-secret-key
export JWT_EXPIRATION=1h
```

**Production**:
```bash
export PORT=8080
export JWT_SECRET=prod-secure-random-32-char-key
export JWT_EXPIRATION=24h
```

## Deployment

### Quick Deploy to Railway (Free - No Credit Card Required)

The easiest way to deploy this API is using Railway's free tier:

1. **Push your code to GitHub**
2. **Sign up at [Railway](https://railway.app)** with GitHub (no credit card needed)
3. **Create a new project**:
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository
   - Railway auto-detects the Dockerfile
4. **Set environment variable**:
   - Go to Variables tab
   - Add: `JWT_SECRET=your-secure-random-32-char-secret`
5. **Generate domain**:
   - Go to Settings → Networking
   - Click "Generate Domain"
   - Your API will be live at `https://your-service-name.up.railway.app`

**Free Tier**: $5 credit/month (~500 hours runtime), no cold starts, no credit card required!

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed instructions and alternative platforms (Render, Fly.io, Docker).

## Project Structure

```
.
├── main.go                 # Application entry point
├── config/
│   └── config.go          # Configuration management
├── handlers/
│   ├── ping.go            # Health check handler
│   ├── echo.go            # Echo handler
│   ├── books.go           # Book CRUD handlers
│   ├── auth.go            # Authentication handler
│   └── *_test.go          # Handler tests
├── middleware/
│   ├── auth.go            # JWT authentication middleware
│   ├── recovery.go        # Panic recovery middleware
│   └── *_test.go          # Middleware tests
├── models/
│   └── book.go            # Data models
├── repository/
│   └── memory.go          # In-memory data store
├── service/
│   ├── book_service.go    # Book business logic
│   └── auth_service.go    # Authentication logic
├── Dockerfile             # Multi-stage Docker build
├── DEPLOYMENT.md          # Deployment guide
└── README.md              # This file
```

## Architecture

The service follows a layered architecture:

1. **Handlers Layer**: HTTP request/response handling
2. **Middleware Layer**: Authentication, error recovery, CORS
3. **Service Layer**: Business logic and orchestration
4. **Repository Layer**: Data access and storage

### Data Flow

```
HTTP Request → Middleware → Handler → Service → Repository → In-Memory Store
                                                              ↓
HTTP Response ← Handler ← Service ← Repository ← In-Memory Store
```

## Development

### Adding New Endpoints

1. Define the model in `models/`
2. Create repository interface and implementation in `repository/`
3. Implement business logic in `service/`
4. Create handler in `handlers/`
5. Register route in `main.go`
6. Write tests

### Code Style

- Follow Go conventions and idioms
- Use `gofmt` for formatting
- Write descriptive comments for exported functions
- Keep functions small and focused
- Handle errors explicitly

### Testing Guidelines

- Write unit tests for all business logic
- Use property-based tests for universal properties
- Test error cases and edge conditions
- Aim for high test coverage
- Use table-driven tests where appropriate

## Security

- **JWT Tokens**: Use strong secrets (32+ characters) in production
- **HTTPS**: Always use HTTPS in production (handled by platforms)
- **Input Validation**: All inputs are validated using Fiber's binding
- **Error Messages**: Generic error messages prevent information leakage
- **Rate Limiting**: Consider adding rate limiting for production
- **CORS**: Configure CORS appropriately for your use case

## Performance

- **Fiber Framework**: High-performance HTTP framework
- **In-Memory Storage**: Fast read/write operations
- **Thread-Safe**: Concurrent request handling with mutex locks
- **Minimal Dependencies**: Small binary size (~10MB)
- **Docker**: Multi-stage build for minimal image size

## Limitations

- **In-Memory Storage**: Data is lost on restart (not persistent)
- **Single Instance**: In-memory store is not shared across instances
- **Basic Auth**: Simple username/password validation (for MVP)

For production use with multiple instances, consider:
- Adding a persistent database (PostgreSQL, MongoDB)
- Implementing proper user authentication
- Adding Redis for shared session storage

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests
5. Run tests: `go test ./...`
6. Submit a pull request

## License

This project is provided as-is for educational and development purposes.

## Support

For issues, questions, or contributions, please open an issue on GitHub.

## Roadmap

Future enhancements:
- [ ] Persistent database support (PostgreSQL)
- [ ] User registration and management
- [ ] Role-based access control (RBAC)
- [ ] Rate limiting middleware
- [ ] Request logging
- [ ] Metrics and monitoring
- [ ] OpenAPI/Swagger documentation
- [ ] GraphQL support
- [ ] WebSocket support

---

Built with ❤️ using Go and Fiber
