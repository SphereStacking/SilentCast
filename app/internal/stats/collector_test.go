package stats

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	tempDir := t.TempDir()

	cfg := Config{
		Enabled:      true,
		DataFile:     filepath.Join(tempDir, "stats.json"),
		SaveInterval: 1 * time.Second,
	}

	collector := NewCollector(cfg)

	if collector.enabled != cfg.Enabled {
		t.Error("Expected collector to be enabled")
	}

	if collector.stats.InstallDate.IsZero() {
		t.Error("Expected install date to be set")
	}
}

func TestRecordLaunch(t *testing.T) {
	collector := &Collector{
		enabled: true,
		stats: &Statistics{
			DailyUsage: make(map[string]*DailyUsage),
		},
	}

	initialLaunches := collector.stats.TotalLaunches
	collector.RecordLaunch()

	if collector.stats.TotalLaunches != initialLaunches+1 {
		t.Errorf("Expected launches to be %d, got %d",
			initialLaunches+1, collector.stats.TotalLaunches)
	}

	if collector.stats.LastUsed.IsZero() {
		t.Error("Expected LastUsed to be updated")
	}

	// Check daily usage
	today := time.Now().Format("2006-01-02")
	if daily, ok := collector.stats.DailyUsage[today]; !ok || daily.Launches != 1 {
		t.Error("Expected daily launch count to be 1")
	}
}

func TestRecordSpellCast(t *testing.T) {
	collector := &Collector{
		enabled: true,
		stats: &Statistics{
			SpellStats: make(map[string]*SpellStat),
			DailyUsage: make(map[string]*DailyUsage),
		},
	}

	spell := "test_spell"
	execTime := 100 * time.Millisecond

	// Record successful cast
	collector.RecordSpellCast(spell, true, execTime)

	stat, ok := collector.stats.SpellStats[spell]
	if !ok {
		t.Fatal("Expected spell stat to be created")
	}

	if stat.CastCount != 1 {
		t.Errorf("Expected cast count 1, got %d", stat.CastCount)
	}

	if stat.SuccessRate != 1.0 {
		t.Errorf("Expected success rate 1.0, got %f", stat.SuccessRate)
	}

	if stat.AvgExecTime != 100.0 {
		t.Errorf("Expected avg exec time 100ms, got %f", stat.AvgExecTime)
	}

	// Record failed cast
	collector.RecordSpellCast(spell, false, 200*time.Millisecond)

	if stat.CastCount != 2 {
		t.Errorf("Expected cast count 2, got %d", stat.CastCount)
	}

	if stat.SuccessRate != 0.5 {
		t.Errorf("Expected success rate 0.5, got %f", stat.SuccessRate)
	}

	expectedAvg := (100.0 + 200.0) / 2
	if stat.AvgExecTime != expectedAvg {
		t.Errorf("Expected avg exec time %f, got %f", expectedAvg, stat.AvgExecTime)
	}
}

func TestRecordHotkeyUse(t *testing.T) {
	collector := &Collector{
		enabled: true,
		stats: &Statistics{
			HotkeyStats: make(map[string]*HotkeyStat),
		},
	}

	sequence := "alt+space,e"
	collector.RecordHotkeyUse(sequence)

	stat, ok := collector.stats.HotkeyStats[sequence]
	if !ok {
		t.Fatal("Expected hotkey stat to be created")
	}

	if stat.UseCount != 1 {
		t.Errorf("Expected use count 1, got %d", stat.UseCount)
	}

	if stat.LastUsed.IsZero() {
		t.Error("Expected LastUsed to be set")
	}
}

func TestGetTopSpells(t *testing.T) {
	collector := &Collector{
		stats: &Statistics{
			SpellStats: map[string]*SpellStat{
				"spell1": {Name: "spell1", CastCount: 10},
				"spell2": {Name: "spell2", CastCount: 25},
				"spell3": {Name: "spell3", CastCount: 15},
				"spell4": {Name: "spell4", CastCount: 5},
			},
		},
	}

	top := collector.GetTopSpells(3)

	if len(top) != 3 {
		t.Fatalf("Expected 3 spells, got %d", len(top))
	}

	// Check order
	expected := []string{"spell2", "spell3", "spell1"}
	for i, spell := range top {
		if spell.Name != expected[i] {
			t.Errorf("Expected spell %s at position %d, got %s",
				expected[i], i, spell.Name)
		}
	}
}

func TestSaveLoad(t *testing.T) {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "stats.json")

	// Create collector with data
	collector := &Collector{
		enabled:  true,
		dataFile: dataFile,
		stats: &Statistics{
			Version:         "v1.0.0",
			TotalLaunches:   10,
			TotalSpellsCast: 50,
			SpellStats: map[string]*SpellStat{
				"test": {Name: "test", CastCount: 5},
			},
			HotkeyStats: make(map[string]*HotkeyStat),
			DailyUsage:  make(map[string]*DailyUsage),
		},
	}

	// Save
	collector.save()

	// Create new collector and load
	newCollector := &Collector{
		dataFile: dataFile,
		stats: &Statistics{
			SpellStats:  make(map[string]*SpellStat),
			HotkeyStats: make(map[string]*HotkeyStat),
			DailyUsage:  make(map[string]*DailyUsage),
		},
	}

	if err := newCollector.load(); err != nil {
		t.Fatalf("Failed to load stats: %v", err)
	}

	// Verify loaded data
	if newCollector.stats.TotalLaunches != 10 {
		t.Errorf("Expected 10 launches, got %d", newCollector.stats.TotalLaunches)
	}

	if newCollector.stats.TotalSpellsCast != 50 {
		t.Errorf("Expected 50 spells cast, got %d", newCollector.stats.TotalSpellsCast)
	}

	if spell, ok := newCollector.stats.SpellStats["test"]; !ok || spell.CastCount != 5 {
		t.Error("Expected test spell with 5 casts")
	}
}

func TestDisabledCollector(t *testing.T) {
	collector := &Collector{
		enabled: false,
		stats: &Statistics{
			SpellStats:  make(map[string]*SpellStat),
			HotkeyStats: make(map[string]*HotkeyStat),
			DailyUsage:  make(map[string]*DailyUsage),
		},
	}

	// These should not update stats when disabled
	collector.RecordLaunch()
	collector.RecordSpellCast("test", true, 100*time.Millisecond)
	collector.RecordHotkeyUse("test")

	if collector.stats.TotalLaunches != 0 {
		t.Error("Expected no launches recorded when disabled")
	}

	if collector.stats.TotalSpellsCast != 0 {
		t.Error("Expected no spells recorded when disabled")
	}

	if len(collector.stats.SpellStats) != 0 {
		t.Error("Expected no spell stats when disabled")
	}
}

func TestGenerateReport(t *testing.T) {
	collector := &Collector{
		stats: &Statistics{
			Version:         "v1.0.0",
			InstallDate:     time.Now().Add(-7 * 24 * time.Hour),
			LastUsed:        time.Now(),
			TotalLaunches:   100,
			TotalSpellsCast: 500,
			SpellStats: map[string]*SpellStat{
				"editor":   {Name: "editor", CastCount: 150, SuccessRate: 0.98},
				"terminal": {Name: "terminal", CastCount: 100, SuccessRate: 1.0},
			},
			HotkeyStats: make(map[string]*HotkeyStat),
			DailyUsage:  make(map[string]*DailyUsage),
		},
	}

	report := collector.GenerateReport()

	if report == "" {
		t.Error("Expected non-empty report")
	}

	// Check report contains expected sections
	expectedStrings := []string{
		"=== Spellbook Usage Statistics ===",
		"Version: v1.0.0",
		"Total Launches: 100",
		"Total Spells Cast: 500",
		"Top Spells:",
		"editor",
		"terminal",
	}

	for _, expected := range expectedStrings {
		if !contains(report, expected) {
			t.Errorf("Expected report to contain '%s'", expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func TestStart(t *testing.T) {
	tempDir := t.TempDir()

	collector := &Collector{
		enabled:      true,
		dataFile:     filepath.Join(tempDir, "stats.json"),
		saveInterval: 100 * time.Millisecond,
		stats: &Statistics{
			SpellStats:  make(map[string]*SpellStat),
			HotkeyStats: make(map[string]*HotkeyStat),
			DailyUsage:  make(map[string]*DailyUsage),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	collector.Start(ctx)

	// Should record launch
	time.Sleep(50 * time.Millisecond)
	if collector.stats.TotalLaunches != 1 {
		t.Error("Expected launch to be recorded")
	}

	// Record some activity
	collector.RecordSpellCast("test", true, 100*time.Millisecond)

	// Wait for auto-save
	time.Sleep(200 * time.Millisecond)

	// Cancel to trigger final save
	cancel()
	time.Sleep(50 * time.Millisecond)

	// Check file was created
	if _, err := os.Stat(collector.dataFile); os.IsNotExist(err) {
		t.Error("Expected stats file to be created")
	}
}
