package benchmarks

import (
	"testing"

	"github.com/SphereStacking/silentcast/internal/hotkey"
)

// BenchmarkKeyParsing measures hotkey sequence parsing performance
func BenchmarkKeyParsing(b *testing.B) {
	parser := hotkey.NewParser()
	
	tests := []struct {
		name     string
		sequence string
	}{
		{"Single key", "e"},
		{"Function key", "f1"},
		{"Modifier key", "ctrl+c"},
		{"Complex modifier", "ctrl+alt+shift+f12"},
		{"Simple sequence", "g,s"},
		{"Complex sequence", "d,o,c,k,e,r"},
		{"Mixed sequence", "ctrl+g,s"},
		{"Long sequence", "a,b,c,d,e,f,g"},
	}
	
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := parser.Parse(tt.sequence)
				if err != nil {
					b.Fatalf("Failed to parse %s: %v", tt.sequence, err)
				}
			}
		})
	}
}

// BenchmarkKeyValidation measures key validation performance
func BenchmarkKeyValidation(b *testing.B) {
	validator := hotkey.NewValidator()
	
	tests := []struct {
		name     string
		sequence string
	}{
		{"Valid single", "a"},
		{"Valid function", "f10"},
		{"Valid modifier", "alt+space"},
		{"Valid sequence", "g,s,t"},
		{"Invalid key", "invalid_key"},
		{"Invalid modifier", "invalid+a"},
		{"Invalid sequence", "a,,b"},
	}
	
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = validator.Validate(tt.sequence, "test")
			}
		})
	}
}

// BenchmarkHotkeyRegistration measures hotkey registration performance
func BenchmarkHotkeyRegistration(b *testing.B) {
	manager := hotkey.NewMockManager()
	
	sequences := []string{
		"e", "t", "b", "f", "c",
		"g,s", "g,p", "g,c", "g,l",
		"d,l", "d,s", "d,r", "d,c",
		"ctrl+c", "ctrl+v", "alt+tab",
		"f1", "f2", "f3", "f4", "f5",
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for j, seq := range sequences {
			spellName := "spell_" + string(rune('a'+j%26))
			err := manager.Register(seq, spellName)
			if err != nil {
				b.Fatalf("Failed to register %s: %v", seq, err)
			}
		}
		
		// Clean up for next iteration
		for _, seq := range sequences {
			_ = manager.Unregister(seq)
		}
	}
}

// BenchmarkHotkeyUnregistration measures hotkey unregistration performance
func BenchmarkHotkeyUnregistration(b *testing.B) {
	sequences := []string{
		"e", "t", "b", "f", "c",
		"g,s", "g,p", "g,c", "g,l",
		"d,l", "d,s", "d,r", "d,c",
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		manager := hotkey.NewMockManager()
		
		// Register all sequences first
		for j, seq := range sequences {
			spellName := "spell_" + string(rune('a'+j%26))
			_ = manager.Register(seq, spellName)
		}
		
		b.StartTimer()
		// Unregister all sequences
		for _, seq := range sequences {
			err := manager.Unregister(seq)
			if err != nil {
				b.Fatalf("Failed to unregister %s: %v", seq, err)
			}
		}
		b.StopTimer()
	}
}

// BenchmarkHotkeyLookup measures spell lookup performance
func BenchmarkHotkeyLookup(b *testing.B) {
	manager := hotkey.NewMockManager()
	
	// Register many sequences
	sequences := make([]string, 100)
	for i := 0; i < 100; i++ {
		seq := "s" + string(rune('a'+i%26)) + string(rune('0'+i%10))
		spellName := "spell_" + string(rune('a'+i%26)) + string(rune('0'+i%10))
		sequences[i] = seq
		_ = manager.Register(seq, spellName)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Look up random sequence
		seq := sequences[i%len(sequences)]
		// Simulate spell lookup through mock manager
		_ = manager.SimulateKeyPress(seq)
	}
}

// BenchmarkSequenceCompletion measures sequence completion detection
func BenchmarkSequenceCompletion(b *testing.B) {
	parser := hotkey.NewParser()
	
	sequences := []string{
		"a,b,c",
		"g,s,t,a",
		"d,o,c,k,e,r",
		"s,y,s,t,e,m",
		"1,2,3,4,5",
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, seq := range sequences {
			keySeq, err := parser.Parse(seq)
			if err != nil {
				b.Fatalf("Failed to parse sequence %s: %v", seq, err)
			}
			
			// Simulate sequence completion detection
			_ = keySeq.String() // Force evaluation
		}
	}
}

// BenchmarkKeyEventProcessing measures key event processing performance
func BenchmarkKeyEventProcessing(b *testing.B) {
	manager := hotkey.NewMockManager()
	parser := hotkey.NewParser()
	
	// Setup handler
	handler := &TestHandler{
		HandleFunc: func(event hotkey.Event) error {
			return nil
		},
	}
	manager.SetHandler(handler)
	
	// Register test sequence
	_ = manager.Register("g,s", "git_status")
	
	keySeq, err := parser.Parse("g,s")
	if err != nil {
		b.Fatalf("Failed to parse sequence: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		event := hotkey.Event{
			Sequence:  keySeq,
			SpellName: "git_status",
		}
		
		err := handler.Handle(event)
		if err != nil {
			b.Fatalf("Handler failed: %v", err)
		}
	}
}

// BenchmarkConcurrentHotkeyRegistration measures concurrent registration performance
func BenchmarkConcurrentHotkeyRegistration(b *testing.B) {
	b.ReportAllocs()
	
	b.RunParallel(func(pb *testing.PB) {
		manager := hotkey.NewMockManager()
		counter := 0
		
		for pb.Next() {
			seq := "concurrent_" + string(rune('a'+counter%26))
			spellName := "spell_" + string(rune('a'+counter%26))
			
			err := manager.Register(seq, spellName)
			if err != nil {
				b.Fatalf("Failed to register %s: %v", seq, err)
			}
			
			counter++
		}
	})
}

// BenchmarkHotkeyMemoryUsage measures memory usage of hotkey operations
func BenchmarkHotkeyMemoryUsage(b *testing.B) {
	RunMemoryBenchmark(b, func() {
		manager := hotkey.NewMockManager()
		
		// Register many hotkeys
		for i := 0; i < 50; i++ {
			seq := "mem_" + string(rune('a'+i%26))
			spellName := "spell_" + string(rune('a'+i%26))
			_ = manager.Register(seq, spellName)
		}
		
		// Simulate key presses
		for i := 0; i < 10; i++ {
			seq := "mem_" + string(rune('a'+i%26))
			_ = manager.SimulateKeyPress(seq)
		}
	})
}

// BenchmarkKeySequenceBuilding measures sequence building performance
func BenchmarkKeySequenceBuilding(b *testing.B) {
	parser := hotkey.NewParser()
	
	// Test different sequence patterns
	patterns := []string{
		"a", "a,b", "a,b,c", "a,b,c,d",
		"ctrl+a", "ctrl+a,b", "ctrl+a,b,c",
		"f1,f2,f3", "shift+f1,f2,f3",
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, pattern := range patterns {
			_, err := parser.Parse(pattern)
			if err != nil {
				b.Fatalf("Failed to parse pattern %s: %v", pattern, err)
			}
		}
	}
}

// TestHandler is a test implementation of hotkey.Handler
type TestHandler struct {
	HandleFunc func(event hotkey.Event) error
}

func (h *TestHandler) Handle(event hotkey.Event) error {
	if h.HandleFunc != nil {
		return h.HandleFunc(event)
	}
	return nil
}