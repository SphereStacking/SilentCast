#!/bin/bash

# TDD Metrics Collection Script
# Measures and tracks TDD cycle performance

set -e

METRICS_DIR="metrics/tdd"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
CYCLE_LOG="$METRICS_DIR/cycles_$TIMESTAMP.json"

# Create metrics directory if it doesn't exist
mkdir -p "$METRICS_DIR"

# Initialize cycle tracking
echo "# TDD Cycle Metrics - $(date)" > "$METRICS_DIR/session_$TIMESTAMP.log"
echo "Starting TDD metrics collection session..." 

# Function to log cycle start
cycle_start() {
    CYCLE_START_TIME=$(date +%s)
    echo "🔴 RED phase started at $(date)" | tee -a "$METRICS_DIR/session_$TIMESTAMP.log"
    echo "{\"cycle_start\": \"$(date -Iseconds)\", \"phase\": \"red\", \"timestamp\": $CYCLE_START_TIME}" >> "$CYCLE_LOG"
}

# Function to log phase transition
phase_transition() {
    local phase=$1
    local current_time=$(date +%s)
    local duration=$((current_time - CYCLE_START_TIME))
    
    echo "🔄 Transitioning to $phase phase ($(format_duration $duration))" | tee -a "$METRICS_DIR/session_$TIMESTAMP.log"
    echo "{\"phase_transition\": \"$phase\", \"timestamp\": $current_time, \"duration_seconds\": $duration}" >> "$CYCLE_LOG"
}

# Function to log cycle completion
cycle_complete() {
    local current_time=$(date +%s)
    local total_duration=$((current_time - CYCLE_START_TIME))
    
    echo "✅ Cycle completed in $(format_duration $total_duration)" | tee -a "$METRICS_DIR/session_$TIMESTAMP.log"
    echo "{\"cycle_complete\": \"$(date -Iseconds)\", \"timestamp\": $current_time, \"total_duration_seconds\": $total_duration}" >> "$CYCLE_LOG"
    
    # Generate coverage report
    go test -tags "nogohook notray" -coverprofile="$METRICS_DIR/coverage_$TIMESTAMP.out" ./... > /dev/null 2>&1
    local coverage=$(go tool cover -func="$METRICS_DIR/coverage_$TIMESTAMP.out" | grep total | awk '{print $3}')
    echo "{\"coverage\": \"$coverage\", \"timestamp\": $current_time}" >> "$CYCLE_LOG"
    
    echo "📊 Coverage: $coverage" | tee -a "$METRICS_DIR/session_$TIMESTAMP.log"
}

# Function to format duration
format_duration() {
    local seconds=$1
    local minutes=$((seconds / 60))
    local remaining_seconds=$((seconds % 60))
    
    if [ $minutes -eq 0 ]; then
        echo "${seconds}s"
    else
        echo "${minutes}m ${remaining_seconds}s"
    fi
}

# Function to run tests and measure timing
timed_test() {
    local phase=$1
    local start_time=$(date +%s)
    
    case $phase in
        "red")
            echo "🔴 Running RED phase tests..."
            go test -tags "nogohook notray" -v -failfast ./... || true
            ;;
        "green")
            echo "🟢 Running GREEN phase tests..."
            go test -tags "nogohook notray" -v -race -cover ./...
            ;;
        "refactor")
            echo "🔵 Running REFACTOR phase tests..."
            go test -tags "nogohook notray" -v -race -cover ./...
            ;;
    esac
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    echo "{\"test_run\": \"$phase\", \"timestamp\": $end_time, \"duration_seconds\": $duration}" >> "$CYCLE_LOG"
    echo "Test run completed in $(format_duration $duration)"
}

# Function to generate daily report
generate_report() {
    local date_filter=${1:-$(date +%Y%m%d)}
    echo "📈 Generating TDD metrics report for $date_filter..."
    
    local report_file="$METRICS_DIR/report_$date_filter.md"
    
    cat > "$report_file" << EOF
# TDD Metrics Report - $date_filter

## Summary

Generated: $(date)

## Cycle Performance
EOF
    
    # Parse cycle logs for the day
    local cycle_files=$(find "$METRICS_DIR" -name "cycles_${date_filter}_*.json" 2>/dev/null || true)
    local total_cycles=0
    local total_time=0
    local under_10min=0
    
    if [ -n "$cycle_files" ]; then
        for file in $cycle_files; do
            while IFS= read -r line; do
                if echo "$line" | grep -q "cycle_complete"; then
                    local duration=$(echo "$line" | grep -o '"total_duration_seconds":[[:space:]]*[0-9]*' | grep -o '[0-9]*')
                    if [ -n "$duration" ]; then
                        total_cycles=$((total_cycles + 1))
                        total_time=$((total_time + duration))
                        if [ "$duration" -le 600 ]; then  # 10 minutes = 600 seconds
                            under_10min=$((under_10min + 1))
                        fi
                    fi
                fi
            done < "$file"
        done
    fi
    
    if [ $total_cycles -gt 0 ]; then
        local avg_time=$((total_time / total_cycles))
        local success_rate=$((under_10min * 100 / total_cycles))
        
        cat >> "$report_file" << EOF

- Total Cycles: $total_cycles
- Average Cycle Time: $(format_duration $avg_time)
- Cycles Under 10 Minutes: $under_10min/$total_cycles ($success_rate%)
- Total Development Time: $(format_duration $total_time)

## Recommendations

EOF
        
        if [ $success_rate -lt 80 ]; then
            echo "- ⚠️  Consider breaking down larger tasks into smaller units" >> "$report_file"
        fi
        
        if [ $avg_time -gt 600 ]; then
            echo "- ⚠️  Average cycle time exceeds 10 minutes - focus on smaller increments" >> "$report_file"
        fi
        
        if [ $success_rate -ge 80 ]; then
            echo "- ✅ Good TDD rhythm maintained!" >> "$report_file"
        fi
    else
        echo "No completed cycles found for $date_filter" >> "$report_file"
    fi
    
    echo "Report generated: $report_file"
}

# Main command handling
case "${1:-help}" in
    "start")
        cycle_start
        ;;
    "red")
        timed_test "red"
        ;;
    "green")
        phase_transition "green"
        timed_test "green"
        ;;
    "refactor")
        phase_transition "refactor"
        timed_test "refactor"
        ;;
    "complete")
        cycle_complete
        ;;
    "report")
        generate_report "$2"
        ;;
    "help"|*)
        echo "TDD Metrics Script"
        echo ""
        echo "Usage: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  start     Start a new TDD cycle"
        echo "  red       Run RED phase tests"
        echo "  green     Run GREEN phase tests" 
        echo "  refactor  Run REFACTOR phase tests"
        echo "  complete  Complete current cycle and generate metrics"
        echo "  report [date]  Generate report for date (YYYYMMDD, default: today)"
        echo ""
        echo "Example TDD session:"
        echo "  $0 start"
        echo "  # Write failing test"
        echo "  $0 red"
        echo "  # Write minimal implementation"
        echo "  $0 green"
        echo "  # Refactor code"
        echo "  $0 refactor"
        echo "  $0 complete"
        ;;
esac