package stats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/SphereStacking/silentcast/pkg/logger"
)

// Collector collects usage statistics
type Collector struct {
	mu           sync.RWMutex
	stats        *Statistics
	dataFile     string
	saveInterval time.Duration
	enabled      bool
}

// Statistics holds usage data
type Statistics struct {
	Version         string                 `json:"version"`
	InstallDate     time.Time              `json:"install_date"`
	LastUsed        time.Time              `json:"last_used"`
	TotalLaunches   int                    `json:"total_launches"`
	TotalSpellsCast int                    `json:"total_spells_cast"`
	SpellStats      map[string]*SpellStat  `json:"spell_stats"`
	HotkeyStats     map[string]*HotkeyStat `json:"hotkey_stats"`
	DailyUsage      map[string]*DailyUsage `json:"daily_usage"`
}

// SpellStat tracks statistics for a single spell
type SpellStat struct {
	Name        string    `json:"name"`
	CastCount   int       `json:"cast_count"`
	LastCast    time.Time `json:"last_cast"`
	SuccessRate float64   `json:"success_rate"`
	AvgExecTime float64   `json:"avg_exec_time_ms"`
}

// HotkeyStat tracks hotkey usage
type HotkeyStat struct {
	Sequence string    `json:"sequence"`
	UseCount int       `json:"use_count"`
	LastUsed time.Time `json:"last_used"`
}

// DailyUsage tracks usage per day
type DailyUsage struct {
	Date       string `json:"date"`
	Launches   int    `json:"launches"`
	SpellsCast int    `json:"spells_cast"`
	ActiveTime int    `json:"active_time_minutes"`
}

// Config holds collector configuration
type Config struct {
	Enabled      bool
	DataFile     string
	SaveInterval time.Duration
}

// NewCollector creates a new statistics collector
func NewCollector(cfg Config) *Collector {
	c := &Collector{
		enabled:      cfg.Enabled,
		dataFile:     cfg.DataFile,
		saveInterval: cfg.SaveInterval,
		stats: &Statistics{
			SpellStats:  make(map[string]*SpellStat),
			HotkeyStats: make(map[string]*HotkeyStat),
			DailyUsage:  make(map[string]*DailyUsage),
		},
	}

	if cfg.SaveInterval == 0 {
		c.saveInterval = 5 * time.Minute
	}

	// Load existing stats if available
	if err := c.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Warn("Failed to load statistics: %v", err)
	}

	// Initialize if new installation
	if c.stats.InstallDate.IsZero() {
		c.stats.InstallDate = time.Now()
	}

	return c
}

// Start begins collecting statistics
func (c *Collector) Start(ctx context.Context) {
	if !c.enabled {
		logger.Info("Statistics collection disabled")
		return
	}

	// Record launch
	c.RecordLaunch()

	// Periodic save
	go func() {
		ticker := time.NewTicker(c.saveInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				c.save()
				return
			case <-ticker.C:
				c.save()
			}
		}
	}()
}

// RecordLaunch records application launch
func (c *Collector) RecordLaunch() {
	if !c.enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats.TotalLaunches++
	c.stats.LastUsed = time.Now()

	// Update daily usage
	today := time.Now().Format("2006-01-02")
	if daily, ok := c.stats.DailyUsage[today]; ok {
		daily.Launches++
	} else {
		c.stats.DailyUsage[today] = &DailyUsage{
			Date:     today,
			Launches: 1,
		}
	}
}

// RecordSpellCast records a spell execution
func (c *Collector) RecordSpellCast(spell string, success bool, execTime time.Duration) {
	if !c.enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats.TotalSpellsCast++

	// Update spell stats
	stat, ok := c.stats.SpellStats[spell]
	if !ok {
		stat = &SpellStat{
			Name: spell,
		}
		c.stats.SpellStats[spell] = stat
	}

	stat.CastCount++
	stat.LastCast = time.Now()

	// Update success rate
	if success {
		stat.SuccessRate = ((stat.SuccessRate * float64(stat.CastCount-1)) + 1.0) / float64(stat.CastCount)
	} else {
		stat.SuccessRate = (stat.SuccessRate * float64(stat.CastCount-1)) / float64(stat.CastCount)
	}

	// Update average execution time
	execMs := float64(execTime.Milliseconds())
	stat.AvgExecTime = ((stat.AvgExecTime * float64(stat.CastCount-1)) + execMs) / float64(stat.CastCount)

	// Update daily usage
	today := time.Now().Format("2006-01-02")
	if daily, ok := c.stats.DailyUsage[today]; ok {
		daily.SpellsCast++
	} else {
		c.stats.DailyUsage[today] = &DailyUsage{
			Date:       today,
			SpellsCast: 1,
		}
	}
}

// RecordHotkeyUse records hotkey usage
func (c *Collector) RecordHotkeyUse(sequence string) {
	if !c.enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	stat, ok := c.stats.HotkeyStats[sequence]
	if !ok {
		stat = &HotkeyStat{
			Sequence: sequence,
		}
		c.stats.HotkeyStats[sequence] = stat
	}

	stat.UseCount++
	stat.LastUsed = time.Now()
}

// GetStatistics returns a copy of current statistics
func (c *Collector) GetStatistics() Statistics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Deep copy to avoid race conditions
	statsCopy := *c.stats
	statsCopy.SpellStats = make(map[string]*SpellStat)
	statsCopy.HotkeyStats = make(map[string]*HotkeyStat)
	statsCopy.DailyUsage = make(map[string]*DailyUsage)

	for k, v := range c.stats.SpellStats {
		statCopy := *v
		statsCopy.SpellStats[k] = &statCopy
	}

	for k, v := range c.stats.HotkeyStats {
		statCopy := *v
		statsCopy.HotkeyStats[k] = &statCopy
	}

	for k, v := range c.stats.DailyUsage {
		statCopy := *v
		statsCopy.DailyUsage[k] = &statCopy
	}

	return statsCopy
}

// GetTopSpells returns the most used spells
func (c *Collector) GetTopSpells(limit int) []SpellStat {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Convert to slice and sort
	spells := make([]SpellStat, 0, len(c.stats.SpellStats))
	for _, stat := range c.stats.SpellStats {
		spells = append(spells, *stat)
	}

	// Simple bubble sort for small dataset
	for i := 0; i < len(spells); i++ {
		for j := i + 1; j < len(spells); j++ {
			if spells[j].CastCount > spells[i].CastCount {
				spells[i], spells[j] = spells[j], spells[i]
			}
		}
	}

	if limit > len(spells) {
		limit = len(spells)
	}

	return spells[:limit]
}

// GenerateReport generates a usage report
func (c *Collector) GenerateReport() string {
	stats := c.GetStatistics()

	report := "=== Spellbook Usage Statistics ===\n\n"
	report += fmt.Sprintf("Version: %s\n", stats.Version)
	report += fmt.Sprintf("Installed: %s\n", stats.InstallDate.Format("2006-01-02"))
	report += fmt.Sprintf("Last Used: %s\n", stats.LastUsed.Format("2006-01-02 15:04:05"))
	report += fmt.Sprintf("Total Launches: %d\n", stats.TotalLaunches)
	report += fmt.Sprintf("Total Spells Cast: %d\n\n", stats.TotalSpellsCast)

	// Top spells
	topSpells := c.GetTopSpells(10)
	if len(topSpells) > 0 {
		report += "Top Spells:\n"
		for i, spell := range topSpells {
			report += fmt.Sprintf("%d. %s - %d casts (%.1f%% success)\n",
				i+1, spell.Name, spell.CastCount, spell.SuccessRate*100)
		}
		report += "\n"
	}

	// Recent daily usage
	report += "Recent Daily Usage:\n"
	// TODO: Sort and show last 7 days

	return report
}

// Reset clears all statistics
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats = &Statistics{
		InstallDate: c.stats.InstallDate, // Preserve install date
		SpellStats:  make(map[string]*SpellStat),
		HotkeyStats: make(map[string]*HotkeyStat),
		DailyUsage:  make(map[string]*DailyUsage),
	}

	c.save()
}

// load loads statistics from file
func (c *Collector) load() error {
	if c.dataFile == "" {
		return nil
	}

	data, err := os.ReadFile(c.dataFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c.stats)
}

// save saves statistics to file
func (c *Collector) save() {
	if c.dataFile == "" {
		return
	}

	c.mu.RLock()
	data, err := json.MarshalIndent(c.stats, "", "  ")
	c.mu.RUnlock()

	if err != nil {
		logger.Error("Failed to marshal statistics: %v", err)
		return
	}

	// Ensure directory exists
	dir := filepath.Dir(c.dataFile)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logger.Error("Failed to create stats directory: %v", err)
		return
	}

	// Write atomically
	tmpFile := c.dataFile + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0o600); err != nil {
		logger.Error("Failed to write statistics: %v", err)
		return
	}

	if err := os.Rename(tmpFile, c.dataFile); err != nil {
		logger.Error("Failed to save statistics: %v", err)
		os.Remove(tmpFile)
	}
}
