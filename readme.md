# Go REST API Design

A learning project for building a structured REST API in Go using **Chi**, **MySQL**, and manual **Dependency Injection**.

The main goal of this repository is to understand how a Go backend can be separated into clear layers:

```text
HTTP Request
     ↓
Router
     ↓
Controller
     ↓
Service
     ↓
Repository
     ↓
MySQL
```

The project currently implements the first user-registration flow and is intended to grow with more CRUD operations, validation, migrations, and additional services.

## Tech Stack

- Go
- Chi Router
- MySQL
- `database/sql`
- `go-sql-driver/mysql`
- `godotenv`

## Project Structure

```text
rest_api_design/
├── app/
│   └── applicatio.go          # Application configuration and dependency wiring
├── config/
│   ├── env.go                 # Environment variable loader
│   └── db/
│       └── mysql.config.go    # MySQL connection setup
├── controller/
│   ├── ping.go                # Basic test handlers
│   └── user.controller.go     # User HTTP handler
├── db/
│   └── repos/
│       ├── storage.go         # Repository container (work in progress)
│       └── user.go            # User model + UserRepository
├── router/
│   ├── router.go              # Main Chi router setup
│   └── user.router.go         # User routes
├── services/
│   └── user.service.go        # User business/service layer
├── main.go                    # Application entry point
├── go.mod
├── go.sum
└── .gitignore
```

## Dependency Injection

Dependencies are created from the outside and passed to the layer that needs them.

Current wiring in `Application.Run()`:

```go
dbConn, err := dbconfig.SetUpDB()
if err != nil {
    return err
}

userRepo := repository.NewUserRepository(dbConn)
userService := services.NewUserService(userRepo)
userController := controller.NewUserController(userService)
userRouter := router.NewUserRouter(userController)
```

This produces the dependency chain:

```text
*sql.DB
   ↓
UserRepository
   ↓
UserService
   ↓
UserController
   ↓
UserRouter
```

A layer does not create its own lower-level dependency. For example, `UserService` receives a `UserRepository` instead of opening a database connection itself.

## Repository Layer

The repository layer is responsible for database operations.

Current interface:

```go
type UserRepository interface {
    Create(u *User) (int64, error)
}
```

The implementation receives `*sql.DB` through its constructor:

```go
func NewUserRepository(db *sql.DB) UserRepository {
    return &UserRepositoryImpl{
        db: db,
    }
}
```

The current `Create` operation executes an SQL `INSERT` and returns the inserted ID.

As the project grows, the repository can contain operations such as:

```go
type UserRepository interface {
    Create(user *User) (int64, error)
    FindByID(id int64) (*User, error)
    FindByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id int64) error
}
```

SQL should remain in the repository layer instead of being placed in controllers or services.

## Service Layer

The service layer sits between the controller and repository.

```go
type UserService interface {
    CreateUser(u *db.User) (int64, error)
}
```

`UserServiceImpl` receives the repository through constructor injection:

```go
func NewUserService(userRepository db.UserRepository) UserService {
    return &UserServiceImpl{
        userRepository: userRepository,
    }
}
```

This layer is the right place for application/business logic such as:

- checking whether an email already exists
- hashing passwords
- deciding whether an operation is allowed
- coordinating multiple repositories
- calling email, cache, payment, or other services

## Controller Layer

Controllers handle HTTP-specific work and call services.

The current user controller receives `UserService`:

```go
type UserController struct {
    userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}
```

As the project develops, controllers should typically:

1. Decode request JSON.
2. Validate request data.
3. Call the appropriate service.
4. Convert service results/errors into HTTP responses.

Business logic and SQL should not be placed directly in controllers.

## Routing

The project uses Chi.

Each feature router implements the common `Router` interface:

```go
type Router interface {
    Register(r chi.Router)
}
```

Current user route:

```text
POST /signup
```

It calls:

```go
u.userController.RegisterUser
```

For larger features, resource-based REST routes are recommended:

```text
GET     /users
POST    /users
GET     /users/{id}
PUT     /users/{id}
DELETE  /users/{id}

GET     /books
POST    /books
GET     /books/{id}
PUT     /books/{id}
DELETE  /books/{id}
```

Prefer HTTP methods to describe actions rather than routes such as `/book/get` or `/book/delete/1`.

### Scaling Router Setup

When more feature routers are added, `SetUpRouter` can be changed to accept a variadic list:

```go
func SetUpRouter(routes ...Router) *chi.Mux {
    r := chi.NewRouter()

    for _, route := range routes {
        route.Register(r)
    }

    return r
}
```

Then application wiring can become:

```go
r := router.SetUpRouter(
    userRouter,
    bookRouter,
    bookingRouter,
)
```

Each feature remains responsible for registering its own routes.

## Validation

Request validation should happen near the HTTP/controller boundary.

A future request DTO could look like:

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

A useful separation is:

```text
Controller  → request/format validation
Service     → business-rule validation
Repository  → persistence
Database    → UNIQUE, NOT NULL, FOREIGN KEY, etc.
```

For example, email syntax belongs to request validation, while checking whether an email is already registered belongs to the service/database flow.

## Environment Configuration

The application reads configuration from a `.env` file.

Create `.env` in the project root:

```env
ADDR=:8080

DB_USER=root
DB_PASS=your_password
DB_HOST=127.0.0.1
DB_NAME=user
```

`.env` is ignored by Git and should not be committed.

## Database

The current repository expects a MySQL table named `user` with fields matching the current insert query:

```sql
CREATE TABLE user (
    id INT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    createdAt DATETIME NOT NULL
);
```

The current code explicitly supplies `id` during insertion. A later improvement is to use an auto-incrementing ID and let MySQL generate it.

Example:

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL
);
```

## Running the Project

### 1. Clone the repository

```bash
git clone https://github.com/srikanta0427/rest_api_design.git
cd rest_api_design
```

### 2. Download dependencies

```bash
go mod download
```

### 3. Configure MySQL

Create the database/table and add the required database values to `.env`.

### 4. Run the server

```bash
go run .
```

By default the server listens on:

```text
http://localhost:8080
```

unless `ADDR` is changed in `.env`.

## Current Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/` | Basic test endpoint |
| `GET` | `/{id}` | Returns the URL ID parameter |
| `POST` | `/signup` | Runs the current user creation flow |

> The `/signup` handler currently creates a hard-coded user in the controller. JSON request decoding is a planned next step.

## Planned Improvements

- Decode user data from request JSON
- Add request DTOs
- Add `go-playground/validator`
- Hash passwords before storage
- Use auto-increment IDs
- Add `FindByID`, `FindByEmail`, `Update`, and `Delete`
- Add proper JSON responses and centralized error handling
- Make `SetUpRouter` support multiple feature routers
- Add Book and other feature modules
- Add database migrations with Goose
- Add graceful shutdown
- Add middleware
- Add tests and repository/service mocks
- Add external services such as email through interfaces and dependency injection

## Future Architecture

As more features are introduced, the same pattern can be repeated:

```text
UserRouter
   ↓
UserController
   ↓
UserService ──────────→ EmailService
   ↓
UserRepository
   ↓
MySQL

BookRouter
   ↓
BookController
   ↓
BookService
   ↓
BookRepository
   ↓
MySQL
```

This keeps routing, HTTP handling, business logic, and persistence responsibilities separated while allowing dependencies to be replaced or mocked during testing.

## Purpose

This repository is primarily for learning and practicing production-style Go REST API design, including:

- interfaces
- dependency injection
- layered architecture
- REST routing
- SQL repositories
- validation
- database migrations
- service composition
- testing-friendly design
