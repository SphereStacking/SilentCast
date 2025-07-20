# Error Handling Guidelines

This document defines the unified error handling patterns for the SilentCast codebase.

## 🎯 Overview

SilentCast uses a structured error handling system based on the `SpellbookError` type, which provides:
- **Type classification** for different error categories
- **Context information** for debugging and user feedback
- **Error wrapping** for error chain traceability
- **User-friendly messages** for end-user communication
- **Structured logging** integration

## 📋 Error Types

### Core Error Types

| Type | Purpose | Usage |
|------|---------|-------|
| `ErrorTypeConfig` | Configuration file issues | Invalid YAML, missing files, validation errors |
| `ErrorTypePermission` | Access control issues | File permissions, OS permissions, elevated actions |
| `ErrorTypeHotkey` | Hotkey system problems | Registration failures, conflicts, platform issues |
| `ErrorTypeExecution` | Action execution failures | Command failures, script errors, app launch issues |
| `ErrorTypeSystem` | System-level problems | OS interactions, resource availability |
| `ErrorTypeValidation` | Input validation errors | Parameter validation, format checking |
| `ErrorTypeNetwork` | Network and HTTP errors | Connection failures, timeouts, HTTP errors |
| `ErrorTypeIO` | File system operations | File not found, read/write errors, permissions |
| `ErrorTypePlatform` | Platform-specific issues | OS-specific functionality, compatibility |
| `ErrorTypeTimeout` | Operation timeouts | Long-running operations, deadlocks |
| `ErrorTypeNotFound` | Resource not found | Missing files, applications, configurations |

## 🔧 Usage Patterns

### 1. Creating New Errors

#### Simple Error Creation
```go
import "github.com/SphereStacking/silentcast/app/internal/errors"

// Simple error
return errors.New(errors.ErrorTypeConfig, "configuration file is empty")

// With context
return errors.New(errors.ErrorTypeConfig, "spell not found").
    WithContext("spell_name", spellName).
    WithContext("available_spells", availableSpells)
```

#### Using Helper Functions
```go
// Configuration errors
return errors.NewConfigError("invalid YAML syntax")

// Validation errors with field context
return errors.NewValidationError("command", "command cannot be empty")

// I/O errors with file path
return errors.NewIOError(filePath, "failed to read file")

// Execution errors with command context
return errors.NewExecutionError(command, "command execution failed")
```

#### Using Sentinel Errors
```go
// Check for specific error conditions
if errors.Is(err, errors.ErrSpellNotFound) {
    // Handle spell not found specifically
}

// Return sentinel errors
if spellName == "" {
    return errors.ErrSpellNotFound.WithContext("spell_name", spellName)
}
```

### 2. Wrapping Existing Errors

#### Basic Wrapping
```go
// Wrap standard errors
file, err := os.Open(path)
if err != nil {
    return errors.Wrap(errors.ErrorTypeIO, "failed to open configuration file", err).
        WithContext("path", path)
}

// Wrap with additional context
err := executeCommand(cmd)
if err != nil {
    return errors.Wrap(errors.ErrorTypeExecution, "script execution failed", err).
        WithContext("command", cmd.String()).
        WithContext("working_dir", cmd.Dir).
        WithContext("exit_code", cmd.ProcessState.ExitCode())
}
```

#### Context-Rich Wrapping
```go
// Multiple context fields
return errors.WrapWithContext(
    errors.ErrorTypeNetwork,
    "failed to download update",
    err,
    map[string]interface{}{
        "url": downloadURL,
        "version": targetVersion,
        "attempt": retryCount,
        "timeout": timeout.String(),
    },
)
```

### 3. Error Checking and Type Detection

#### Type-Based Error Handling
```go
// Check error type
if errors.IsType(err, errors.ErrorTypePermission) {
    // Handle permission errors specifically
    log.Error("Permission error occurred", "error", err)
    return showPermissionDialog()
}

// Extract SpellbookError for detailed handling
var spellErr *errors.SpellbookError
if errors.As(err, &spellErr) {
    // Access structured context
    if spellName, ok := spellErr.Context["spell_name"]; ok {
        log.Debug("Failed spell", "spell", spellName)
    }
    
    // Get structured logging fields
    logFields := spellErr.LogFields()
    log.Error("Operation failed", logFields)
}
```

#### Sentinel Error Checking
```go
// Check for specific conditions
if errors.Is(err, errors.ErrConfigNotFound) {
    // Create default configuration
    return createDefaultConfig()
}

if errors.Is(err, errors.ErrBrowserNotFound) {
    // Fallback to system default
    return openWithSystemDefault(url)
}
```

## 📝 Error Context Guidelines

### Essential Context Fields

#### Configuration Errors
```go
errors.NewConfigError("invalid spell definition").
    WithContext("spell_name", spellName).
    WithContext("file_path", configPath).
    WithContext("line_number", lineNum)
```

#### Execution Errors
```go
errors.NewExecutionError(command, "execution failed").
    WithContext("command", command).
    WithContext("args", args).
    WithContext("working_dir", workingDir).
    WithContext("timeout", timeout).
    WithContext("exit_code", exitCode)
```

#### File I/O Errors
```go
errors.NewIOError(path, "file operation failed").
    WithContext("operation", "read|write|delete").
    WithContext("file_size", fileSize).
    WithContext("permissions", fileMode)
```

#### Network Errors
```go
errors.NewNetworkError(url, "request failed").
    WithContext("method", "GET|POST|...").
    WithContext("status_code", statusCode).
    WithContext("timeout", timeout).
    WithContext("retry_count", retryCount)
```

### Context Best Practices

1. **Include actionable information** for debugging
2. **Avoid sensitive data** in context (passwords, tokens)
3. **Use consistent key names** across similar operations
4. **Include suggestions** for resolution when possible

```go
// Good context
return errors.NewConfigError("missing required field").
    WithContext("field", "command").
    WithContext("section", "grimoire.editor").
    WithContext("suggestion", "add 'command' field to action definition")

// Avoid sensitive data
return errors.NewNetworkError(url, "authentication failed").
    WithContext("url", sanitizeURL(url)).  // Remove credentials
    WithContext("auth_method", "token")    // Don't include the actual token
```

## 🚨 Migration from fmt.Errorf

### Before (fmt.Errorf)
```go
// Old pattern
if err != nil {
    return fmt.Errorf("failed to load %s: %w", path, err)
}

return fmt.Errorf("spell %s not found", spellName)
```

### After (SpellbookError)
```go
// New pattern with context
if err != nil {
    return errors.Wrap(errors.ErrorTypeIO, "failed to load configuration", err).
        WithContext("path", path)
}

return errors.NewNotFoundError("spell", spellName)
```

### Migration Checklist

- [ ] Replace `fmt.Errorf` with appropriate `errors.New` or `errors.Wrap`
- [ ] Add appropriate error type classification
- [ ] Include relevant context information
- [ ] Use helper functions for common patterns
- [ ] Convert to sentinel errors where appropriate
- [ ] Update error checking to use `errors.Is` and `errors.As`

## 📊 Logging Integration

### Structured Logging
```go
// Get structured fields for logging
var spellErr *errors.SpellbookError
if errors.As(err, &spellErr) {
    // Use LogFields() for structured logging
    logger.Error("Operation failed", spellErr.LogFields())
}

// Manual field extraction
logger.Error("Spell execution failed",
    "error_type", "execution",
    "spell_name", spellName,
    "error", err.Error(),
)
```

### User-Friendly Messages
```go
// Get user-friendly error message
userMessage := errors.GetUserMessage(err)
showNotification(userMessage)

// Example outputs:
// "Configuration error: spell not found"
// "Permission error: access denied. Please check the permissions guide."
// "Network error: connection timeout. Please check your connection."
```

## 🧪 Testing Error Handling

### Error Creation Tests
```go
func TestConfigError(t *testing.T) {
    err := errors.NewConfigError("test error").
        WithContext("field", "command")
    
    assert.True(t, errors.IsType(err, errors.ErrorTypeConfig))
    assert.Equal(t, "test error", err.Error())
    assert.Equal(t, "command", err.Context["field"])
}
```

### Error Wrapping Tests
```go
func TestErrorWrapping(t *testing.T) {
    cause := fmt.Errorf("original error")
    wrapped := errors.Wrap(errors.ErrorTypeIO, "wrapped error", cause)
    
    assert.True(t, errors.Is(wrapped, cause))
    assert.Contains(t, wrapped.Error(), "wrapped error")
    assert.Contains(t, wrapped.Error(), "original error")
}
```

### Sentinel Error Tests
```go
func TestSentinelErrors(t *testing.T) {
    err := errors.ErrSpellNotFound.WithContext("spell", "test")
    
    assert.True(t, errors.Is(err, errors.ErrSpellNotFound))
    assert.True(t, errors.IsType(err, errors.ErrorTypeConfig))
}
```

## 🔄 Error Recovery Patterns

### Graceful Degradation
```go
func loadConfig() error {
    err := loadPrimaryConfig()
    if errors.Is(err, errors.ErrConfigNotFound) {
        // Fallback to default configuration
        log.Warn("Primary config not found, using defaults")
        return loadDefaultConfig()
    }
    if errors.IsType(err, errors.ErrorTypeValidation) {
        // Try to fix common validation issues
        return fixAndReloadConfig(err)
    }
    return err
}
```

### Retry Logic
```go
func executeWithRetry(action func() error, maxRetries int) error {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        err := action()
        if err == nil {
            return nil
        }
        
        // Only retry on specific error types
        if errors.IsType(err, errors.ErrorTypeNetwork) || 
           errors.IsType(err, errors.ErrorTypeTimeout) {
            lastErr = err
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        
        // Don't retry other error types
        return err
    }
    
    return errors.Wrap(errors.ErrorTypeSystem, "max retries exceeded", lastErr).
        WithContext("max_retries", maxRetries)
}
```

## 📈 Performance Considerations

### Error Context Efficiency
```go
// Efficient context building
context := map[string]interface{}{
    "operation": "file_read",
    "path": path,
    "size": fileSize,
}
return errors.WrapWithContext(errors.ErrorTypeIO, "read failed", err, context)

// Avoid expensive operations in error paths
// Bad: Don't compute heavy context unless error actually occurs
if heavyComputation() != nil {
    return errors.New(errors.ErrorTypeValidation, "validation failed").
        WithContext("expensive_data", computeExpensiveData()) // Avoid this
}
```

### Memory Usage
```go
// Limit context data size for long-running operations
if len(data) > 1000 {
    return errors.NewIOError(path, "data too large").
        WithContext("data_size", len(data)).
        WithContext("data_preview", string(data[:100])+"...") // Truncate large data
}
```

## 🔗 Related Documentation

- [Debugging Guide](../troubleshooting/debugging.md)
- [TDD Development Guide](../guide/tdd-development.md)
- [Troubleshooting](../troubleshooting/)
- [Development Documentation](./)

---

**Remember: Consistent error handling improves debugging, user experience, and code maintainability. When in doubt, provide more context rather than less.**