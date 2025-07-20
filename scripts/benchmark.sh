#!/bin/bash

# SilentCast Performance Benchmark Runner
# Provides a comprehensive interface for running various benchmark scenarios

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
APP_DIR="$PROJECT_ROOT/app"
RESULTS_DIR="$PROJECT_ROOT/benchmark-results"
REPORTS_DIR="$PROJECT_ROOT/performance-reports"
BUILD_TAGS="nogohook notray"

# Create directories if they don't exist
mkdir -p "$RESULTS_DIR" "$REPORTS_DIR"

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

log_benchmark() {
    echo -e "${PURPLE}[BENCHMARK]${NC} $1"
}

show_help() {
    cat << EOF
SilentCast Performance Benchmark Runner

Usage: $0 [OPTIONS] COMMAND

Commands:
  all                 Run all benchmarks (comprehensive suite)
  startup            Run startup performance benchmarks
  hotkey             Run hotkey processing benchmarks
  action             Run action execution benchmarks
  config             Run configuration loading benchmarks
  memory             Run memory usage benchmarks
  stress             Run stress test benchmarks
  notification       Run notification system benchmarks
  watcher            Run file watcher benchmarks
  compare            Compare two benchmark results
  report             Generate performance report
  baseline           Update performance baseline
  regression         Check for performance regressions

Options:
  -h, --help         Show this help message
  -v, --verbose      Enable verbose output
  -o, --output DIR   Output directory (default: benchmark-results)
  -t, --timeout SEC  Benchmark timeout in seconds (default: 1800)
  -c, --cpu-profile  Enable CPU profiling
  -m, --mem-profile  Enable memory profiling
  -r, --runs N       Number of benchmark runs (default: 1)
  -f, --format TYPE  Output format (text|json|csv) (default: text)
  --tags TAGS        Build tags (default: "$BUILD_TAGS")

Examples:
  $0 all                          # Run all benchmarks
  $0 startup -c -m               # Run startup benchmarks with profiling
  $0 compare baseline.txt new.txt # Compare two benchmark results
  $0 stress -t 3600              # Run stress tests for 1 hour
  $0 memory --runs 5             # Run memory benchmarks 5 times

Environment Variables:
  BENCHMARK_TIMEOUT   Default timeout in seconds (default: 1800)
  BENCHMARK_RUNS      Default number of runs (default: 1)
  BENCHMARK_FORMAT    Default output format (default: text)
EOF
}

# Parse command line arguments
VERBOSE=false
OUTPUT_DIR="$RESULTS_DIR"
TIMEOUT=${BENCHMARK_TIMEOUT:-1800}
CPU_PROFILE=false
MEM_PROFILE=false
RUNS=${BENCHMARK_RUNS:-1}
FORMAT=${BENCHMARK_FORMAT:-text}
COMMAND=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -c|--cpu-profile)
            CPU_PROFILE=true
            shift
            ;;
        -m|--mem-profile)
            MEM_PROFILE=true
            shift
            ;;
        -r|--runs)
            RUNS="$2"
            shift 2
            ;;
        -f|--format)
            FORMAT="$2"
            shift 2
            ;;
        --tags)
            BUILD_TAGS="$2"
            shift 2
            ;;
        all|startup|hotkey|action|config|memory|stress|notification|watcher|compare|report|baseline|regression)
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

# Timestamp for unique filenames
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Build profiling flags
PROFILE_FLAGS=""
if [[ "$CPU_PROFILE" == "true" ]]; then
    PROFILE_FLAGS="$PROFILE_FLAGS -cpuprofile=$OUTPUT_DIR/cpu_${COMMAND}_${TIMESTAMP}.prof"
fi

if [[ "$MEM_PROFILE" == "true" ]]; then
    PROFILE_FLAGS="$PROFILE_FLAGS -memprofile=$OUTPUT_DIR/mem_${COMMAND}_${TIMESTAMP}.prof"
fi

# System information collection
collect_system_info() {
    local info_file="$OUTPUT_DIR/system_info_${TIMESTAMP}.txt"
    
    log_info "Collecting system information..."
    
    cat > "$info_file" << EOF
# SilentCast Benchmark System Information
# Generated: $(date)

## System Details
Platform: $PLATFORM
Architecture: $ARCH
Hostname: $(hostname)
Kernel: $(uname -r)

## Go Environment
$(go version)
GOROOT: $(go env GOROOT)
GOPATH: $(go env GOPATH)
GOOS: $(go env GOOS)
GOARCH: $(go env GOARCH)

## Hardware Information
EOF

    # Platform-specific hardware info
    case "$PLATFORM" in
        linux)
            echo "CPU Info:" >> "$info_file"
            grep "model name\|cpu cores\|siblings" /proc/cpuinfo | head -20 >> "$info_file" || true
            echo "" >> "$info_file"
            echo "Memory Info:" >> "$info_file"
            grep "MemTotal\|MemAvailable" /proc/meminfo >> "$info_file" || true
            ;;
        darwin)
            echo "CPU Info:" >> "$info_file"
            sysctl -n machdep.cpu.brand_string >> "$info_file" || true
            echo "CPU Cores: $(sysctl -n hw.ncpu)" >> "$info_file" || true
            echo "Memory: $(( $(sysctl -n hw.memsize) / 1024 / 1024 / 1024 ))GB" >> "$info_file" || true
            ;;
        *)
            echo "System info collection not implemented for $PLATFORM" >> "$info_file"
            ;;
    esac

    echo "" >> "$info_file"
    echo "## Benchmark Configuration" >> "$info_file"
    echo "Build Tags: $BUILD_TAGS" >> "$info_file"
    echo "Timeout: ${TIMEOUT}s" >> "$info_file"
    echo "Runs: $RUNS" >> "$info_file"
    echo "CPU Profiling: $CPU_PROFILE" >> "$info_file"
    echo "Memory Profiling: $MEM_PROFILE" >> "$info_file"
    
    log_success "System information saved to $info_file"
}

# Run benchmark function
run_benchmark() {
    local benchmark_type="$1"
    local output_file="$OUTPUT_DIR/${benchmark_type}_${PLATFORM}_${ARCH}_${TIMESTAMP}.txt"
    
    log_benchmark "Running $benchmark_type benchmarks..."
    
    cd "$APP_DIR"
    
    local start_time=$(date +%s)
    
    case "$benchmark_type" in
        all)
            go test -tags "$BUILD_TAGS" -benchmem -timeout=${TIMEOUT}s $PROFILE_FLAGS -bench=. ./test/benchmarks/... | tee "$output_file"
            ;;
        startup)
            go test -tags "$BUILD_TAGS" -benchmem -timeout=${TIMEOUT}s $PROFILE_FLAGS -bench=BenchmarkStartup ./test/benchmarks/startup_test.go ./test/benchmarks/framework.go | tee "$output_file"
            ;;
        hotkey)
            go test -tags "$BUILD_TAGS" -benchmem -timeout=${TIMEOUT}s $PROFILE_FLAGS -bench=BenchmarkKey ./test/benchmarks/hotkey_test.go ./test/benchmarks/framework.go | tee "$output_file"
            ;;
        action)
            go test -tags "$BUILD_TAGS" -benchmem -timeout=${TIMEOUT}s $PROFILE_FLAGS -bench=BenchmarkAction ./test/benchmarks/action_test.go ./test/benchmarks/framework.go | tee "$output_file"
            ;;
        config)
            go test -tags "$BUILD_TAGS" -benchmem -timeout=${TIMEOUT}s $PROFILE_FLAGS -bench=BenchmarkConfig ./test/benchmarks/config_test.go ./test/benchmarks/framework.go | tee "$output_file"
            ;;
        memory)
            go test -tags "$BUILD_TAGS" -benchmem -timeout=${TIMEOUT}s $PROFILE_FLAGS -bench=BenchmarkMemory ./test/benchmarks/memory_test.go ./test/benchmarks/framework.go | tee "$output_file"
            ;;
        stress)
            go test -tags "$BUILD_TAGS" -benchmem -timeout=${TIMEOUT}s $PROFILE_FLAGS -bench=BenchmarkStress ./test/benchmarks/stress_test.go ./test/benchmarks/framework.go | tee "$output_file"
            ;;
        notification)
            go test -tags "$BUILD_TAGS" -benchmem -timeout=${TIMEOUT}s $PROFILE_FLAGS -bench=BenchmarkNotification ./test/benchmarks/notification_test.go ./test/benchmarks/framework.go | tee "$output_file"
            ;;
        watcher)
            go test -tags "$BUILD_TAGS" -benchmem -timeout=${TIMEOUT}s $PROFILE_FLAGS -bench=BenchmarkConfig ./test/benchmarks/watcher_test.go ./test/benchmarks/framework.go | tee "$output_file"
            ;;
        *)
            log_error "Unknown benchmark type: $benchmark_type"
            return 1
            ;;
    esac
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    log_success "$benchmark_type benchmarks completed in ${duration}s"
    log_info "Results saved to: $output_file"
    
    # Convert format if needed
    if [[ "$FORMAT" != "text" ]]; then
        convert_format "$output_file" "$FORMAT"
    fi
}

# Format conversion function
convert_format() {
    local input_file="$1"
    local format="$2"
    local base_name="${input_file%.*}"
    
    case "$format" in
        json)
            local json_file="${base_name}.json"
            log_info "Converting to JSON format: $json_file"
            
            # Simple JSON conversion (this could be enhanced with a proper parser)
            echo "{" > "$json_file"
            echo '  "timestamp": "'$(date -Iseconds)'",' >> "$json_file"
            echo '  "platform": "'$PLATFORM'",' >> "$json_file"
            echo '  "architecture": "'$ARCH'",' >> "$json_file"
            echo '  "benchmarks": [' >> "$json_file"
            
            local first=true
            while IFS= read -r line; do
                if [[ "$line" =~ ^Benchmark.*[[:space:]]+[0-9]+.*ns/op ]]; then
                    if [[ "$first" == "false" ]]; then
                        echo ',' >> "$json_file"
                    fi
                    
                    # Parse benchmark line (this is a simple parser)
                    local name=$(echo "$line" | awk '{print $1}')
                    local iterations=$(echo "$line" | awk '{print $2}')
                    local ns_per_op=$(echo "$line" | awk '{print $3}')
                    
                    echo -n '    {"name": "'$name'", "iterations": '$iterations', "ns_per_op": "'$ns_per_op'"}' >> "$json_file"
                    first=false
                fi
            done < "$input_file"
            
            echo '' >> "$json_file"
            echo '  ]' >> "$json_file"
            echo '}' >> "$json_file"
            ;;
        csv)
            local csv_file="${base_name}.csv"
            log_info "Converting to CSV format: $csv_file"
            
            echo "name,iterations,ns_per_op,mb_per_sec,allocs_per_op,b_per_op" > "$csv_file"
            
            while IFS= read -r line; do
                if [[ "$line" =~ ^Benchmark.*[[:space:]]+[0-9]+.*ns/op ]]; then
                    # Parse and format for CSV (basic implementation)
                    echo "$line" | awk '{print $1","$2","$3","$4","$5","$6}' >> "$csv_file"
                fi
            done < "$input_file"
            ;;
        *)
            log_warning "Unknown format: $format"
            ;;
    esac
}

# Compare benchmark results
compare_benchmarks() {
    local file1="$1"
    local file2="$2"
    
    if [[ ! -f "$file1" ]]; then
        log_error "First benchmark file not found: $file1"
        return 1
    fi
    
    if [[ ! -f "$file2" ]]; then
        log_error "Second benchmark file not found: $file2"
        return 1
    fi
    
    log_info "Comparing benchmarks: $file1 vs $file2"
    
    # Check if benchcmp is available
    if ! command -v benchcmp &> /dev/null; then
        log_warning "benchcmp not found, installing..."
        go install golang.org/x/tools/cmd/benchcmp@latest
    fi
    
    local compare_file="$OUTPUT_DIR/comparison_${TIMESTAMP}.txt"
    
    benchcmp "$file1" "$file2" | tee "$compare_file"
    
    log_success "Comparison saved to: $compare_file"
}

# Generate performance report
generate_report() {
    local report_file="$REPORTS_DIR/performance_report_${TIMESTAMP}.md"
    
    log_info "Generating performance report..."
    
    cat > "$report_file" << EOF
# SilentCast Performance Report

Generated: $(date)
Platform: $PLATFORM ($ARCH)

## Summary

This report contains performance benchmarks for SilentCast across different components.

## System Information

- Platform: $PLATFORM
- Architecture: $ARCH
- Go Version: $(go version)
- Build Tags: $BUILD_TAGS

## Benchmark Results

EOF

    # Include latest benchmark results
    for result_file in "$OUTPUT_DIR"/*_"$PLATFORM"_"$ARCH"_*.txt; do
        if [[ -f "$result_file" ]]; then
            local basename=$(basename "$result_file")
            local benchmark_type=$(echo "$basename" | cut -d'_' -f1)
            
            echo "### $benchmark_type Benchmarks" >> "$report_file"
            echo "" >> "$report_file"
            echo "\`\`\`" >> "$report_file"
            grep "^Benchmark" "$result_file" | head -20 >> "$report_file" || true
            echo "\`\`\`" >> "$report_file"
            echo "" >> "$report_file"
        fi
    done
    
    cat >> "$report_file" << EOF

## Performance Guidelines

### Target Performance Metrics

- Application startup: < 100ms
- Hotkey response time: < 10ms
- Memory usage (normal load): < 50MB
- Configuration reload: < 50ms
- Action execution overhead: < 5ms

### Optimization Areas

1. **Startup Performance**
   - Configuration loading optimization
   - Component initialization efficiency
   - Memory allocation reduction

2. **Runtime Performance**
   - Hotkey processing efficiency
   - Action execution optimization
   - Memory pool utilization

3. **Resource Management**
   - Garbage collection optimization
   - File descriptor management
   - Memory leak prevention

EOF

    log_success "Performance report generated: $report_file"
}

# Update performance baseline
update_baseline() {
    local baseline_dir="$PROJECT_ROOT/benchmark-baselines"
    mkdir -p "$baseline_dir"
    
    log_info "Updating performance baseline..."
    
    # Run a comprehensive benchmark for baseline
    local baseline_file="$baseline_dir/baseline_${PLATFORM}_${ARCH}.txt"
    
    cd "$APP_DIR"
    go test $BENCH_FLAGS -bench=. ./test/benchmarks/... > "$baseline_file"
    
    log_success "Baseline updated: $baseline_file"
}

# Check for performance regressions
check_regression() {
    local baseline_dir="$PROJECT_ROOT/benchmark-baselines"
    local baseline_file="$baseline_dir/baseline_${PLATFORM}_${ARCH}.txt"
    
    if [[ ! -f "$baseline_file" ]]; then
        log_warning "No baseline found for $PLATFORM/$ARCH. Creating new baseline..."
        update_baseline
        return 0
    fi
    
    log_info "Checking for performance regressions..."
    
    # Run current benchmarks
    local current_file="$OUTPUT_DIR/current_${PLATFORM}_${ARCH}_${TIMESTAMP}.txt"
    cd "$APP_DIR"
    go test $BENCH_FLAGS -bench=. ./test/benchmarks/... > "$current_file"
    
    # Compare with baseline
    local regression_file="$OUTPUT_DIR/regression_check_${TIMESTAMP}.txt"
    
    if command -v benchcmp &> /dev/null; then
        benchcmp "$baseline_file" "$current_file" > "$regression_file" 2>&1 || true
        
        # Check for significant regressions (>20% slower)
        if grep -q "+[2-9][0-9]\|+[0-9][0-9][0-9]" "$regression_file"; then
            log_error "Performance regression detected!"
            log_error "Check: $regression_file"
            cat "$regression_file"
            return 1
        else
            log_success "No significant performance regressions detected"
        fi
    else
        log_warning "benchcmp not available, manual comparison required"
        log_info "Baseline: $baseline_file"
        log_info "Current: $current_file"
    fi
}

# Main execution
main() {
    log_info "SilentCast Performance Benchmark Runner"
    log_info "Platform: $PLATFORM/$ARCH"
    log_info "Output Directory: $OUTPUT_DIR"
    
    # Collect system information for all runs
    collect_system_info
    
    case "$COMMAND" in
        compare)
            if [[ $# -lt 2 ]]; then
                log_error "Compare command requires two file arguments"
                exit 1
            fi
            compare_benchmarks "$1" "$2"
            ;;
        report)
            generate_report
            ;;
        baseline)
            update_baseline
            ;;
        regression)
            check_regression
            ;;
        *)
            # Run multiple times if specified
            for ((i=1; i<=RUNS; i++)); do
                if [[ "$RUNS" -gt 1 ]]; then
                    log_info "Benchmark run $i of $RUNS"
                fi
                run_benchmark "$COMMAND"
            done
            
            # Generate report if multiple runs
            if [[ "$RUNS" -gt 1 ]]; then
                generate_report
            fi
            ;;
    esac
    
    log_success "Benchmark execution completed"
}

# Execute main function with remaining arguments
main "$@"