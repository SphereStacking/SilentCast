package updater

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// ProgressReporter reports download progress
type ProgressReporter struct {
	Total      int64
	Downloaded int64
	StartTime  time.Time
	writer     io.Writer
}

// NewProgressReporter creates a new progress reporter
func NewProgressReporter(total int64, writer io.Writer) *ProgressReporter {
	return &ProgressReporter{
		Total:     total,
		StartTime: time.Now(),
		writer:    writer,
	}
}

// Write implements io.Writer for progress tracking
func (pr *ProgressReporter) Write(p []byte) (int, error) {
	n := len(p)
	atomic.AddInt64(&pr.Downloaded, int64(n))
	pr.Report()
	return n, nil
}

// Report prints the current progress
func (pr *ProgressReporter) Report() {
	downloaded := atomic.LoadInt64(&pr.Downloaded)

	if pr.Total <= 0 {
		// Unknown total size
		fmt.Fprintf(pr.writer, "\rDownloading... %s", formatBytes(downloaded))
		return
	}

	// Calculate progress
	percent := float64(downloaded) / float64(pr.Total) * 100

	// Calculate speed
	elapsed := time.Since(pr.StartTime).Seconds()
	speed := float64(downloaded) / elapsed

	// Calculate ETA
	var eta string
	if speed > 0 {
		remaining := pr.Total - downloaded
		etaSeconds := float64(remaining) / speed
		if etaSeconds > 0 {
			eta = formatDuration(time.Duration(etaSeconds * float64(time.Second)))
		}
	}

	// Create progress bar
	barWidth := 30
	filledWidth := int(percent / 100 * float64(barWidth))
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filledWidth {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	// Print progress
	fmt.Fprintf(pr.writer, "\r[%s] %.1f%% %s/%s @ %s/s",
		bar,
		percent,
		formatBytes(downloaded),
		formatBytes(pr.Total),
		formatBytes(int64(speed)),
	)

	if eta != "" {
		fmt.Fprintf(pr.writer, " ETA: %s", eta)
	}
}

// Done prints the final status
func (pr *ProgressReporter) Done() {
	downloaded := atomic.LoadInt64(&pr.Downloaded)
	elapsed := time.Since(pr.StartTime)
	speed := float64(downloaded) / elapsed.Seconds()

	fmt.Fprintf(pr.writer, "\r✅ Downloaded %s in %s @ %s/s\n",
		formatBytes(downloaded),
		formatDuration(elapsed),
		formatBytes(int64(speed)),
	)
}

// formatBytes formats bytes to a human-readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// formatDuration formats a duration to a human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "< 1s"
	}

	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}

	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60

	if minutes < 60 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}

	hours := minutes / 60
	minutes %= 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// ProgressWriter wraps an io.Writer with progress reporting
type ProgressWriter struct {
	writer   io.Writer
	reporter *ProgressReporter
}

// NewProgressWriter creates a new progress writer
func NewProgressWriter(w io.Writer, total int64, output io.Writer) *ProgressWriter {
	return &ProgressWriter{
		writer:   w,
		reporter: NewProgressReporter(total, output),
	}
}

// Write writes data and reports progress
func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.writer.Write(p)
	if err != nil {
		return n, err
	}

	// Report progress
	_, _ = pw.reporter.Write(p[:n]) // Error ignored as this is for progress reporting only
	return n, nil
}

// Done finalizes the progress report
func (pw *ProgressWriter) Done() {
	pw.reporter.Done()
}
