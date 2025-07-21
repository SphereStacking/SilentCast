package config

import "time"

// PerformanceConfig contains performance optimization settings
type PerformanceConfig struct {
	// EnableOptimization enables performance optimizations
	EnableOptimization bool `yaml:"enable_optimization"`

	// BufferSize sets the buffer pool size
	BufferSize int `yaml:"buffer_size"`

	// GCPercent sets the garbage collection target percentage
	GCPercent int `yaml:"gc_percent"`

	// MaxIdleTime sets the maximum idle time for pooled resources
	MaxIdleTime time.Duration `yaml:"max_idle_time"`

	// EnableProfiling enables performance profiling
	EnableProfiling bool `yaml:"enable_profiling"`

	// ProfileHost sets the profiling server host
	ProfileHost string `yaml:"profile_host"`

	// ProfilePort sets the profiling server port
	ProfilePort int `yaml:"profile_port"`
}

// DefaultPerformanceConfig returns default performance configuration
func DefaultPerformanceConfig() PerformanceConfig {
	return PerformanceConfig{
		EnableOptimization: true,
		BufferSize:         1024,
		GCPercent:          100,
		MaxIdleTime:        5 * time.Minute,
		EnableProfiling:    false,
		ProfileHost:        "localhost",
		ProfilePort:        6060,
	}
}

// Validate validates the performance configuration
//
//nolint:unparam // This function validates and modifies config, error return is for interface consistency
func (pc *PerformanceConfig) Validate() error {
	if pc.BufferSize <= 0 {
		pc.BufferSize = 1024
	}
	if pc.GCPercent <= 0 {
		pc.GCPercent = 100
	}
	if pc.MaxIdleTime <= 0 {
		pc.MaxIdleTime = 5 * time.Minute
	}
	if pc.ProfilePort <= 0 || pc.ProfilePort > 65535 {
		pc.ProfilePort = 6060
	}
	return nil
}
