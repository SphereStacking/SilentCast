package shell

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}
}

func TestManagerGetShell(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()

	tests := []struct {
		name      string
		shellName string
		wantErr   bool
	}{
		{"Empty name returns default", "", false},
		{"Find sh", "sh", false},
		{"Find nonexistent", "nonexistentshell", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, err := manager.GetShell(ctx, tt.shellName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShell(%q) error = %v, wantErr %v", tt.shellName, err, tt.wantErr)
				return
			}

			if !tt.wantErr && shell != nil {
				t.Logf("Got shell: %s", FormatShellInfo(shell))
			}
		})
	}
}

func TestManagerListShells(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()

	shells, err := manager.ListShells(ctx)
	if err != nil {
		t.Fatalf("ListShells() error: %v", err)
	}

	if len(shells) == 0 {
		t.Error("ListShells() returned no shells")
	}

	t.Logf("Available shells:")
	for _, shell := range shells {
		t.Logf("  %s", FormatShellInfo(&shell))
	}
}

func TestManagerCreateCommand(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()

	tests := []struct {
		name   string
		shell  *Shell
		script string
		env    []string
	}{
		{
			name:   "Default shell",
			shell:  nil,
			script: "echo hello",
			env:    nil,
		},
		{
			name: "Specific shell",
			shell: &Shell{
				Name:       "sh",
				Executable: "/bin/sh",
				Type:       ShellTypeBourne,
				Args:       []string{"-c"},
			},
			script: "echo world",
			env:    []string{"TEST=1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := manager.CreateCommand(ctx, tt.shell, tt.script, tt.env)
			if err != nil {
				t.Fatalf("CreateCommand() error: %v", err)
			}

			if cmd == nil {
				t.Fatal("CreateCommand() returned nil")
			}

			// Check command
			if cmd.Path == "" {
				t.Error("Command path is empty")
			}

			// Check args
			if len(cmd.Args) < 2 {
				t.Error("Command should have at least 2 args (shell and script)")
			}

			// Check environment
			if len(tt.env) > 0 {
				found := false
				for _, e := range cmd.Env {
					if e == tt.env[0] {
						found = true
						break
					}
				}
				if !found {
					t.Error("Custom environment variable not found")
				}
			}
		})
	}
}

func TestManagerGetShellForScript(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()

	tests := []struct {
		name           string
		script         string
		preferredShell string
		checkName      string // Expected shell name (flexible check)
	}{
		{
			name:           "No shebang, no preference",
			script:         "echo hello",
			preferredShell: "",
			checkName:      "", // Any shell is ok
		},
		{
			name:           "With preference",
			script:         "echo hello",
			preferredShell: "sh",
			checkName:      "sh",
		},
		{
			name:           "Invalid preference falls back",
			script:         "echo hello",
			preferredShell: "invalidshell",
			checkName:      "", // Should use default
		},
		{
			name:           "Python shebang",
			script:         "#!/usr/bin/env python\nprint('hello')",
			preferredShell: "",
			checkName:      "python",
		},
		{
			name:           "Bash shebang",
			script:         "#!/bin/bash\necho hello",
			preferredShell: "",
			checkName:      "bash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, err := manager.GetShellForScript(ctx, tt.script, tt.preferredShell)
			if err != nil {
				t.Fatalf("GetShellForScript() error: %v", err)
			}

			if shell == nil {
				t.Fatal("GetShellForScript() returned nil")
			}

			t.Logf("Got shell: %s", FormatShellInfo(shell))

			if tt.checkName != "" {
				if !strings.Contains(shell.Name, tt.checkName) {
					t.Errorf("Expected shell containing %q, got %q", tt.checkName, shell.Name)
				}
			}
		})
	}
}

func TestDetectShebang(t *testing.T) {
	manager := &manager{
		detector: NewDetector(),
	}
	ctx := context.Background()

	tests := []struct {
		name      string
		script    string
		wantShell bool
		shellName string
	}{
		{
			name:      "No shebang",
			script:    "echo hello",
			wantShell: false,
		},
		{
			name:      "Empty shebang",
			script:    "#!\n",
			wantShell: false,
		},
		{
			name:      "Bash shebang",
			script:    "#!/bin/bash\necho hello",
			wantShell: true,
			shellName: "bash",
		},
		{
			name:      "Env shebang",
			script:    "#!/usr/bin/env python3\nprint('hello')",
			wantShell: true,
			shellName: "python",
		},
		{
			name:      "Shebang with args",
			script:    "#!/bin/bash -e\necho hello",
			wantShell: true,
			shellName: "bash",
		},
		{
			name:      "Invalid shebang path",
			script:    "#!/nonexistent/shell\necho hello",
			wantShell: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell := manager.detectShebang(ctx, tt.script)

			if tt.wantShell {
				if shell == nil {
					t.Error("detectShebang() returned nil, expected shell")
				} else if tt.shellName != "" && !strings.Contains(shell.Name, tt.shellName) {
					t.Errorf("detectShebang() returned shell %q, want containing %q", shell.Name, tt.shellName)
				}
			} else {
				if shell != nil {
					t.Errorf("detectShebang() returned shell %s, expected nil", shell.Name)
				}
			}
		})
	}
}

func TestFormatShellInfo(t *testing.T) {
	tests := []struct {
		name     string
		shell    *Shell
		contains []string
	}{
		{
			name:     "Nil shell",
			shell:    nil,
			contains: []string{"No shell"},
		},
		{
			name: "Basic shell",
			shell: &Shell{
				Name: "bash",
				Type: ShellTypeBourne,
			},
			contains: []string{"bash", "Bourne-compatible"},
		},
		{
			name: "Shell with version",
			shell: &Shell{
				Name:    "python",
				Type:    ShellTypeInterpreter,
				Version: "Python 3.9.0",
			},
			contains: []string{"python", "Interpreter", "Python 3.9.0"},
		},
		{
			name: "Default shell",
			shell: &Shell{
				Name:      "zsh",
				Type:      ShellTypeBourne,
				IsDefault: true,
			},
			contains: []string{"zsh", "[default]"},
		},
		{
			name: "Long version truncated",
			shell: &Shell{
				Name:    "bash",
				Version: "GNU bash, version 5.1.16(1)-release (x86_64-pc-linux-gnu) with lots of extra information that should be truncated",
			},
			contains: []string{"bash", "...", "GNU bash"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatShellInfo(tt.shell)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("FormatShellInfo() = %q, want containing %q", result, expected)
				}
			}
		})
	}
}

func TestGetShellEnvironment(t *testing.T) {
	tests := []struct {
		name          string
		shell         *Shell
		additionalEnv map[string]string
		checkEnv      []string
	}{
		{
			name:          "Nil shell",
			shell:         nil,
			additionalEnv: map[string]string{"TEST": "1"},
			checkEnv:      []string{"TEST=1"},
		},
		{
			name: "Python shell",
			shell: &Shell{
				Name: "python",
				Type: ShellTypeInterpreter,
			},
			additionalEnv: nil,
			checkEnv:      []string{"PYTHONUNBUFFERED=1"},
		},
		{
			name: "Node shell",
			shell: &Shell{
				Name: "node",
				Type: ShellTypeInterpreter,
			},
			additionalEnv: nil,
			checkEnv:      []string{"NODE_NO_WARNINGS=1"},
		},
		{
			name: "PowerShell",
			shell: &Shell{
				Name: "powershell",
				Type: ShellTypePowerShell,
			},
			additionalEnv: nil,
			checkEnv:      []string{"PSModulePath="},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := GetShellEnvironment(tt.shell, tt.additionalEnv)

			// Should include system environment
			if len(env) <= len(os.Environ()) {
				t.Error("Shell environment should include system environment")
			}

			// Check for expected environment variables
			for _, expected := range tt.checkEnv {
				found := false
				for _, e := range env {
					if e == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected environment variable %q not found", expected)
				}
			}
		})
	}
}

func TestManagerCaching(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()

	// First call
	shell1, err := manager.GetShell(ctx, "sh")
	if err != nil {
		t.Fatalf("First GetShell() error: %v", err)
	}

	// Second call should use cache
	shell2, err := manager.GetShell(ctx, "sh")
	if err != nil {
		t.Fatalf("Second GetShell() error: %v", err)
	}

	// Should be same instance (from cache)
	if shell1 != shell2 {
		t.Error("GetShell() not using cache")
	}
}

func BenchmarkDetectShells(b *testing.B) {
	detector := NewDetector()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.DetectShells(ctx)
	}
}

func BenchmarkFindShell(b *testing.B) {
	detector := NewDetector()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detector.FindShell(ctx, "sh")
	}
}

func TestCreateInterpreterCommand(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()

	// Test with Python interpreter
	python := &Shell{
		Name:       "python",
		Executable: "/usr/bin/python3",
		Type:       ShellTypeInterpreter,
		Args:       []string{"-c"},
	}

	code := "print('Hello from interpreter mode')"
	cmd, err := manager.CreateInterpreterCommand(ctx, python, code, nil, nil)
	if err != nil {
		t.Fatalf("CreateInterpreterCommand() error: %v", err)
	}

	if cmd == nil {
		t.Fatal("CreateInterpreterCommand() returned nil")
	}

	// Check command structure
	if cmd.Path != python.Executable {
		t.Errorf("Command path = %q, want %q", cmd.Path, python.Executable)
	}

	// Check args
	expectedArgs := []string{python.Executable, "-c", code}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("Command args length = %d, want %d", len(cmd.Args), len(expectedArgs))
	}

	// Test with non-interpreter shell
	bash := &Shell{
		Name: "bash",
		Type: ShellTypeBourne,
	}

	_, err = manager.CreateInterpreterCommand(ctx, bash, code, nil, nil)
	if err == nil {
		t.Error("CreateInterpreterCommand() with non-interpreter should return error")
	}
}

func TestGetInterpreterForScript(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()

	tests := []struct {
		name        string
		script      string
		preferred   string
		expectError bool
		expectType  string
	}{
		{
			name:        "Python script with shebang",
			script:      "#!/usr/bin/env python3\nimport sys\nprint(sys.version)",
			preferred:   "",
			expectError: false,
			expectType:  "python",
		},
		{
			name:        "Python script with content detection",
			script:      "import os\nprint('hello')\ndef main():\n    pass",
			preferred:   "",
			expectError: false,
			expectType:  "python",
		},
		{
			name:        "Node.js script",
			script:      "const fs = require('fs');\nconsole.log('hello');",
			preferred:   "",
			expectError: false,
			expectType:  "node",
		},
		{
			name:        "Preferred interpreter",
			script:      "print('hello')",
			preferred:   "python",
			expectError: false,
			expectType:  "python",
		},
		{
			name:        "Shell script (not an interpreter)",
			script:      "echo hello\nif [ -f file ]; then\n  echo found\nfi",
			preferred:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter, err := manager.GetInterpreterForScript(ctx, tt.script, tt.preferred)

			if tt.expectError {
				if err == nil {
					t.Error("GetInterpreterForScript() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GetInterpreterForScript() unexpected error: %v", err)
				return
			}

			if interpreter == nil {
				t.Error("GetInterpreterForScript() returned nil interpreter")
				return
			}

			if !IsInterpreter(interpreter) {
				t.Error("GetInterpreterForScript() returned non-interpreter")
			}

			if tt.expectType != "" && !strings.Contains(interpreter.Name, tt.expectType) {
				t.Errorf("GetInterpreterForScript() returned %q, expected containing %q", interpreter.Name, tt.expectType)
			}
		})
	}
}

func TestDetectInterpreterFromContent(t *testing.T) {
	manager := &manager{
		detector: NewDetector(),
	}
	ctx := context.Background()

	tests := []struct {
		name        string
		script      string
		expectType  string
		expectFound bool
	}{
		{
			name:        "Python with imports",
			script:      "import sys\nimport os\nprint('hello')",
			expectType:  "python",
			expectFound: true,
		},
		{
			name:        "Node.js with requires",
			script:      "const fs = require('fs');\nconsole.log('hello');",
			expectType:  "node",
			expectFound: true,
		},
		{
			name:        "Ruby script",
			script:      "puts 'hello'\n[1,2,3].each { |x| puts x }\nend",
			expectType:  "ruby",
			expectFound: true,
		},
		{
			name:        "Shell script",
			script:      "echo hello\nif [ -f file ]; then\n  echo found\nfi",
			expectFound: false,
		},
		{
			name:        "Empty script",
			script:      "",
			expectFound: false,
		},
		{
			name:        "Single Python statement",
			script:      "print('hello')",
			expectFound: false, // Only one match, needs 2+
		},
		{
			name:        "Python with function",
			script:      "def hello():\n    print('hello')\nhello()",
			expectType:  "python",
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip Ruby test if Ruby is not available in the environment
			if tt.expectType == "ruby" {
				if _, err := manager.detector.FindShell(ctx, "ruby"); err != nil {
					t.Skip("Ruby not available in test environment")
				}
			}

			interpreter := manager.detectInterpreterFromContent(ctx, tt.script)

			if tt.expectFound {
				if interpreter == nil {
					t.Error("detectInterpreterFromContent() expected interpreter but got nil")
				} else if tt.expectType != "" && !strings.Contains(interpreter.Name, tt.expectType) {
					t.Errorf("detectInterpreterFromContent() returned %q, expected containing %q", interpreter.Name, tt.expectType)
				}
			} else {
				if interpreter != nil {
					t.Errorf("detectInterpreterFromContent() expected nil but got %s", interpreter.Name)
				}
			}
		})
	}
}

func ExampleManager_GetShellForScript() {
	manager := NewManager()
	ctx := context.Background()

	// Script with Python shebang
	script := `#!/usr/bin/env python3
print("Hello from Python")
`

	shell, err := manager.GetShellForScript(ctx, script, "")
	if err != nil {
		panic(err)
	}

	// Will detect Python from shebang
	_ = shell
}

func ExampleManager_CreateInterpreterCommand() {
	manager := NewManager()
	ctx := context.Background()

	// Get Python interpreter
	python, err := manager.GetShell(ctx, "python")
	if err != nil {
		panic(err)
	}

	// Create command for direct interpreter execution
	code := "print('Hello from interpreter mode')"
	cmd, err := manager.CreateInterpreterCommand(ctx, python, code, nil, nil)
	if err != nil {
		panic(err)
	}

	// Command will execute: python -c "print('Hello from interpreter mode')"
	_ = cmd
}
