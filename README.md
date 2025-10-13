# LoveApp Backend API 💕

A todo management API designed for couples, allowing two users (Anyel and Alexis) to create, manage, and complete todos together.

## 🚀 Features

- **User Authentication**: JWT-based authentication system
- **Todo Management**: Create, list, complete, and delete todos
- **Collaborative Completion**: Todos require both users to mark as completed
- **Filtering**: Filter todos by status (pending/completed) and creator
- **API Documentation**: Complete Swagger/OpenAPI documentation
- **PostgreSQL Database**: Robust database with proper migrations
- **Docker Support**: Full containerization with Docker Compose
- **Health Checks**: Built-in health monitoring

## 🏗️ Architecture

The project follows a clean architecture pattern with the following structure:

```
├── cmd/server/          # Application entry points
├── internal/
│   ├── handlers/        # HTTP handlers (controllers)
│   ├── services/        # Business logic layer
│   ├── repository/      # Data access layer
│   └── models/          # Data models and DTOs
├── pkg/
│   ├── auth/           # Authentication utilities
│   ├── database/       # Database connection
│   └── response/       # HTTP response utilities
├── middleware/         # HTTP middleware
├── routes/            # Route definitions
├── config/            # Configuration management
├── database/          # Database migrations and seeds
├── docs/              # Swagger documentation
└── scripts/           # Database initialization scripts
```

## 🛠️ Tech Stack

- **Language**: Go 1.23.1
- **Database**: PostgreSQL 17.4
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **HTTP Router**: Gorilla Mux

## 📋 Prerequisites

- Go 1.23.1 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd LoveApp-Backend
   ```

2. **Start the application**
   ```bash
   make docker-run
   # or
   docker-compose up -d
   ```

3. **Access the application**
   - API: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/index.html
   - Health Check: http://localhost:8080/health

### Local Development

1. **Install dependencies**
   ```bash
   make deps
   ```

2. **Set up environment variables**
   ```bash
   cp example.env .env
   # Edit .env with your configuration
   ```

3. **Start PostgreSQL** (using Docker)
   ```bash
   docker-compose up -d postgres
   ```

4. **Run the application**
   ```bash
   make run
   ```

## 🔧 Configuration

Environment variables can be set in `.env` file:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=loveapp
DB_PASSWORD=loveapp123
DB_NAME=loveapp
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-please

# Server Configuration
SERVER_PORT=8080

# Environment
ENV=development
```

## 📚 API Documentation

### Authentication

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "username": "anyel",
  "password": "password"
}
```

**Response:**
```json
{
  "message": "Inicio de sesión exitoso",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "anyel",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Change Password
```http
POST /auth/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "password",
  "new_password": "newpassword"
}
```

### Todos

#### Create Todo
```http
POST /todos
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Comprar comida",
  "description": "Ir al supermercado y comprar frutas y verduras"
}
```

#### List Todos
```http
GET /todos?status=pending&creator_id=1
Authorization: Bearer <token>
```

**Query Parameters:**
- `status`: `all`, `pending`, `completed`
- `creator_id`: Filter by creator user ID

#### Complete Todo
```http
POST /todos/{id}/complete
Authorization: Bearer <token>
```

#### Delete Todo
```http
DELETE /todos/{id}
Authorization: Bearer <token>
```

## 👥 Default Users

The application comes with two pre-seeded users:

| Username | Password | Role |
|----------|----------|------|
| anyel    | password | User |
| alexis   | password | User |

## 🐳 Docker Commands

```bash
# Start all services
make docker-run

# Start with pgAdmin
make docker-run-with-tools

# Stop services
make docker-stop

# View logs
make logs

# Clean up
make docker-clean
```

## 🔨 Development Commands

```bash
# Install dependencies
make deps

# Generate Swagger docs
make swagger

# Build application
make build

# Run tests
make test

# Run development server
make dev

# Check health
make health

# Open Swagger UI
make swagger-ui
```

## 🗄️ Database

### Migrations

Database migrations are automatically run on application startup. The migrations create:

- `users` table with authentication data
- `todos` table with todo items
- Proper indexes for performance
- Triggers for automatic timestamp updates

### Seeding

The application automatically seeds the database with default users on startup.

## 📊 Monitoring

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Servicio funcionando correctamente"
}
```

### pgAdmin (Optional)

When running with tools profile, pgAdmin is available at:
- URL: http://localhost:5050
- Email: admin@loveapp.com
- Password: admin123

## 🔒 Security

- JWT tokens for authentication
- Password hashing with bcrypt
- SQL injection protection with parameterized queries
- CORS middleware for cross-origin requests
- Input validation and sanitization

## 🚀 Deployment

### Production Build

```bash
make prod-build
```

### Environment Variables for Production

Ensure to set secure values for:
- `JWT_SECRET`: Use a strong, random secret
- `DB_PASSWORD`: Use a secure database password
- `DB_SSLMODE`: Set to `require` for production

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License.

## 🆘 Support

For support, please open an issue in the repository or contact the development team.

---

Made with ❤️ for couples who want to organize their lives together!