# Contributing to SilentCast

Thank you for your interest in contributing to SilentCast! This guide will help you get started.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- **Be respectful** and inclusive
- **Be patient** with new contributors
- **Be constructive** in feedback
- **Be collaborative** and helpful

## How to Contribute

### Reporting Issues

1. **Search existing issues** to avoid duplicates
2. **Use issue templates** when available
3. **Provide details**:
   - SilentCast version
   - Operating system
   - Steps to reproduce
   - Expected vs actual behavior
   - Error messages and logs

### Suggesting Features

1. **Check the roadmap** first
2. **Open a discussion** before implementation
3. **Explain the use case** clearly
4. **Consider alternatives**

### Contributing Code

1. **Fork the repository**
2. **Create a feature branch**
3. **Make your changes**
4. **Test thoroughly**
5. **Submit a pull request**

## Development Setup

### Prerequisites

```bash
# Install Go 1.23+
# Install Git
# Install Make (optional)
```

### Fork and Clone

```bash
# Fork on GitHub, then:
git clone https://github.com/YOUR-USERNAME/silentcast.git
cd silentcast
git remote add upstream https://github.com/SphereStacking/silentcast.git
```

### Development Workflow

```bash
# Create feature branch
git checkout -b feature/your-feature

# Make changes
# ...

# Run tests
make test

# Commit with conventional commit message
git commit -m "feat: add new spell type"

# Push to your fork
git push origin feature/your-feature

# Create pull request on GitHub
```

## Coding Standards

### Go Code Style

Follow standard Go conventions:

```go
// Package comment
package action

import (
    "context"
    "fmt"
    
    "github.com/SphereStacking/silentcast/pkg/logger"
)

// Executor handles action execution.
type Executor struct {
    logger logger.Logger
}

// Execute runs the specified action.
func (e *Executor) Execute(ctx context.Context, action Action) error {
    // Implementation
    return nil
}
```

### Code Organization

- **One package per directory**
- **Interface definitions** at the top
- **Public types** before private
- **Methods** grouped by receiver

### Error Handling

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to execute action: %w", err)
}

// Custom error types for specific cases
type NotFoundError struct {
    Spell string
}

func (e NotFoundError) Error() string {
    return fmt.Sprintf("spell not found: %s", e.Spell)
}
```

### Testing

Write table-driven tests:

```go
func TestParser_Parse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Hotkey
        wantErr bool
    }{
        {
            name:  "simple key",
            input: "ctrl+a",
            want:  Hotkey{Modifiers: []string{"ctrl"}, Key: "a"},
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Parse() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Commit Messages

Use conventional commits:

```
type(scope): subject

body

footer
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `chore`: Build/tooling changes

### Examples

```
feat(hotkey): add support for mouse buttons

- Add mouse button parsing
- Update key validation
- Add tests for mouse events

Closes #123
```

```
fix(config): handle empty yaml files gracefully

Previously crashed when config file was empty.
Now returns default configuration.
```

## Pull Request Process

### Before Submitting

1. **Update from upstream**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all checks**:
   ```bash
   make fmt
   make lint
   make test
   ```

3. **Update documentation** if needed

### PR Description

Include:
- **What** changes were made
- **Why** they were needed
- **How** to test them
- **Breaking changes** if any

### Review Process

1. **Automated checks** must pass
2. **Code review** by maintainers
3. **Testing** on multiple platforms
4. **Documentation** review
5. **Merge** when approved

## Testing Guidelines

### Unit Tests

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/config/...

# Run with coverage
make test-coverage
```

### Integration Tests

```bash
# Run integration tests
make test-integration

# Run on specific platform
GOOS=windows make test-integration
```

### Manual Testing

Test on all supported platforms:
- [ ] Windows 10/11
- [ ] macOS 12+
- [ ] Linux (Ubuntu/Fedora)

## Documentation

### Code Documentation

```go
// Package logger provides structured logging with rotation support.
//
// The logger supports multiple output targets, log levels, and
// automatic rotation based on size and age.
package logger

// Logger is the main logging interface.
// It provides methods for different log levels and structured logging.
type Logger interface {
    // Debug logs a debug message with optional fields.
    Debug(msg string, fields ...Field)
    
    // Info logs an informational message.
    Info(msg string, fields ...Field)
}
```

### User Documentation

- Update relevant `.md` files in `docs/`
- Include examples
- Update configuration references
- Add to FAQ if applicable

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features
- **PATCH**: Bug fixes

### Release Checklist

1. [ ] Update version in `version.go`
2. [ ] Update CHANGELOG.md
3. [ ] Run full test suite
4. [ ] Build all platforms
5. [ ] Create GitHub release
6. [ ] Update documentation
7. [ ] Announce release

## Getting Help

### For Contributors

- **Development questions**: Open a GitHub Discussion
- **Bug in your PR**: Update the PR with fixes
- **Stuck on something**: Ask in the PR comments

### Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Recognition

Contributors are recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to SilentCast! 🎉