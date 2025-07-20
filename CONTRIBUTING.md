# Contributing to SilentCast

Thank you for your interest in contributing to SilentCast! We welcome contributions from everyone. This guide will help you get started.

## 🎯 Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- **Be respectful**: Treat everyone with respect
- **Be inclusive**: Welcome newcomers and help them get started
- **Be collaborative**: Work together to solve problems
- **Be professional**: Keep discussions focused and constructive

## 🚀 Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Set up development environment** - See [Development Setup Guide](docs/development/setup.md)
4. **Create a branch** for your changes
5. **Make your changes** following our guidelines
6. **Submit a pull request**

## 🎫 Using the Ticket System

SilentCast uses a ticket-based development system. Before starting work:

```bash
# Check available tickets
./ticket list --status todo

# View ticket details
./ticket show T001

# Claim a ticket
./ticket status T001 in_progress
```

### Creating New Tickets

For new features or bugs not covered by existing tickets:

```bash
# Create a new ticket
./ticket new --title "Add new feature" --type feature --priority medium
```

## 🏗️ Development Process

### 1. Before You Start

- Check existing [issues](https://github.com/SphereStacking/silentcast/issues) and [pull requests](https://github.com/SphereStacking/silentcast/pulls)
- Look for tickets marked `good-first-issue` or `help-wanted`
- Discuss major changes in an issue first

### 2. Development Workflow

```bash
# 1. Sync with upstream
git checkout main
git pull upstream main

# 2. Create feature branch
git checkout -b feature/your-feature-name

# 3. Make changes following TDD
# - Write tests first
# - Implement feature
# - Refactor

# 4. Commit changes
git add .
git commit -m "feat: add amazing feature"

# 5. Push to your fork
git push origin feature/your-feature-name
```

### 3. Pull Request Process

1. **Update documentation** for any changed functionality
2. **Add tests** for new features
3. **Run all checks** locally:
   ```bash
   make check  # Runs fmt, lint, and test
   ```
4. **Update CHANGELOG.md** if applicable
5. **Submit PR** with clear description

## 📝 Coding Standards

### Go Code Style

We follow standard Go conventions:

```go
// Package comment explains the purpose
package action

import (
    "context"
    "fmt"
    
    "github.com/SphereStacking/silentcast/internal/config"
)

// Executor handles action execution.
// It implements the ActionExecutor interface.
type Executor struct {
    config *config.Config
    logger Logger
}

// Execute runs the specified action.
// It returns an error if the action fails.
func (e *Executor) Execute(ctx context.Context, name string) error {
    action, err := e.config.GetAction(name)
    if err != nil {
        return fmt.Errorf("get action %q: %w", name, err)
    }
    
    // Implementation...
    return nil
}
```

### Best Practices

1. **Error Handling**:
   ```go
   // Wrap errors with context
   if err != nil {
       return fmt.Errorf("failed to execute %s: %w", name, err)
   }
   ```

2. **Table-Driven Tests**:
   ```go
   func TestExecutor(t *testing.T) {
       tests := []struct {
           name    string
           input   string
           want    string
           wantErr bool
       }{
           {"valid input", "test", "result", false},
           {"empty input", "", "", true},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // Test implementation
           })
       }
   }
   ```

3. **Comments**:
   - Export comments start with the name
   - Explain why, not what
   - Use complete sentences

### Magic Terminology

Always use our magic terminology consistently:

- ✅ **Spells** (not shortcuts or hotkeys)
- ✅ **Grimoire** (not actions or commands)
- ✅ **Spellbook** (not config or configuration file)

## 🧪 Test-Driven Development (TDD)

We follow strict TDD practices:

### The TDD Cycle

1. **🔴 Red**: Write a failing test
2. **🟢 Green**: Make it pass with minimal code
3. **🔵 Refactor**: Improve while keeping tests green

### Example TDD Session

```go
// 1. RED - Write failing test
func TestSpellParser_ParseSequence(t *testing.T) {
    parser := NewSpellParser()
    result, err := parser.Parse("g,s")
    
    assert.NoError(t, err)
    assert.Equal(t, []string{"g", "s"}, result.Sequence)
}

// 2. GREEN - Minimal implementation
func (p *SpellParser) Parse(input string) (*Spell, error) {
    if strings.Contains(input, ",") {
        parts := strings.Split(input, ",")
        return &Spell{Sequence: parts}, nil
    }
    return &Spell{Key: input}, nil
}

// 3. REFACTOR - Improve design
func (p *SpellParser) Parse(input string) (*Spell, error) {
    if p.isSequence(input) {
        return p.parseSequence(input)
    }
    return p.parseSingleKey(input)
}
```

### Test Requirements

- **Coverage**: Aim for >80% test coverage
- **Edge Cases**: Test error conditions and boundaries
- **Mocking**: Use interfaces for external dependencies
- **Clarity**: Tests should document behavior

## 📋 Commit Message Format

We use conventional commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process or auxiliary tool changes

### Examples

```bash
# Feature
git commit -m "feat(hotkey): add support for function keys"

# Bug fix
git commit -m "fix(config): resolve yaml parsing error for sequences"

# Documentation
git commit -m "docs(readme): update installation instructions"

# With body
git commit -m "feat(action): add elevated execution support

- Add platform-specific elevation methods
- Implement permission checking
- Add tests for Windows and macOS

Closes #123"
```

## 🔍 Code Review Process

### For Contributors

1. **Self-review** your PR first
2. **Respond** to review comments promptly
3. **Update** based on feedback
4. **Be patient** - reviewers are volunteers

### For Reviewers

1. **Be constructive** - suggest improvements
2. **Be specific** - point to examples
3. **Be timely** - review within 48 hours
4. **Be thorough** - check tests and docs

### Review Checklist

- [ ] Code follows project style
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No unnecessary dependencies added
- [ ] Performance impact considered
- [ ] Security implications reviewed

## 📚 Documentation

### When to Update Docs

Update documentation when you:
- Add new features
- Change existing behavior
- Add new configuration options
- Change CLI commands
- Fix confusing documentation

### Documentation Types

1. **Code Comments**: Explain why, not what
2. **API Documentation**: GoDoc comments
3. **User Guides**: Markdown in `/docs`
4. **Examples**: Working examples in `/examples`

## 🏷️ Issue and PR Labels

### Priority Labels
- `critical`: Urgent, breaking issues
- `high`: Important features/fixes
- `medium`: Normal priority
- `low`: Nice to have

### Type Labels
- `bug`: Something isn't working
- `feature`: New feature request
- `docs`: Documentation improvements
- `refactor`: Code improvement

### Status Labels
- `help-wanted`: Looking for contributors
- `good-first-issue`: Good for newcomers
- `blocked`: Waiting on something
- `in-progress`: Being worked on

## 🚀 Release Process

### Version Numbering

We use semantic versioning (MAJOR.MINOR.PATCH):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### Release Checklist

1. Update version in `version.go`
2. Update `CHANGELOG.md`
3. Create release PR
4. Tag release after merge
5. Build and upload binaries
6. Update documentation

## 💡 Tips for Contributors

### First Time Contributors

1. Start with `good-first-issue` labeled issues
2. Read existing code to understand patterns
3. Ask questions - we're here to help!
4. Small PRs are easier to review

### Making Good PRs

- **One feature per PR**: Keep PRs focused
- **Clear description**: Explain what and why
- **Screenshots**: For UI changes
- **Tests**: Always include tests
- **Small commits**: Logical, atomic changes

### Getting Help

- **GitHub Discussions**: For questions
- **Issues**: For bugs and features
- **PR Comments**: For code-specific questions

## 🙏 Recognition

We value all contributions:
- Code contributions
- Documentation improvements
- Bug reports
- Feature suggestions
- Helping others

Contributors are recognized in:
- Release notes
- Contributors file
- Project documentation

## 📜 License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT).

---

<div align="center">
  <p><strong>Thank you for contributing to SilentCast! 🎉</strong></p>
  <p>Your contributions make the magic happen ✨</p>
</div>