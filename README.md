# Learning Go - Clean Architecture Example

A learning project demonstrating clean architecture principles in Go.

## Project Structure

```
.
├── cmd/              # Application entry points
│   └── main.go       # Main application
├── internal/         # Private application code
│   ├── domain/       # Core business logic and entities
│   ├── usecase/      # Business rules and use cases
│   ├── interface/    # External interfaces (HTTP handlers, repos)
│   └── infrastructure/ # External services and drivers
├── go.mod            # Go module file
├── go.sum            # Go dependencies lock file
└── . air.toml         # Hot reload configuration
```

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Air (for hot reloading)

### Installation

```bash
# Clone the repository
git clone https://github.com/QuangNV23062004/learning-go.git
cd learning-go

# Install dependencies
go mod download

# Install Air for hot reloading (optional)
go install github.com/cosmtrek/air@latest
```

### Running the Application

**With Air (hot reload)**:
```bash
air
```

**Without Air**:
```bash
go run ./cmd
```

## Architecture Overview

This project follows the **Clean Architecture** principles:

- **Domain Layer**: Business entities and rules (independent of frameworks)
- **Use Case Layer**: Application-specific business rules
- **Interface Layer**: Controllers, gateways, and presenters
- **Infrastructure Layer**: External services, databases, HTTP clients

## Development

### Project Structure Benefits

- **Testability**: Each layer can be tested independently
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Easy to swap implementations
- **Scalability**:  Simple to add new features without affecting existing code

## Dependencies

See `go.mod` for the complete list of dependencies.

## Contributing

Feel free to fork this project and submit pull requests! 

## License

This project is open source and available under the MIT License. 
