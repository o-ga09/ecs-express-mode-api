---
description: バックエンド開発のためのGo言語に関する指示
applyTo: "**.go,**_test.go,go.mod,go.sum"
---

# Go API Project - Coding Agent Instructions

## Project Overview

This is a Go-based API project. Follow these instructions carefully when generating, modifying, or reviewing code.

---

## Coding Standards

### General Guidelines

- **Go Version**: Use Go 1.25
- **Error Handling**: Always handle errors explicitly. Never ignore errors with `_`
- **Naming Conventions**:
  - Use `camelCase` for private functions/variables
  - Use `PascalCase` for exported functions/variables
  - Use descriptive names; avoid single-letter variables except in short scopes (e.g., loop indices)
- **Comments**:
  - Add package comments at the top of each package
  - Document all exported functions, types, and constants
  - Use `//` for single-line comments, not `/* */`

### Project Structure

```
project-root/
├── cmd/                    # Application entrypoints
│   └── api/
│       └── main.go
├── internal/               # Private application code
│   ├── handler/           # HTTP handlers
│   ├── usecase/           # Business logic
│   ├── repository/        # Data access layer
│   ├── domain/            # Domain models
│   ├── server/            # init server, request / response type
│   └── middleware/        # Middleware components
├── pkg/                    # Public libraries
├── docs/                   # Documentation (see docs section)
├── config/                 # Configuration files
├── migrations/             # Database migrations
├── scripts/                # Build and utility scripts
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Code Organization

- **Separation of Concerns**: Keep handlers, usecases, and repositories separated
- **Dependency Injection**: Use constructor functions that accept dependencies
- **Interface Usage**: Define interfaces in the consumer package, not the implementer
- **Context Propagation**: Always pass `context.Context` as the first parameter

### Error Handling

```go
// Good: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to fetch user %d: %w", userID, err)
}

// Bad: Generic error messages
if err != nil {
    return err
}
```

### API Response Format

```go
// Success Response
type SuccessResponse struct {
    Data    interface{} `json:"data"`
    Message string      `json:"message,omitempty"`
}

// Error Response
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code"`
    Details string `json:"details,omitempty"`
}
```

---

## Lint & Format Standards

### Required Tools

- **gofmt**: Standard Go formatter
- **goimports**: Import organization
- **golangci-lint**: Comprehensive linter (v2.0+)

### Pre-commit Checks

Run these commands before committing:

```bash
# Format code
gofmt -w .
goimports -w .

# Run linters
golangci-lint run ./...

# Run tests
go test ./... -race -cover
```

### golangci-lint Configuration

Create `.golangci.yml` in project root:

```yaml
linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - unused
    - ineffassign
    - typecheck
    - bodyclose
    - errname
    - errorlint
    - exhaustive
    - goconst
    - gocritic
    - revive
    - misspell

linters-settings:
  goimports:
    local-prefixes: github.com/your-org/your-project
  errcheck:
    check-blank: true
  goconst:
    min-len: 3
    min-occurrences: 3
  revive:
    rules:
      - name: exported
        severity: warning
      - name: unexported-return
        severity: warning

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

### Makefile Targets

```makefile
.PHONY: fmt lint test

fmt:
	@echo "Formatting code..."
	@gofmt -w .
	@goimports -w .

lint:
	@echo "Running linters..."
	@golangci-lint run ./...

test:
	@echo "Running tests..."
	@go test ./... -race -cover -coverprofile=coverage.out

all: fmt lint test
```

---

## Documentation Management

### docs/ Directory Structure

```
docs/
├── 00_getting_started/
│   ├── README.md
│   ├── setup.md
│   └── quick_start.md
├── 01_architecture/
│   ├── README.md
│   ├── overview.md
│   ├── database_schema.md
│   └── api_design.md
├── 02_api_reference/
│   ├── README.md
│   ├── authentication.md
│   ├── users.md
│   └── resources.md
├── 03_development/
│   ├── README.md
│   ├── coding_standards.md
│   ├── testing_guide.md
│   └── deployment.md
├── 04_operations/
│   ├── README.md
│   ├── monitoring.md
│   └── troubleshooting.md
└── 99_misc/
    ├── README.md
    └── changelog.md
```

### Documentation Rules

#### Naming Convention

- Use `00_`, `01_`, `02_`, etc. as directory prefixes for ordering
- Use lowercase with underscores for directory and file names
- Each directory must contain a `README.md` as the index

#### Documentation Content Guidelines

**00_getting_started/**:

- Installation instructions
- Environment setup
- Quick start guide
- Prerequisites

**01_architecture/**:

- System architecture diagrams
- Component relationships
- Database schemas (ERD)
- Design decisions and ADRs (Architecture Decision Records)

**02_api_reference/**:

- API endpoint documentation
- Request/response examples
- Authentication methods
- Error codes reference

**03_development/**:

- Development workflow
- Testing strategies
- Code review process
- Deployment procedures

**04_operations/**:

- Monitoring and alerting
- Backup procedures
- Incident response
- Performance tuning

**99_misc/**:

- Changelog
- FAQ
- Known issues
- Legacy documentation

#### Markdown Standards

- Use ATX-style headers (`#` notation)
- Include a table of contents for documents over 200 lines
- Use code fences with language identifiers
- Include examples for all API endpoints
- Keep line length under 100 characters for readability

#### API Documentation Example

````markdown
## GET /api/v1/users/:id

Retrieves a user by ID.

### Parameters

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| id   | int  | Yes      | User ID     |

### Response

**Success (200 OK)**

```json
{
  "data": {
    "id": 123,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```
````

**Error (404 Not Found)**

```json
{
  "error": "User not found",
  "code": "USER_NOT_FOUND"
}
```

````

### Documentation Maintenance
- Update docs when adding/modifying APIs
- Run spell check before committing
- Review docs during code review
- Keep examples up-to-date with actual implementation
- Archive outdated docs to `99_misc/archive/`

---

## Testing Standards

### Test File Organization
- Place tests in `*_test.go` files alongside the code
- **Use table-driven tests** for all test cases with multiple scenarios
- Mock external dependencies using interfaces

### Test Coverage
- Maintain minimum 80% code coverage
- Focus on business logic and error paths
- Use `go test -cover` to verify coverage

### Table-Driven Test Pattern
**Always use table-driven tests** for testing multiple scenarios:

```go
func TestUserUsecase_GetUser(t *testing.T) {
    tests := []struct {
        name    string
        userID  int
        want    *domain.User
        wantErr bool
    }{
        {
            name:    "valid user",
            userID:  1,
            want:    &domain.User{ID: 1, Name: "Test"},
            wantErr: false,
        },
        {
            name:    "user not found",
            userID:  999,
            want:    nil,
            wantErr: true,
        },
        {
            name:    "invalid user ID",
            userID:  -1,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := &MockUserRepository{}
            usecase := NewUserUsecase(mockRepo)

            // Act
            got, err := usecase.GetUser(context.Background(), tt.userID)

            // Assert
            if (err != nil) != tt.wantErr {
                t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetUser() = %v, want %v", got, tt.want)
            }
        })
    }
}
````

---

## Git Commit Standards

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

**Example**:

```
feat(api): add user authentication endpoint

- Implement JWT-based authentication
- Add middleware for token validation
- Update user model with password hash

Closes #123
```

---

## Security Guidelines

- Never commit secrets or API keys
- Use environment variables for configuration
- Validate all user inputs
- Use parameterized queries to prevent SQL injection
- Implement rate limiting on public endpoints
- Set appropriate CORS policies
- Use HTTPS in production

---

## Performance Guidelines

- Use connection pooling for databases
- Implement caching where appropriate (Redis)
- Add indexes to frequently queried database columns
- Use goroutines for concurrent operations (with proper synchronization)
- Profile code with `pprof` for optimization
- Set appropriate timeouts for external requests

---

## Code Review Checklist

- [ ] Code follows project structure
- [ ] All errors are properly handled
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] Linters pass without warnings
- [ ] No sensitive information in code
- [ ] Performance considerations addressed
- [ ] Security best practices followed

---

## Additional Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
