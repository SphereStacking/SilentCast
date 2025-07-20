# TDD Best Practices for SilentCast

This document captures proven TDD patterns and best practices specific to the SilentCast project based on our implementation experience.

## 🎯 Project-Specific TDD Patterns

### 1. Testing Platform-Specific Code

#### Pattern: Abstract Platform Behavior
```go
// ✅ Good: Test the interface, not the implementation
func TestNotifier_SendsNotificationSuccessfully(t *testing.T) {
    notifier := notify.NewMockNotifier(true)
    
    err := notifier.Notify(context.Background(), notify.Notification{
        Title:   "Test",
        Message: "Message",
        Level:   notify.LevelInfo,
    })
    
    assert.NoError(t, err)
}
```

#### Pattern: Skip Platform-Dependent Tests
```go
func TestHotkeyManager_RegisterHotkey(t *testing.T) {
    if runtime.GOOS != "darwin" {
        t.Skip("macOS-specific test")
    }
    
    // Platform-specific test logic
}
```

### 2. Testing Configuration System

#### Pattern: Cascade Configuration Testing
```go
func TestConfigLoader_CascadeLoading(t *testing.T) {
    tempDir := t.TempDir()
    
    // TDD approach: Start with one config file
    baseConfig := `
spells:
  e: editor
`
    
    // Add platform-specific override
    macosConfig := `
spells:
  e: vscode-mac
`
    
    loader := config.NewLoader()
    cfg, err := loader.LoadCascade(tempDir)
    
    assert.NoError(t, err)
    assert.Equal(t, "vscode-mac", cfg.Spells["e"])
}
```

### 3. Testing Error Handling

#### Pattern: Unified Error Context Testing
```go
func TestActionExecutor_ErrorWithContext(t *testing.T) {
    executor := action.NewExecutor()
    
    err := executor.Execute(action.Config{
        Type: "invalid",
        App:  "nonexistent",
    })
    
    require.Error(t, err)
    
    var spellErr *errors.SpellbookError
    require.True(t, errors.As(err, &spellErr))
    assert.Equal(t, errors.ErrorTypeConfig, spellErr.Type)
    assert.Contains(t, spellErr.Context, "action_type")
}
```

### 4. Testing Hotkey System

#### Pattern: Mock Hotkey Events
```go
func TestHotkeyManager_TriggersCallback(t *testing.T) {
    manager := hotkey.NewMockManager()
    
    var triggered bool
    callback := func(spell string) {
        triggered = true
    }
    
    err := manager.Register("alt+space", "e", callback)
    require.NoError(t, err)
    
    // Simulate hotkey press
    manager.SimulateKeyPress("alt+space", "e")
    
    assert.True(t, triggered)
}
```

## 🔄 TDD Workflow Optimizations

### 1. Fast Feedback Loops

#### Use Build Tags for Testing
```go
//go:build !integration

package action

// Unit tests that run quickly
```

#### Parallel Test Execution
```go
func TestConfigLoader_MultipleFiles(t *testing.T) {
    t.Parallel() // Run in parallel with other tests
    
    // Test implementation
}
```

### 2. Effective Mocking Strategy

#### Interface Segregation for Testing
```go
// Small, focused interfaces are easier to mock
type Notifier interface {
    Notify(context.Context, Notification) error
}

type OutputNotifier interface {
    Notifier
    ShowWithOutput(context.Context, OutputNotification) error
}
```

#### Mock Builder Pattern
```go
func NewMockNotifier() *MockNotifier {
    return &MockNotifier{
        notifications: make([]Notification, 0),
        available:     true,
    }
}

func (m *MockNotifier) WithError(err error) *MockNotifier {
    m.simulateError = err
    return m
}

func (m *MockNotifier) WithAvailability(available bool) *MockNotifier {
    m.available = available
    return m
}
```

### 3. Test Data Management

#### Use Table-Driven Tests Effectively
```go
func TestSpellParser_ParseVariousFormats(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Spell
        wantErr bool
    }{
        // Start with simplest case in TDD
        {"single key", "e", Spell{Key: "e"}, false},
        
        // Add complexity incrementally
        {"sequence", "g,s", Spell{Sequence: []string{"g", "s"}}, false},
        
        // Edge cases discovered during TDD
        {"empty input", "", Spell{}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewSpellParser()
            got, err := parser.Parse(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

#### Fixture Management
```go
func createTestSpellbook(t *testing.T) string {
    content := `
spells:
  e: editor
  g,s: git-status
grimoire:
  editor:
    app: code
  git-status:
    script: git status
`
    
    file := filepath.Join(t.TempDir(), "spellbook.yml")
    err := os.WriteFile(file, []byte(content), 0644)
    require.NoError(t, err)
    
    return file
}
```

## 📊 Quality Metrics and Standards

### 1. Coverage Guidelines

#### Natural Coverage (90%+)
- Don't write tests for coverage numbers
- Focus on behavior and edge cases
- Use coverage to find missed scenarios

#### Coverage Exclusions
```go
//nolint:gocov // Platform-specific code tested manually
func (m *DarwinHotkeyManager) platformSpecificInit() {
    // Code that requires manual testing
}
```

### 2. Test Quality Metrics

#### Test Naming Convention
```go
// Pattern: TestSubject_Behavior_Context
func TestSpellParser_ReturnsError_WhenInputIsEmpty(t *testing.T) {}
func TestActionExecutor_LaunchesApp_WhenValidConfig(t *testing.T) {}
func TestConfigLoader_MergesConfigs_WithCascadeLoading(t *testing.T) {}
```

#### Test Organization
```
internal/action/
├── executor.go
├── executor_test.go          // Unit tests
├── integration_test.go       // Integration tests
└── testdata/
    ├── valid_config.yml
    └── invalid_config.yml
```

### 3. TDD Cycle Metrics

#### Measuring Cycle Time
```bash
# Target metrics for SilentCast TDD
Red Phase:    2-3 minutes (write failing test)
Green Phase:  3-5 minutes (make test pass)  
Refactor:     2-4 minutes (improve code)
Total Cycle:  7-12 minutes
```

#### Quality Indicators
- Tests fail for right reasons
- Minimal code to pass tests
- Refactoring doesn't break tests
- High confidence in changes

## 🛠️ Tools and Automation

### 1. TDD-Friendly Make Targets

```makefile
# Fast TDD feedback
.PHONY: tdd-test
tdd-test:
	@go test -short -v ./...

# Watch mode for continuous testing
.PHONY: tdd-watch
tdd-watch:
	@echo "👀 TDD Watch mode - save files to run tests"
	@find . -name "*.go" | entr -c $(MAKE) tdd-test

# TDD cycle timing
.PHONY: tdd-cycle
tdd-cycle:
	@./scripts/tdd-metrics.sh cycle_start
```

### 2. CI/CD Integration

```yaml
# .github/workflows/tdd.yml
name: TDD Validation
on: [push, pull_request]

jobs:
  tdd-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.19
      
      - name: Verify TDD practices
        run: |
          # Check test coverage
          make test-coverage
          
          # Verify test quality
          make lint-tests
          
          # Check cycle time metrics
          make tdd-validate
```

### 3. IDE Configuration

#### VS Code Settings
```json
{
    "go.testTimeout": "30s",
    "go.testFlags": ["-v", "-short"],
    "go.coverageDecorator": {
        "type": "gutter"
    },
    "files.watcherExclude": {
        "**/coverage.out": true
    }
}
```

## 🎓 Training and Adoption

### 1. Onboarding New Developers

#### TDD Kata Sessions
```go
// Start with simple problems
func TestFizzBuzz_ReturnsFizz_WhenDivisibleBy3(t *testing.T) {
    result := FizzBuzz(3)
    assert.Equal(t, "Fizz", result)
}

// Progress to SilentCast features
func TestSpellParser_ParsesSingleKey(t *testing.T) {
    parser := NewSpellParser()
    result, err := parser.Parse("e")
    assert.NoError(t, err)
    assert.Equal(t, Spell{Key: "e"}, result)
}
```

#### Code Review Checklist
- [ ] Test written before implementation?
- [ ] Test fails for expected reason?
- [ ] Minimal implementation to pass?
- [ ] Refactoring maintains green tests?
- [ ] Test name describes behavior?

### 2. Common TDD Mistakes

#### Writing Tests After Implementation
```go
// ❌ Don't do this
func (s *SpellParser) Parse(input string) (Spell, error) {
    // Implementation first
    return Spell{Key: input}, nil
}

// Then write test to match implementation
```

#### Over-Engineering in Green Phase
```go
// ❌ Too complex for green phase
func (s *SpellParser) Parse(input string) (Spell, error) {
    // Complex regex parsing, caching, validation
    // Save for refactor phase!
}

// ✅ Minimal implementation for green
func (s *SpellParser) Parse(input string) (Spell, error) {
    if input == "e" {
        return Spell{Key: "e"}, nil
    }
    return Spell{}, errors.New("not implemented")
}
```

## 📈 Continuous Improvement

### 1. TDD Retrospectives

#### Weekly Team Review
- What TDD practices worked well?
- Where did we struggle with TDD?
- How can we improve cycle times?
- What tools would help us?

#### Metrics Dashboard
```bash
# Generate weekly TDD report
make tdd-report-weekly

# Show trends over time
./scripts/tdd-metrics.sh trends
```

### 2. Evolving Practices

#### Experiment with New Patterns
- Property-based testing for complex inputs
- Mutation testing for test quality
- Approval testing for legacy code

#### Share Learnings
- Document new patterns in this file
- Update CLAUDE.md with insights
- Create example implementations

## 🎯 Success Criteria

### Project-Level Success
- [ ] 90%+ test coverage achieved naturally
- [ ] New features always use TDD
- [ ] Refactoring is safe and confident
- [ ] Bug rate decreases over time

### Team-Level Success
- [ ] All developers comfortable with TDD
- [ ] Average cycle time under 10 minutes
- [ ] Tests serve as living documentation
- [ ] Code reviews focus on design

### Technical Success
- [ ] Fast test suite (under 30 seconds)
- [ ] High-quality test names and structure
- [ ] Effective mocking strategy
- [ ] CI/CD validates TDD practices

---

## 📝 Contributing to This Document

This document should evolve with our TDD practice. When you discover new patterns or best practices:

1. Add them to the relevant section
2. Include code examples
3. Explain why the pattern helps
4. Share with the team for feedback

Remember: The goal is to make TDD natural and effective for SilentCast development.