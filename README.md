# Go Fiber Starter Kit

A production-ready Go REST API starter kit built with Fiber framework. Includes authentication (JWT + Google OAuth), database management, pagination, and more.

## Features

- 🔐 **Authentication System**
  - JWT-based authentication
  - Google OAuth integration
  - Token blacklisting for secure logout
  - Role-based access control

- 🗄️ **Database**
  - GORM ORM with SQLite (development) and PostgreSQL (production) support
  - Automatic migrations
  - Seeding with superadmin user

- 📝 **API Features**
  - RESTful API design
  - Request validation
  - Pagination support
  - Structured JSON responses

- 🛡️ **Security**
  - Password hashing with bcrypt
  - CORS configuration
  - JWT middleware protection
  - Environment variable configuration

- 📚 **Documentation**
  - Swagger/OpenAPI documentation
  - Auto-generated API docs

## Tech Stack

- **Framework**: [Fiber](https://gofiber.io/)
- **Database**: SQLite (development), PostgreSQL (production)
- **ORM**: [GORM](https://gorm.io/)
- **Authentication**: JWT, Google OAuth
- **Validation**: Custom validator
- **Documentation**: Swagger

## Project Structure

```
gofiber_starter_kit/
├── blacklist/          # Token blacklist management
├── cmd/               # Application entry point
├── config/            # Configuration management
├── controllers/       # HTTP controllers
├── database/          # Database connection and operations
├── docs/              # Swagger documentation
├── middlewares/       # Custom middleware (JWT, roles)
├── models/            # Database models
├── requests/          # Request structs
├── responses/         # Response structs
├── routes/            # Route definitions
├── services/          # Business logic
└── utils/             # Utility functions
```

## Getting Started

### Prerequisites

- Go 1.23.4 or higher
- SQLite (for development)
- PostgreSQL (for production, optional)

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd gofiber_starter_kit
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:
```env
DB_PATH=gofiber.db
JWT_SECRET=your-super-secret-jwt-key
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8000/api/auth/google/callback
```

4. Run the application:
```bash
# With air for live reload (recommended for development)
air

# Or directly with go
go run cmd/main.go
```

The server will start on `http://localhost:8000`

### Seeded Superadmin

A superadmin user is automatically created with:
- Email: `admin@example.com`
- Password: `admin1234`

## API Documentation

Once the server is running, access the Swagger documentation at:
- Swagger UI: `http://localhost:8000/swagger/index.html`

## Available Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/google` - Google OAuth authentication
- `POST /api/auth/logout` - User logout

### User
- `GET /api/user/profile` - Get user profile (protected)

## Development

### Using Air for Live Reload

The project includes [Air](https://github.com/cosmtrek/air) for live reload during development. Configuration is in `.air.toml`.

### Database Migrations

Migrations are handled automatically by GORM's AutoMigrate feature.

### Adding New Features

1. Create models in `models/` directory
2. Add request/response structs in their respective directories
3. Implement business logic in `services/`
4. Create controllers in `controllers/`
5. Define routes in `routes/routes.go`

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_PATH | Database file path | gofiber.db |
| JWT_SECRET | Secret key for JWT tokens | secret |
| GOOGLE_CLIENT_ID | Google OAuth client ID | - |
| GOOGLE_CLIENT_SECRET | Google OAuth client secret | - |
| GOOGLE_REDIRECT_URL | Google OAuth redirect URL | - |

## Deployment

### Building for Production

```bash
go build -o build/main ./cmd/main.go
```

### Docker Deployment

A sample Dockerfile might look like:

```dockerfile
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8000

CMD ["./main"]
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

If you have any questions or issues, please open an issue on the GitHub repository.
```
