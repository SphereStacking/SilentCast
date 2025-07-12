package hotkey

import (
	"testing"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name       string
		sequence   string
		spellName  string
		registered map[string]string // pre-registered sequences
		wantErr    bool
		errMsg     string
	}{
		{
			name:      "Valid single key",
			sequence:  "a",
			spellName: "test",
			wantErr:   false,
		},
		{
			name:      "Valid sequence",
			sequence:  "g,s",
			spellName: "git_status",
			wantErr:   false,
		},
		{
			name:      "Empty sequence",
			sequence:  "",
			spellName: "test",
			wantErr:   true,
			errMsg:    "empty sequence",
		},
		{
			name:      "Invalid key",
			sequence:  "invalid_key",
			spellName: "test",
			wantErr:   true,
			errMsg:    "unknown key",
		},
		{
			name:      "Duplicate exact sequence",
			sequence:  "ctrl+a",
			spellName: "new_spell",
			registered: map[string]string{
				"ctrl+a": "existing_spell",
			},
			wantErr: true,
			errMsg:  "already registered",
		},
		{
			name:      "Same spell re-registration",
			sequence:  "ctrl+a",
			spellName: "same_spell",
			registered: map[string]string{
				"ctrl+a": "same_spell",
			},
			wantErr: false,
		},
		{
			name:      "Prefix conflict - shorter",
			sequence:  "g",
			spellName: "new_spell",
			registered: map[string]string{
				"g,s": "git_status",
			},
			wantErr: true,
			errMsg:  "conflicts with longer sequence",
		},
		{
			name:      "Prefix conflict - longer",
			sequence:  "g,s,p",
			spellName: "new_spell",
			registered: map[string]string{
				"g,s": "git_status",
			},
			wantErr: true,
			errMsg:  "conflicts with shorter sequence",
		},
		{
			name:      "No conflict - different prefix",
			sequence:  "h,s",
			spellName: "help_search",
			registered: map[string]string{
				"g,s": "git_status",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator()

			// Pre-register sequences
			for seq, spell := range tt.registered {
				if err := validator.Register(seq, spell); err != nil {
					t.Fatalf("Failed to pre-register sequence: %v", err)
				}
			}

			err := validator.Validate(tt.sequence, tt.spellName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errMsg, err.Error())
				}
			}
		})
	}
}

func TestValidator_Register(t *testing.T) {
	validator := NewValidator()

	// Test successful registration
	err := validator.Register("ctrl+a", "select_all")
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// Test duplicate registration
	err = validator.Register("ctrl+a", "different_spell")
	if err == nil {
		t.Error("Expected error for duplicate registration")
	}

	// Test registration after unregister
	validator.Unregister("ctrl+a")
	err = validator.Register("ctrl+a", "new_spell")
	if err != nil {
		t.Errorf("Register() after unregister error = %v", err)
	}
}

func TestValidator_GetRegistered(t *testing.T) {
	validator := NewValidator()

	// Register some sequences
	sequences := map[string]string{
		"ctrl+a": "select_all",
		"g,s":    "git_status",
		"f1":     "help",
	}

	for seq, spell := range sequences {
		if err := validator.Register(seq, spell); err != nil {
			t.Fatalf("Failed to register sequence: %v", err)
		}
	}

	// Get registered sequences
	registered := validator.GetRegistered()

	if len(registered) != len(sequences) {
		t.Errorf("Expected %d registered sequences, got %d", len(sequences), len(registered))
	}

	// Verify each sequence is registered correctly
	for seq, expectedSpell := range sequences {
		// The validator normalizes sequences, so we need to check the normalized form
		found := false
		for _, regSpell := range registered {
			if regSpell == expectedSpell {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Sequence '%s' with spell '%s' not found in registered", seq, expectedSpell)
		}
	}
}

func TestValidator_Clear(t *testing.T) {
	validator := NewValidator()

	// Register some sequences
	validator.Register("ctrl+a", "select_all")
	validator.Register("g,s", "git_status")

	// Clear all registrations
	validator.Clear()

	// Verify all cleared
	registered := validator.GetRegistered()
	if len(registered) != 0 {
		t.Errorf("Expected 0 registered sequences after clear, got %d", len(registered))
	}

	// Should be able to register previously registered sequences
	err := validator.Register("ctrl+a", "new_spell")
	if err != nil {
		t.Errorf("Register() after clear error = %v", err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}
