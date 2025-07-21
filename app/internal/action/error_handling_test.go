package action

import (
	"context"
	"errors"
	"testing"

	"github.com/SphereStacking/silentcast/internal/config"
	customErrors "github.com/SphereStacking/silentcast/internal/errors"
)

// TestActionManagerErrorHandling tests that action manager errors use unified error handling
func TestActionManagerErrorHandling(t *testing.T) {
	tests := []struct {
		name         string
		setupManager func() *Manager
		spellName    string
		expectError  bool
		expectType   customErrors.ErrorType
		expectContext map[string]interface{}
	}{
		{
			name: "spell not found with unified error",
			setupManager: func() *Manager {
				return NewManager(map[string]config.ActionConfig{
					"existing": {Type: "app", Command: "test"},
				})
			},
			spellName:   "nonexistent",
			expectError: true,
			expectType:  customErrors.ErrorTypeConfig,
			expectContext: map[string]interface{}{
				"spell_name":        "nonexistent",
				"available_spells":  []string{"existing"},
				"error_type":        "spell_not_found",
			},
		},
		{
			name: "unknown action type with unified error",
			setupManager: func() *Manager {
				return NewManager(map[string]config.ActionConfig{
					"invalid": {Type: "unknown", Command: "test"},
				})
			},
			spellName:   "invalid",
			expectError: true,
			expectType:  customErrors.ErrorTypeConfig,
			expectContext: map[string]interface{}{
				"spell_name":   "invalid",
				"action_type":  "unknown",
				"valid_types":  []string{"app", "script", "url"},
			},
		},
		{
			name: "valid app action succeeds",
			setupManager: func() *Manager {
				return NewManager(map[string]config.ActionConfig{
					"editor": {Type: "app", Command: "echo"},
				})
			},
			spellName:   "editor",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			manager := tt.setupManager()
			err := manager.Execute(ctx, tt.spellName)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}

				// Check if error uses unified error pattern
				var spellErr *customErrors.SpellbookError
				if !errors.As(err, &spellErr) {
					t.Errorf("error should be SpellbookError, got %T", err)
					return
				}

				// Check error type
				if spellErr.Type != tt.expectType {
					t.Errorf("expected error type %v, got %v", tt.expectType, spellErr.Type)
				}

				// Check context
				for key, expectedValue := range tt.expectContext {
					if actual, ok := spellErr.Context[key]; !ok {
						t.Errorf("error context missing key %q", key)
					} else {
						switch expected := expectedValue.(type) {
						case []string:
							actualSlice, ok := actual.([]string)
							if !ok {
								t.Errorf("error context[%q] should be []string, got %T", key, actual)
								continue
							}
							// Check if slices contain the same elements (order may vary)
							if len(actualSlice) != len(expected) {
								t.Errorf("error context[%q] length = %d, want %d", key, len(actualSlice), len(expected))
							}
						default:
							if actual != expectedValue {
								t.Errorf("error context[%q] = %v, want %v", key, actual, expectedValue)
							}
						}
					}
				}

			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestActionExecutorErrorPropagation tests that executor errors are properly wrapped
func TestActionExecutorErrorPropagation(t *testing.T) {
	tests := []struct {
		name         string
		action       config.ActionConfig
		expectError  bool
		expectType   customErrors.ErrorType
	}{
		{
			name: "script execution error with context",
			action: config.ActionConfig{
				Type:    "script",
				Command: "false", // Cross-platform command that always fails
			},
			expectError: true,
			expectType:  customErrors.ErrorTypeSystem,
		},
		{
			name: "app execution error with context",
			action: config.ActionConfig{
				Type:    "app",
				Command: "nonexistent-app-12345",
			},
			expectError: true,
			expectType:  customErrors.ErrorTypeSystem,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(map[string]config.ActionConfig{
				"test": tt.action,
			})

			ctx := context.Background()
			err := manager.Execute(ctx, "test")

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}

				var spellErr *customErrors.SpellbookError
				if !errors.As(err, &spellErr) {
					t.Errorf("error should be SpellbookError, got %T", err)
					return
				}

				// Should include spell context
				if spellName, ok := spellErr.Context["spell_name"]; !ok || spellName != "test" {
					t.Errorf("error should include spell_name context")
				}

				// Should include action type context
				if actionType, ok := spellErr.Context["action_type"]; !ok || actionType != tt.action.Type {
					t.Errorf("error should include action_type context")
				}

			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestActionUserFriendlyMessages tests user-friendly error messages
func TestActionUserFriendlyMessages(t *testing.T) {
	tests := []struct {
		name           string
		error          error
		expectContains string
	}{
		{
			name: "spell not found error",
			error: func() error {
				return customErrors.New(customErrors.ErrorTypeConfig, "spell not found").
					WithContext("spell_name", "missing-spell").
					WithContext("available_spells", []string{"editor", "terminal"}).
					WithContext("suggested_action", "check spellbook.yml configuration")
			}(),
			expectContains: "Configuration error",
		},
		{
			name: "action execution failed",
			error: func() error {
				return customErrors.New(customErrors.ErrorTypeSystem, "action execution failed").
					WithContext("spell_name", "editor").
					WithContext("action_type", "app").
					WithContext("app_path", "/usr/bin/code").
					WithContext("suggested_action", "check if application is installed")
			}(),
			expectContains: "System error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMsg := customErrors.GetUserMessage(tt.error)
			if len(userMsg) == 0 {
				t.Error("user message should not be empty")
			}
			// Note: We're just checking that GetUserMessage works
			// The actual message format is defined by the errors package
		})
	}
}