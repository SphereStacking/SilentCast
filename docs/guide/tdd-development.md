# TDD Development Guide

This guide provides comprehensive instructions for Test-Driven Development (TDD) in the SilentCast project.

## 📋 Quick Start

### Prerequisites
- Go 1.19+ installed
- Understanding of Go testing basics
- 10-minute timer (physical or app)

### TDD Commands
```bash
# Start TDD watch mode
make tdd-watch

# Start new TDD cycle
make tdd-cycle

# Measure TDD metrics
make tdd-metrics

# Complete TDD cycle with metrics
make tdd-metrics-complete
```

## 🔄 Red-Green-Refactor Cycle

### 🔴 RED Phase (2-3 minutes)
**Goal**: Write the smallest failing test

1. **Write ONE failing test**
   ```go
   func TestSpellParser_ParseSingleKey(t *testing.T) {
       parser := NewSpellParser()
       result, err := parser.Parse("e")
       
       // This will fail - implementation doesn't exist yet
       assert.NoError(t, err)
       assert.Equal(t, Spell{Key: "e"}, result)
   }
   ```

2. **Verify test fails**
   ```bash
   go test -v ./internal/spell/
   # Should show: FAIL
   ```

3. **Write minimal production code to compile**
   ```go
   type SpellParser struct{}
   func NewSpellParser() *SpellParser { return &SpellParser{} }
   func (p *SpellParser) Parse(input string) (Spell, error) {
       return Spell{}, errors.New("not implemented")
   }
   ```

### 🟢 GREEN Phase (3-5 minutes)
**Goal**: Make the test pass with minimum code

1. **Write simplest implementation**
   ```go
   func (p *SpellParser) Parse(input string) (Spell, error) {
       if input == "e" {
           return Spell{Key: "e"}, nil
       }
       return Spell{}, errors.New("not implemented")
   }
   ```

2. **Verify test passes**
   ```bash
   go test -v ./internal/spell/
   # Should show: PASS
   ```

3. **All tests must pass**
   ```bash
   make test  # Run full test suite
   ```

### ♻️ REFACTOR Phase (2-4 minutes)
**Goal**: Improve code while keeping tests green

1. **Improve without changing behavior**
   ```go
   func (p *SpellParser) Parse(input string) (Spell, error) {
       // More general implementation
       if len(input) == 1 && unicode.IsLetter(rune(input[0])) {
           return Spell{Key: input}, nil
       }
       return Spell{}, errors.New("invalid spell format")
   }
   ```

2. **Run tests after each change**
   ```bash
   go test -v ./internal/spell/
   ```

3. **Keep all tests passing**

## 📊 TDD Metrics and Quality

### Cycle Time Tracking
```bash
# Start timing a cycle
./scripts/tdd-metrics.sh cycle_start

# End and record cycle
./scripts/tdd-metrics.sh cycle_end "implemented spell parsing"
```

### Coverage Goals
- **Natural coverage**: 90%+ through TDD cycles
- **No forced coverage**: Don't write tests just for coverage
- **Focus on behavior**: Test what the code should do

### Quality Indicators
- **Fast feedback**: Tests run in < 5 seconds
- **One assertion per test**: Clear failure messages
- **Descriptive names**: Test name explains the scenario
- **Independent tests**: No test depends on another

## 🎯 TDD Patterns for SilentCast

### Testing Action Execution
```go
func TestActionExecutor_ExecuteSpell(t *testing.T) {
    tests := []struct {
        name     string
        spell    Spell
        want     ActionResult
        wantErr  bool
    }{
        {
            name:  "editor spell launches VS Code",
            spell: Spell{Key: "e"},
            want:  ActionResult{Success: true, App: "code"},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            executor := NewActionExecutor()
            got, err := executor.Execute(tt.spell)
            
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

### Testing Configuration Loading
```go
func TestConfigLoader_LoadSpellbook(t *testing.T) {
    tempDir := t.TempDir()
    configPath := filepath.Join(tempDir, "spellbook.yml")
    
    yamlContent := `
spells:
  e: editor
grimoire:
  editor:
    app: code
`
    
    err := os.WriteFile(configPath, []byte(yamlContent), 0644)
    require.NoError(t, err)
    
    loader := NewConfigLoader()
    config, err := loader.Load(configPath)
    
    assert.NoError(t, err)
    assert.Contains(t, config.Spells, "e")
    assert.Equal(t, "editor", config.Spells["e"])
}
```

### Testing Platform-Specific Code
```go
func TestHotkeyManager_RegisterSpell(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping hotkey test in short mode")
    }
    
    manager := NewHotkeyManager()
    defer manager.Cleanup()
    
    spell := Spell{
        Prefix: "alt+space",
        Key:    "e",
    }
    
    err := manager.Register(spell)
    
    // Platform-specific behavior
    if runtime.GOOS == "darwin" {
        // macOS requires accessibility permissions
        if err != nil && strings.Contains(err.Error(), "accessibility") {
            t.Skip("macOS accessibility permissions required")
        }
    }
    
    assert.NoError(t, err)
}
```

## 🛠️ TDD Tools and Setup

### IDE Integration

#### VS Code
Install recommended extensions:
- Go extension
- Test Explorer
- TDD cycle timer

#### Vim/Neovim
```vim
" Run tests on save
autocmd BufWritePost *.go !go test -v ./...

" Quick test commands
nnoremap <leader>t :!go test -v %:h<CR>
nnoremap <leader>T :!go test -v ./...<CR>
```

### Continuous Testing
```bash
# Use entr for automatic test running
find . -name "*.go" | entr -c go test -v ./...

# Or use the make target
make tdd-watch
```

## 🎓 TDD Best Practices

### Do's
- ✅ Write failing test first (RED)
- ✅ Write minimal code to pass (GREEN)
- ✅ Refactor with confidence (REFACTOR)
- ✅ Keep cycles under 10 minutes
- ✅ Commit after each green cycle
- ✅ Use descriptive test names
- ✅ Test behavior, not implementation

### Don'ts
- ❌ Write implementation before test
- ❌ Write multiple tests at once
- ❌ Skip the refactor phase
- ❌ Write tests for 100% coverage
- ❌ Test private methods directly
- ❌ Use TDD for simple getters/setters

### Common Anti-Patterns

#### Testing Implementation Details
```go
// ❌ Bad: Testing internal state
func TestSpellParser_setsInternalField(t *testing.T) {
    parser := NewSpellParser()
    parser.Parse("e")
    assert.Equal(t, "e", parser.lastParsed) // Don't test private fields
}

// ✅ Good: Testing behavior
func TestSpellParser_ParseSingleKey(t *testing.T) {
    parser := NewSpellParser()
    result, err := parser.Parse("e")
    assert.NoError(t, err)
    assert.Equal(t, Spell{Key: "e"}, result)
}
```

#### Overly Complex Tests
```go
// ❌ Bad: Testing too much at once
func TestCompleteWorkflow(t *testing.T) {
    // 50 lines of setup
    // Testing configuration, parsing, execution, notification
    // Multiple assertions
}

// ✅ Good: One behavior per test
func TestSpellParser_ParsesValidKey(t *testing.T) {
    // Single responsibility test
}
```

## 📈 Measuring TDD Success

### Metrics to Track
- **Cycle time**: Average time per Red-Green-Refactor cycle
- **Test coverage**: Should naturally reach 90%+
- **Defect rate**: Bugs found in production
- **Refactoring safety**: Confidence in making changes

### Weekly TDD Review
```bash
# Generate TDD metrics report
make tdd-report

# Review cycle times
./scripts/tdd-metrics.sh summary

# Check test quality
go test -cover ./...
```

### Team TDD Maturity

#### Level 1: Basic TDD
- Writes tests before code
- Follows Red-Green-Refactor
- Achieves good coverage

#### Level 2: Proficient TDD
- Fast, confident cycles (< 10 minutes)
- Tests drive design decisions
- Comfortable refactoring

#### Level 3: Expert TDD
- Uses TDD for complex problems
- Teaches TDD to others
- Contributes to TDD practices

## 🚀 Next Steps

1. **Practice with Katas**: Implement classic problems (FizzBuzz, Roman Numerals)
2. **Apply to Real Features**: Use TDD for next SilentCast feature
3. **Measure and Improve**: Track metrics and optimize cycle time
4. **Share Knowledge**: Document lessons learned

## 📚 Resources

- [Test-Driven Development by Example](https://www.amazon.com/Test-Driven-Development-Kent-Beck/dp/0321146530) - Kent Beck
- [Growing Object-Oriented Software, Guided by Tests](https://www.amazon.com/Growing-Object-Oriented-Software-Guided-Tests/dp/0321503627) - Freeman & Pryce
- [t-wada's TDD cycle](https://github.com/twada-js/tdd-cycle)
- [TDD Best Practices](https://docs.microsoft.com/en-us/dotnet/core/testing/unit-testing-best-practices)

## 🤔 Troubleshooting

### "I can't think of a failing test"
- Start with the simplest case
- Think about edge cases
- Consider error conditions

### "My test is too complex"
- Break it into smaller tests
- Test one behavior at a time
- Use helper functions for setup

### "I don't know what to refactor"
- Look for duplication
- Improve naming
- Extract methods
- Simplify conditionals

### "Tests are slow"
- Use mocks for external dependencies
- Minimize file I/O
- Avoid network calls
- Use `testing.Short()` for integration tests