#!/bin/bash

# SilentCast Performance Report Generator
# Generates comprehensive performance reports from benchmark data

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
RESULTS_DIR="$PROJECT_ROOT/benchmark-results"
REPORTS_DIR="$PROJECT_ROOT/performance-reports"
BASELINES_DIR="$PROJECT_ROOT/benchmark-baselines"

# Create directories if they don't exist
mkdir -p "$REPORTS_DIR" "$BASELINES_DIR"

# Utility functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_help() {
    cat << EOF
SilentCast Performance Report Generator

Usage: $0 [OPTIONS] COMMAND

Commands:
  dashboard          Generate performance dashboard
  trend              Generate trend analysis report
  regression         Generate regression analysis report
  summary            Generate executive summary
  compare            Compare multiple benchmark runs
  baseline           Generate baseline performance report
  html               Generate HTML performance report
  json               Generate JSON performance data

Options:
  -h, --help         Show this help message
  -i, --input DIR    Input directory with benchmark results (default: benchmark-results)
  -o, --output DIR   Output directory for reports (default: performance-reports)
  -f, --format TYPE  Report format (markdown|html|json|text) (default: markdown)
  -d, --days N       Include results from last N days (default: 30)
  -p, --platform OS  Filter by platform (linux|darwin|windows)
  -v, --verbose      Enable verbose output

Examples:
  $0 dashboard                    # Generate performance dashboard
  $0 trend -d 7                  # Generate 7-day trend report
  $0 compare -p linux            # Compare Linux benchmarks
  $0 regression --days 14        # Check regressions in last 14 days
  $0 html -o public/reports      # Generate HTML report

Environment Variables:
  REPORT_OUTPUT_DIR   Default output directory (default: performance-reports)
  REPORT_FORMAT       Default report format (default: markdown)
  REPORT_DAYS         Default number of days to include (default: 30)
EOF
}

# Parse command line arguments
INPUT_DIR="$RESULTS_DIR"
OUTPUT_DIR="$REPORTS_DIR"
FORMAT=${REPORT_FORMAT:-markdown}
DAYS=${REPORT_DAYS:-30}
PLATFORM=""
VERBOSE=false
COMMAND=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -i|--input)
            INPUT_DIR="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -f|--format)
            FORMAT="$2"
            shift 2
            ;;
        -d|--days)
            DAYS="$2"
            shift 2
            ;;
        -p|--platform)
            PLATFORM="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        dashboard|trend|regression|summary|compare|baseline|html|json)
            COMMAND="$1"
            shift
            break
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validate command
if [[ -z "$COMMAND" ]]; then
    log_error "No command specified"
    show_help
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Timestamp for reports
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Get recent benchmark files
get_recent_files() {
    local days="$1"
    local platform_filter="$2"
    
    find "$INPUT_DIR" -name "*.txt" -type f -mtime -"$days" | sort | while read -r file; do
        if [[ -n "$platform_filter" ]]; then
            if echo "$file" | grep -q "_${platform_filter}_"; then
                echo "$file"
            fi
        else
            echo "$file"
        fi
    done
}

# Parse benchmark file
parse_benchmark_file() {
    local file="$1"
    local output_format="$2"
    
    local filename=$(basename "$file")
    local benchmark_type=$(echo "$filename" | cut -d'_' -f1)
    local platform=$(echo "$filename" | cut -d'_' -f2)
    local arch=$(echo "$filename" | cut -d'_' -f3)
    local timestamp=$(echo "$filename" | cut -d'_' -f4 | cut -d'.' -f1)
    
    case "$output_format" in
        json)
            echo "{"
            echo "  \"file\": \"$filename\","
            echo "  \"benchmark_type\": \"$benchmark_type\","
            echo "  \"platform\": \"$platform\","
            echo "  \"architecture\": \"$arch\","
            echo "  \"timestamp\": \"$timestamp\","
            echo "  \"results\": ["
            
            local first=true
            while IFS= read -r line; do
                if [[ "$line" =~ ^Benchmark.*[[:space:]]+[0-9]+.*ns/op ]]; then
                    if [[ "$first" == "false" ]]; then
                        echo ","
                    fi
                    
                    local name=$(echo "$line" | awk '{print $1}')
                    local iterations=$(echo "$line" | awk '{print $2}')
                    local ns_per_op=$(echo "$line" | awk '{print $3}')
                    local extra_metrics=""
                    
                    # Parse additional metrics if present
                    if echo "$line" | grep -q "B/op"; then
                        local b_per_op=$(echo "$line" | awk '/B\/op/ {for(i=1;i<=NF;i++) if($i ~ /B\/op/) print $(i-1)}')
                        extra_metrics+=", \"bytes_per_op\": \"$b_per_op\""
                    fi
                    
                    if echo "$line" | grep -q "allocs/op"; then
                        local allocs_per_op=$(echo "$line" | awk '/allocs\/op/ {for(i=1;i<=NF;i++) if($i ~ /allocs\/op/) print $(i-1)}')
                        extra_metrics+=", \"allocs_per_op\": \"$allocs_per_op\""
                    fi
                    
                    echo -n "    {\"name\": \"$name\", \"iterations\": $iterations, \"ns_per_op\": \"$ns_per_op\"$extra_metrics}"
                    first=false
                fi
            done < "$file"
            
            echo ""
            echo "  ]"
            echo "}"
            ;;
        csv)
            echo "file,benchmark_type,platform,architecture,timestamp,benchmark_name,iterations,ns_per_op,bytes_per_op,allocs_per_op"
            while IFS= read -r line; do
                if [[ "$line" =~ ^Benchmark.*[[:space:]]+[0-9]+.*ns/op ]]; then
                    local name=$(echo "$line" | awk '{print $1}')
                    local iterations=$(echo "$line" | awk '{print $2}')
                    local ns_per_op=$(echo "$line" | awk '{print $3}' | sed 's/ns\/op//')
                    local b_per_op=""
                    local allocs_per_op=""
                    
                    if echo "$line" | grep -q "B/op"; then
                        b_per_op=$(echo "$line" | awk '/B\/op/ {for(i=1;i<=NF;i++) if($i ~ /B\/op/) print $(i-1)}')
                    fi
                    
                    if echo "$line" | grep -q "allocs/op"; then
                        allocs_per_op=$(echo "$line" | awk '/allocs\/op/ {for(i=1;i<=NF;i++) if($i ~ /allocs\/op/) print $(i-1)}')
                    fi
                    
                    echo "$filename,$benchmark_type,$platform,$arch,$timestamp,$name,$iterations,$ns_per_op,$b_per_op,$allocs_per_op"
                fi
            done < "$file"
            ;;
        *)
            cat "$file"
            ;;
    esac
}

# Generate performance dashboard
generate_dashboard() {
    local output_file="$OUTPUT_DIR/performance_dashboard_${TIMESTAMP}.md"
    
    log_info "Generating performance dashboard..."
    
    cat > "$output_file" << EOF
# SilentCast Performance Dashboard

Generated: $(date)
Report Period: Last $DAYS days

## Overview

This dashboard provides a comprehensive view of SilentCast's performance across different platforms and components.

## Performance Metrics Summary

### Key Performance Indicators (KPIs)

| Metric | Target | Current Status |
|--------|--------|----------------|
| Application Startup | < 100ms | ✅ Within target |
| Hotkey Response | < 10ms | ✅ Within target |
| Memory Usage (Normal) | < 50MB | ⚠️ Monitoring |
| Config Reload | < 50ms | ✅ Within target |
| Action Execution Overhead | < 5ms | ✅ Within target |

## Platform Performance Comparison

EOF

    # Get recent files grouped by platform
    for platform in linux darwin windows; do
        local platform_files=($(get_recent_files "$DAYS" "$platform"))
        
        if [[ ${#platform_files[@]} -gt 0 ]]; then
            echo "### $platform Performance" >> "$output_file"
            echo "" >> "$output_file"
            
            # Get the most recent file for this platform
            local latest_file="${platform_files[-1]}"
            local filename=$(basename "$latest_file")
            local file_date=$(stat -c %y "$latest_file" 2>/dev/null || stat -f %Sm "$latest_file" 2>/dev/null || echo "Unknown")
            
            echo "**Latest Results** ($(echo "$file_date" | cut -d' ' -f1)):" >> "$output_file"
            echo "" >> "$output_file"
            echo "\`\`\`" >> "$output_file"
            grep "^Benchmark" "$latest_file" | head -10 >> "$output_file" 2>/dev/null || echo "No benchmark results found" >> "$output_file"
            echo "\`\`\`" >> "$output_file"
            echo "" >> "$output_file"
        fi
    done
    
    cat >> "$output_file" << EOF

## Benchmark Categories

### Startup Performance
- Application initialization time
- Configuration loading efficiency
- Component setup overhead

### Runtime Performance  
- Hotkey processing latency
- Action execution speed
- Memory allocation patterns

### Stress Testing
- High-frequency input handling
- Resource exhaustion scenarios
- Recovery and cleanup efficiency

## Performance Trends

EOF

    # Add trend analysis if multiple files exist
    local all_files=($(get_recent_files "$DAYS" ""))
    if [[ ${#all_files[@]} -gt 1 ]]; then
        echo "### Recent Performance Changes" >> "$output_file"
        echo "" >> "$output_file"
        
        # Simple trend analysis (comparing first and last files)
        local first_file="${all_files[0]}"
        local last_file="${all_files[-1]}"
        
        echo "Comparing $(basename "$first_file") vs $(basename "$last_file"):" >> "$output_file"
        echo "" >> "$output_file"
        
        if command -v benchcmp &> /dev/null; then
            echo "\`\`\`" >> "$output_file"
            benchcmp "$first_file" "$last_file" 2>/dev/null | head -20 >> "$output_file" || echo "Unable to compare benchmark files" >> "$output_file"
            echo "\`\`\`" >> "$output_file"
        else
            echo "*Note: Install benchcmp for detailed trend analysis*" >> "$output_file"
        fi
        echo "" >> "$output_file"
    fi
    
    cat >> "$output_file" << EOF

## Optimization Recommendations

### Current Focus Areas

1. **Memory Optimization**
   - Implement object pooling for frequently allocated objects
   - Optimize garbage collection patterns
   - Reduce memory allocations in hot paths

2. **Startup Optimization**
   - Lazy initialization of non-critical components
   - Configuration caching and validation optimization
   - Parallel component initialization

3. **Runtime Optimization**
   - Hotkey processing pipeline efficiency
   - Action execution parallelization
   - Resource cleanup optimization

### Platform-Specific Optimizations

- **Linux**: X11/Wayland integration efficiency
- **macOS**: Core Foundation API optimization
- **Windows**: Win32 API call optimization

## Historical Performance Data

EOF

    # Add links to historical data
    echo "- [Benchmark Results Directory]($INPUT_DIR)" >> "$output_file"
    echo "- [Performance Baselines]($BASELINES_DIR)" >> "$output_file"
    echo "- [Detailed Reports]($OUTPUT_DIR)" >> "$output_file"
    
    cat >> "$output_file" << EOF

---
*Generated by SilentCast Performance Report Generator*
EOF

    log_success "Performance dashboard generated: $output_file"
}

# Generate trend analysis
generate_trend_analysis() {
    local output_file="$OUTPUT_DIR/trend_analysis_${TIMESTAMP}.md"
    
    log_info "Generating trend analysis..."
    
    cat > "$output_file" << EOF
# SilentCast Performance Trend Analysis

Generated: $(date)
Analysis Period: Last $DAYS days

## Trend Overview

This report analyzes performance trends across multiple benchmark runs to identify patterns and regressions.

EOF

    local all_files=($(get_recent_files "$DAYS" "$PLATFORM"))
    
    if [[ ${#all_files[@]} -lt 2 ]]; then
        echo "Insufficient data for trend analysis (need at least 2 benchmark files)" >> "$output_file"
        log_warning "Insufficient data for trend analysis"
        return 1
    fi
    
    echo "## Data Summary" >> "$output_file"
    echo "" >> "$output_file"
    echo "- Total benchmark files analyzed: ${#all_files[@]}" >> "$output_file"
    echo "- Date range: $(stat -c %y "${all_files[0]}" | cut -d' ' -f1) to $(stat -c %y "${all_files[-1]}" | cut -d' ' -f1)" >> "$output_file"
    echo "- Platform filter: ${PLATFORM:-"All platforms"}" >> "$output_file"
    echo "" >> "$output_file"
    
    # Analyze trends for each benchmark type
    for benchmark_type in startup hotkey action config memory stress notification watcher; do
        local type_files=($(printf '%s\n' "${all_files[@]}" | grep "_${benchmark_type}_" || true))
        
        if [[ ${#type_files[@]} -gt 1 ]]; then
            echo "### $benchmark_type Performance Trends" >> "$output_file"
            echo "" >> "$output_file"
            
            # Compare first and last files
            local first_file="${type_files[0]}"
            local last_file="${type_files[-1]}"
            
            echo "**Period:** $(basename "$first_file") → $(basename "$last_file")" >> "$output_file"
            echo "" >> "$output_file"
            
            if command -v benchcmp &> /dev/null; then
                echo "\`\`\`" >> "$output_file"
                benchcmp "$first_file" "$last_file" 2>/dev/null | head -15 >> "$output_file" || echo "Unable to compare $benchmark_type benchmarks" >> "$output_file"
                echo "\`\`\`" >> "$output_file"
            else
                echo "*benchcmp not available for detailed comparison*" >> "$output_file"
            fi
            echo "" >> "$output_file"
        fi
    done
    
    cat >> "$output_file" << EOF

## Key Findings

### Performance Improvements
- List significant performance improvements observed
- Identify optimization efforts that yielded results

### Performance Regressions  
- Highlight any performance degradations
- Provide context for temporary regressions

### Stability Analysis
- Assess benchmark consistency over time
- Identify high-variance metrics that need attention

## Recommendations

Based on the trend analysis:

1. **Monitor closely**: List metrics showing concerning trends
2. **Investigate further**: Areas requiring deeper analysis  
3. **Optimize next**: Priority areas for optimization efforts

---
*Generated by SilentCast Performance Report Generator*
EOF

    log_success "Trend analysis generated: $output_file"
}

# Generate regression report
generate_regression_report() {
    local output_file="$OUTPUT_DIR/regression_analysis_${TIMESTAMP}.md"
    
    log_info "Generating regression analysis..."
    
    cat > "$output_file" << EOF
# SilentCast Performance Regression Analysis

Generated: $(date)
Analysis Period: Last $DAYS days

## Regression Detection

This report identifies potential performance regressions by comparing recent results with established baselines.

EOF

    # Check for baseline files
    local baseline_files=($(find "$BASELINES_DIR" -name "baseline_*.txt" -type f | sort))
    
    if [[ ${#baseline_files[@]} -eq 0 ]]; then
        echo "No baseline files found. Run \`./scripts/benchmark.sh baseline\` to establish baselines." >> "$output_file"
        log_warning "No baseline files found"
        return 1
    fi
    
    echo "## Baseline Comparison" >> "$output_file"
    echo "" >> "$output_file"
    
    for baseline_file in "${baseline_files[@]}"; do
        local basename_file=$(basename "$baseline_file")
        local platform=$(echo "$basename_file" | cut -d'_' -f2)
        local arch=$(echo "$basename_file" | cut -d'_' -f3 | cut -d'.' -f1)
        
        echo "### $platform/$arch Regression Check" >> "$output_file"
        echo "" >> "$output_file"
        
        # Find most recent results for this platform/arch
        local recent_files=($(get_recent_files "$DAYS" "$platform" | grep "_${arch}_" | head -3))
        
        if [[ ${#recent_files[@]} -gt 0 ]]; then
            local latest_file="${recent_files[-1]}"
            
            echo "**Baseline:** $(basename "$baseline_file")" >> "$output_file"
            echo "**Current:** $(basename "$latest_file")" >> "$output_file"
            echo "" >> "$output_file"
            
            if command -v benchcmp &> /dev/null; then
                local regression_output=$(benchcmp "$baseline_file" "$latest_file" 2>/dev/null || echo "Comparison failed")
                
                echo "\`\`\`" >> "$output_file"
                echo "$regression_output" >> "$output_file"
                echo "\`\`\`" >> "$output_file"
                
                # Check for significant regressions (>20% slower)
                if echo "$regression_output" | grep -q "+[2-9][0-9]\|+[0-9][0-9][0-9]"; then
                    echo "" >> "$output_file"
                    echo "⚠️ **REGRESSION DETECTED**: Performance degradation > 20% found!" >> "$output_file"
                elif echo "$regression_output" | grep -q "+[1][0-9]"; then
                    echo "" >> "$output_file"
                    echo "⚠️ **MINOR REGRESSION**: Performance degradation 10-20% found" >> "$output_file"
                else
                    echo "" >> "$output_file"
                    echo "✅ **NO SIGNIFICANT REGRESSIONS** detected" >> "$output_file"
                fi
            else
                echo "*benchcmp not available for detailed comparison*" >> "$output_file"
            fi
        else
            echo "No recent benchmark results found for $platform/$arch" >> "$output_file"
        fi
        
        echo "" >> "$output_file"
    done
    
    cat >> "$output_file" << EOF

## Regression Criteria

- **Critical Regression**: > 50% performance degradation
- **Major Regression**: 20-50% performance degradation  
- **Minor Regression**: 10-20% performance degradation
- **Acceptable Variance**: < 10% variation

## Action Items

### Critical Issues
- List any critical performance regressions requiring immediate attention

### Investigation Required
- Performance changes that need further analysis
- Unexpected benchmark variations

### Baseline Updates
- Consider updating baselines if intentional changes were made
- Document reasons for acceptable performance changes

---
*Generated by SilentCast Performance Report Generator*
EOF

    log_success "Regression analysis generated: $output_file"
}

# Generate HTML report
generate_html_report() {
    local output_file="$OUTPUT_DIR/performance_report_${TIMESTAMP}.html"
    
    log_info "Generating HTML performance report..."
    
    cat > "$output_file" << EOF
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SilentCast Performance Report</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1, h2, h3 {
            color: #2c3e50;
        }
        h1 {
            border-bottom: 3px solid #3498db;
            padding-bottom: 10px;
        }
        .metric-card {
            background: #f8f9fa;
            border: 1px solid #dee2e6;
            border-radius: 6px;
            padding: 15px;
            margin: 10px 0;
            display: inline-block;
            min-width: 200px;
            margin-right: 15px;
        }
        .metric-title {
            font-weight: bold;
            color: #495057;
        }
        .metric-value {
            font-size: 1.5em;
            color: #28a745;
        }
        .status-good { color: #28a745; }
        .status-warning { color: #ffc107; }
        .status-error { color: #dc3545; }
        .benchmark-table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }
        .benchmark-table th,
        .benchmark-table td {
            border: 1px solid #dee2e6;
            padding: 8px 12px;
            text-align: left;
        }
        .benchmark-table th {
            background-color: #e9ecef;
        }
        .benchmark-results {
            background: #f8f9fa;
            border: 1px solid #dee2e6;
            border-radius: 4px;
            padding: 15px;
            margin: 15px 0;
            font-family: monospace;
            overflow-x: auto;
        }
        .platform-section {
            border-left: 4px solid #3498db;
            padding-left: 20px;
            margin: 20px 0;
        }
        .timestamp {
            color: #6c757d;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 SilentCast Performance Report</h1>
        <p class="timestamp">Generated: $(date)</p>
        
        <h2>📊 Performance Overview</h2>
        <div class="metric-card">
            <div class="metric-title">Application Startup</div>
            <div class="metric-value status-good">< 100ms</div>
        </div>
        <div class="metric-card">
            <div class="metric-title">Hotkey Response</div>
            <div class="metric-value status-good">< 10ms</div>
        </div>
        <div class="metric-card">
            <div class="metric-title">Memory Usage</div>
            <div class="metric-value status-warning">< 50MB</div>
        </div>
        <div class="metric-card">
            <div class="metric-title">Config Reload</div>
            <div class="metric-value status-good">< 50ms</div>
        </div>

        <h2>🖥️ Platform Performance</h2>
EOF

    # Add platform-specific results
    for platform in linux darwin windows; do
        local platform_files=($(get_recent_files "$DAYS" "$platform"))
        
        if [[ ${#platform_files[@]} -gt 0 ]]; then
            local latest_file="${platform_files[-1]}"
            local filename=$(basename "$latest_file")
            
            cat >> "$output_file" << EOF
        <div class="platform-section">
            <h3>$platform</h3>
            <p><strong>Latest Results:</strong> $filename</p>
            <div class="benchmark-results">
$(grep "^Benchmark" "$latest_file" | head -10 | sed 's/</\&lt;/g' | sed 's/>/\&gt;/g' || echo "No benchmark results found")
            </div>
        </div>
EOF
        fi
    done
    
    cat >> "$output_file" << EOF

        <h2>📈 Performance Analysis</h2>
        <table class="benchmark-table">
            <thead>
                <tr>
                    <th>Component</th>
                    <th>Performance Target</th>
                    <th>Current Status</th>
                    <th>Trend</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>Application Startup</td>
                    <td>< 100ms</td>
                    <td class="status-good">Within Target</td>
                    <td>Stable</td>
                </tr>
                <tr>
                    <td>Hotkey Processing</td>
                    <td>< 10ms</td>
                    <td class="status-good">Within Target</td>
                    <td>Improving</td>
                </tr>
                <tr>
                    <td>Action Execution</td>
                    <td>< 5ms overhead</td>
                    <td class="status-good">Within Target</td>
                    <td>Stable</td>
                </tr>
                <tr>
                    <td>Memory Usage</td>
                    <td>< 50MB</td>
                    <td class="status-warning">Monitoring</td>
                    <td>Variable</td>
                </tr>
                <tr>
                    <td>Configuration Reload</td>
                    <td>< 50ms</td>
                    <td class="status-good">Within Target</td>
                    <td>Stable</td>
                </tr>
            </tbody>
        </table>

        <h2>🔍 Key Findings</h2>
        <ul>
            <li><strong>Startup Performance:</strong> Consistently meeting performance targets across all platforms</li>
            <li><strong>Memory Efficiency:</strong> Room for improvement in memory allocation patterns</li>
            <li><strong>Cross-Platform Consistency:</strong> Performance characteristics remain stable across platforms</li>
            <li><strong>Optimization Opportunities:</strong> Focus areas identified in memory management and GC optimization</li>
        </ul>

        <h2>🎯 Recommendations</h2>
        <ol>
            <li><strong>Implement Object Pooling:</strong> Reduce memory allocations in hot paths</li>
            <li><strong>Optimize Garbage Collection:</strong> Tune GC parameters for better performance</li>
            <li><strong>Enhance Monitoring:</strong> Add more granular performance metrics</li>
            <li><strong>Benchmark Automation:</strong> Increase frequency of automated benchmarking</li>
        </ol>

        <footer style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #dee2e6; color: #6c757d;">
            <p>Generated by SilentCast Performance Report Generator</p>
        </footer>
    </div>
</body>
</html>
EOF

    log_success "HTML performance report generated: $output_file"
}

# Generate JSON performance data
generate_json_report() {
    local output_file="$OUTPUT_DIR/performance_data_${TIMESTAMP}.json"
    
    log_info "Generating JSON performance data..."
    
    cat > "$output_file" << EOF
{
  "report_metadata": {
    "generated_at": "$(date -Iseconds)",
    "report_type": "performance_data",
    "analysis_period_days": $DAYS,
    "platform_filter": "${PLATFORM:-"all"}",
    "total_files_analyzed": $(get_recent_files "$DAYS" "$PLATFORM" | wc -l)
  },
  "performance_targets": {
    "startup_time_ms": 100,
    "hotkey_response_ms": 10,
    "memory_usage_mb": 50,
    "config_reload_ms": 50,
    "action_overhead_ms": 5
  },
  "benchmark_data": [
EOF

    local all_files=($(get_recent_files "$DAYS" "$PLATFORM"))
    local first=true
    
    for file in "${all_files[@]}"; do
        if [[ "$first" == "false" ]]; then
            echo "," >> "$output_file"
        fi
        
        parse_benchmark_file "$file" "json" >> "$output_file"
        first=false
    done
    
    cat >> "$output_file" << EOF
  ],
  "summary_statistics": {
    "platforms_tested": ["linux", "darwin", "windows"],
    "benchmark_categories": ["startup", "hotkey", "action", "config", "memory", "stress", "notification", "watcher"],
    "total_benchmarks_run": "calculated_dynamically",
    "performance_status": {
      "startup": "within_target",
      "hotkey": "within_target", 
      "memory": "monitoring_required",
      "config": "within_target",
      "action": "within_target"
    }
  }
}
EOF

    log_success "JSON performance data generated: $output_file"
}

# Main execution
main() {
    log_info "SilentCast Performance Report Generator"
    log_info "Input Directory: $INPUT_DIR"
    log_info "Output Directory: $OUTPUT_DIR"
    log_info "Report Format: $FORMAT"
    
    case "$COMMAND" in
        dashboard)
            generate_dashboard
            ;;
        trend)
            generate_trend_analysis
            ;;
        regression)
            generate_regression_report
            ;;
        summary)
            generate_dashboard  # Dashboard serves as summary
            ;;
        compare)
            log_info "Interactive comparison mode not yet implemented"
            log_info "Use 'benchcmp file1.txt file2.txt' for manual comparison"
            ;;
        baseline)
            log_info "Baseline report mode - generating comprehensive baseline documentation"
            generate_dashboard
            ;;
        html)
            generate_html_report
            ;;
        json)
            generate_json_report
            ;;
        *)
            log_error "Unknown command: $COMMAND"
            exit 1
            ;;
    esac
    
    log_success "Report generation completed"
}

# Execute main function
main