package hotkey

import (
	"runtime"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		input    string
		wantLen  int
		wantErr  bool
		validate func(t *testing.T, seq KeySequence)
	}{
		{
			name:    "Single key",
			input:   "a",
			wantLen: 1,
			validate: func(t *testing.T, seq KeySequence) {
				if seq.Keys[0].Name != "a" {
					t.Errorf("Expected key name 'a', got '%s'", seq.Keys[0].Name)
				}
				if len(seq.Keys[0].Modifiers) != 0 {
					t.Errorf("Expected no modifiers, got %v", seq.Keys[0].Modifiers)
				}
			},
		},
		{
			name:    "Key with modifier",
			input:   "ctrl+a",
			wantLen: 1,
			validate: func(t *testing.T, seq KeySequence) {
				if seq.Keys[0].Name != "a" {
					t.Errorf("Expected key name 'a', got '%s'", seq.Keys[0].Name)
				}
				if len(seq.Keys[0].Modifiers) != 1 || seq.Keys[0].Modifiers[0] != "ctrl" {
					t.Errorf("Expected modifier [ctrl], got %v", seq.Keys[0].Modifiers)
				}
			},
		},
		{
			name:    "Multiple modifiers",
			input:   "ctrl+shift+a",
			wantLen: 1,
			validate: func(t *testing.T, seq KeySequence) {
				if seq.Keys[0].Name != "a" {
					t.Errorf("Expected key name 'a', got '%s'", seq.Keys[0].Name)
				}
				if len(seq.Keys[0].Modifiers) != 2 {
					t.Errorf("Expected 2 modifiers, got %d", len(seq.Keys[0].Modifiers))
				}
			},
		},
		{
			name:    "Sequential keys",
			input:   "g,s",
			wantLen: 2,
			validate: func(t *testing.T, seq KeySequence) {
				if seq.Keys[0].Name != "g" {
					t.Errorf("Expected first key 'g', got '%s'", seq.Keys[0].Name)
				}
				if seq.Keys[1].Name != "s" {
					t.Errorf("Expected second key 's', got '%s'", seq.Keys[1].Name)
				}
			},
		},
		{
			name:    "Sequential with modifiers",
			input:   "ctrl+g,s",
			wantLen: 2,
			validate: func(t *testing.T, seq KeySequence) {
				if len(seq.Keys[0].Modifiers) != 1 {
					t.Errorf("Expected first key to have modifier")
				}
				if len(seq.Keys[1].Modifiers) != 0 {
					t.Errorf("Expected second key to have no modifiers")
				}
			},
		},
		{
			name:    "Function key",
			input:   "f1",
			wantLen: 1,
			validate: func(t *testing.T, seq KeySequence) {
				if seq.Keys[0].Name != "f1" {
					t.Errorf("Expected key name 'f1', got '%s'", seq.Keys[0].Name)
				}
			},
		},
		{
			name:    "Special key",
			input:   "space",
			wantLen: 1,
			validate: func(t *testing.T, seq KeySequence) {
				if seq.Keys[0].Name != "space" {
					t.Errorf("Expected key name 'space', got '%s'", seq.Keys[0].Name)
				}
			},
		},
		{
			name:    "Platform modifier",
			input:   getPlatformModifier() + "+space",
			wantLen: 1,
			validate: func(t *testing.T, seq KeySequence) {
				if len(seq.Keys[0].Modifiers) != 1 {
					t.Errorf("Expected platform modifier")
				}
			},
		},
		{
			name:    "Empty sequence",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Empty key in sequence",
			input:   "a,,b",
			wantErr: true,
		},
		{
			name:    "Unknown modifier",
			input:   "invalid+a",
			wantErr: true,
		},
		{
			name:    "Unknown key",
			input:   "unknownkey",
			wantErr: true,
		},
		{
			name:    "Case insensitive",
			input:   "CTRL+A",
			wantLen: 1,
			validate: func(t *testing.T, seq KeySequence) {
				if seq.Keys[0].Name != "a" {
					t.Errorf("Expected lowercase key 'a', got '%s'", seq.Keys[0].Name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, err := parser.Parse(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if len(seq.Keys) != tt.wantLen {
				t.Errorf("Parse() got %d keys, want %d", len(seq.Keys), tt.wantLen)
			}

			if tt.validate != nil {
				tt.validate(t, seq)
			}
		})
	}
}

func TestKeySequence_String(t *testing.T) {
	tests := []struct {
		name     string
		sequence KeySequence
		want     string
	}{
		{
			name: "Single key",
			sequence: KeySequence{
				Keys: []Key{{Name: "a", Modifiers: []string{}}},
			},
			want: "a",
		},
		{
			name: "Key with modifier",
			sequence: KeySequence{
				Keys: []Key{{Name: "a", Modifiers: []string{"ctrl"}}},
			},
			want: "ctrl+a",
		},
		{
			name: "Multiple modifiers",
			sequence: KeySequence{
				Keys: []Key{{Name: "a", Modifiers: []string{"ctrl", "shift"}}},
			},
			want: "ctrl+shift+a",
		},
		{
			name: "Sequential keys",
			sequence: KeySequence{
				Keys: []Key{
					{Name: "g", Modifiers: []string{}},
					{Name: "s", Modifiers: []string{}},
				},
			},
			want: "g,s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sequence.String(); got != tt.want {
				t.Errorf("KeySequence.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getPlatformModifier() string {
	switch runtime.GOOS {
	case "darwin":
		return "cmd"
	case "windows":
		return "win"
	default:
		return "ctrl"
	}
}
