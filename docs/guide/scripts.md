# Script Execution Guide

Unlock the full power of automation with SilentCast's script execution capabilities. From simple one-liners to complex workflows, this guide shows you how to run any command or script with a keyboard shortcut.

## Script Basics

Scripts in SilentCast are actions with `type: script` that execute shell commands. They're perfect for automation, system tasks, and developer workflows.

### Simple Script

```yaml
grimoire:
  hello_world:
    type: script
    command: "echo 'Hello, World!'"
    description: "My first script"
```

### Script Anatomy

```yaml
grimoire:
  complete_script:
    type: script
    command: "npm test"           # Command to execute
    working_dir: "~/project"      # Where to run it
    shell: "bash"                 # Which shell to use
    env:                          # Environment variables
      NODE_ENV: "test"
    show_output: true            # Show result notification
    keep_open: true              # Keep terminal open
    timeout: 300                 # Max execution time (seconds)
    description: "Run test suite"
```

## Shell Selection

### Default Shells by Platform

| Platform | Default Shell | Alternative Shells |
|----------|--------------|-------------------|
| macOS | `zsh` (Catalina+) or `bash` | `sh`, `fish` |
| Windows | `cmd` | `powershell`, `pwsh`, `bash` (WSL) |

### Specifying Shell

```yaml
grimoire:
  # Use specific shell
  bash_script:
    type: script
    shell: "bash"
    command: "echo $BASH_VERSION"
    
  powershell_script:
    type: script
    shell: "pwsh"  # PowerShell Core
    command: "Get-Host | Select-Object Version"
    
  python_script:
    type: script
    shell: "python3"
    command: "print('Python ' + sys.version)"
```

## Multi-line Scripts

### Using YAML Multi-line Syntax

```yaml
grimoire:
  multi_line_example:
    type: script
    command: |
      echo "Starting process..."
      echo "Step 1: Checking environment"
      
      if [ -d ".git" ]; then
        echo "✓ Git repository found"
      else
        echo "✗ Not a git repository"
        exit 1
      fi
      
      echo "Step 2: Running tests"
      npm test
      
      echo "Process complete!"
```

### Complex Workflows

```yaml
grimoire:
  deploy_workflow:
    type: script
    command: |
      #!/bin/bash
      set -e  # Exit on any error
      
      # Color output
      RED='\033[0;31m'
      GREEN='\033[0;32m'
      YELLOW='\033[1;33m'
      NC='\033[0m' # No Color
      
      echo -e "${YELLOW}🚀 Starting deployment...${NC}"
      
      # Pre-flight checks
      echo -e "${YELLOW}→ Running pre-flight checks${NC}"
      
      # Check Node.js version
      required_node="16"
      current_node=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
      if [ "$current_node" -lt "$required_node" ]; then
        echo -e "${RED}✗ Node.js $required_node+ required${NC}"
        exit 1
      fi
      
      # Check for uncommitted changes
      if ! git diff-index --quiet HEAD --; then
        echo -e "${RED}✗ Uncommitted changes detected${NC}"
        echo "Please commit or stash your changes first"
        exit 1
      fi
      
      # Run tests
      echo -e "${YELLOW}→ Running tests${NC}"
      if npm test; then
        echo -e "${GREEN}✓ Tests passed${NC}"
      else
        echo -e "${RED}✗ Tests failed${NC}"
        exit 1
      fi
      
      # Build
      echo -e "${YELLOW}→ Building application${NC}"
      npm run build
      echo -e "${GREEN}✓ Build complete${NC}"
      
      # Deploy
      echo -e "${YELLOW}→ Deploying to production${NC}"
      rsync -avz --delete dist/ user@prod-server:/var/www/app/
      
      # Post-deployment
      echo -e "${YELLOW}→ Running post-deployment tasks${NC}"
      ssh user@prod-server 'cd /var/www/app && npm run migrate'
      
      # Notify
      echo -e "${GREEN}✅ Deployment successful!${NC}"
      
      # Send notification (platform-specific)
      if command -v notify-send &> /dev/null; then
        notify-send "Deployment Complete" "Successfully deployed to production"
      elif command -v osascript &> /dev/null; then
        osascript -e 'display notification "Successfully deployed to production" with title "Deployment Complete"'
      fi
    keep_open: true
    timeout: 600  # 10 minutes
    description: "Full deployment workflow"
```

## Working with Output

### Show Output in Notifications

```yaml
grimoire:
  quick_status:
    type: script
    command: |
      echo "System: $(uname -s)"
      echo "Uptime: $(uptime -p)"
      echo "Memory: $(free -h | grep Mem | awk '{print $3 "/" $2}')"
    show_output: true  # Display in notification
```

### Keep Terminal Open

```yaml
grimoire:
  monitoring_script:
    type: script
    command: "tail -f /var/log/application.log"
    keep_open: true  # Terminal stays open
```

### Capture and Process Output

```yaml
grimoire:
  process_output:
    type: script
    command: |
      # Capture command output
      OUTPUT=$(curl -s https://api.github.com/user/repos | jq length)
      
      # Process the output
      if [ "$OUTPUT" -gt 50 ]; then
        echo "You have $OUTPUT repositories! Time to clean up?"
      else
        echo "You have $OUTPUT repositories."
      fi
    show_output: true
```

## Environment Variables

### Using System Environment

```yaml
grimoire:
  use_system_env:
    type: script
    command: |
      echo "Home: $HOME"
      echo "User: $USER"
      echo "Path: $PATH"
      echo "Editor: ${EDITOR:-not set}"
```

### Setting Custom Environment

```yaml
grimoire:
  custom_env:
    type: script
    env:
      API_KEY: "secret-key-123"
      API_URL: "https://api.example.com"
      DEBUG: "true"
    command: |
      echo "Calling API at $API_URL"
      curl -H "Authorization: Bearer $API_KEY" "$API_URL/status"
```

### Environment Variable Expansion

```yaml
grimoire:
  # Variables are expanded before execution
  open_project:
    type: script
    command: "cd ${PROJECT_DIR:-$HOME/projects} && code ."
    
  # Mix system and custom variables
  build_with_env:
    type: script
    env:
      BUILD_ENV: "production"
    command: |
      echo "Building for $BUILD_ENV on $HOSTNAME"
      npm run build:$BUILD_ENV
```

## Error Handling

### Basic Error Handling

```yaml
grimoire:
  safe_script:
    type: script
    command: |
      set -e  # Exit on error
      
      # Check prerequisites
      command -v git >/dev/null 2>&1 || { echo "Git not installed"; exit 1; }
      command -v npm >/dev/null 2>&1 || { echo "NPM not installed"; exit 1; }
      
      # Run commands
      git pull
      npm install
      npm test
```

### Advanced Error Handling

```yaml
grimoire:
  robust_script:
    type: script
    command: |
      #!/bin/bash
      
      # Error handler
      handle_error() {
        local exit_code=$1
        local line_number=$2
        echo "Error on line $line_number: Command exited with status $exit_code"
        
        # Cleanup on error
        rm -f /tmp/script.lock
        
        # Notify user
        notify-send "Script Failed" "Check the logs for details"
        
        exit $exit_code
      }
      
      # Set error trap
      trap 'handle_error $? $LINENO' ERR
      
      # Create lock file
      touch /tmp/script.lock
      
      # Your script logic here
      echo "Processing..."
      # Simulate potential failure
      [ -f "required-file.txt" ] || exit 1
      
      # Cleanup on success
      rm -f /tmp/script.lock
      echo "Success!"
    show_output: true
```

## Script Timeouts

### Basic Timeout Configuration

Set a maximum execution time for scripts to prevent hanging processes:

```yaml
grimoire:
  quick_task:
    type: script
    command: "echo 'Quick task'"
    timeout: 10  # 10 seconds timeout
    
  long_running_task:
    type: script
    command: "./build.sh && ./test.sh"
    timeout: 300  # 5 minutes timeout
    
  no_timeout_task:
    type: script
    command: "tail -f /var/log/app.log"
    timeout: 0  # No timeout (default)
```

### Timeout Behavior

When a script exceeds its timeout:
1. The process receives SIGTERM (or equivalent on Windows) for graceful shutdown
2. A configurable grace period allows the process to clean up
3. If the process doesn't exit within the grace period, SIGKILL is sent
4. User receives a timeout notification
5. Exit code indicates timeout occurred

```yaml
grimoire:
  test_with_timeout:
    type: script
    command: "npm test"
    timeout: 60  # 1 minute timeout
    show_output: true  # Will show partial output before timeout
    description: "Run tests with timeout"
```

### Timeout with Different Execution Modes

```yaml
grimoire:
  # Timeout applies to script execution
  background_with_timeout:
    type: script
    command: "sleep 30 && echo 'Done'"
    timeout: 10  # Will timeout after 10 seconds
    
  # Terminal launch timeout only
  terminal_with_timeout:
    type: script
    command: "npm run dev"
    terminal: true
    timeout: 5  # Only applies to terminal launch, not the script
    
  # Keep open with timeout
  monitor_with_timeout:
    type: script
    command: "watch -n 1 df -h"
    keep_open: true
    timeout: 5  # Only applies to terminal launch
```

### Best Practices for Timeouts

1. **Set reasonable timeouts**: Consider worst-case execution time
2. **Use timeouts for CI/CD**: Prevent hanging builds
3. **No timeout for monitors**: Long-running processes like logs
4. **Test timeout values**: Ensure they work in production
5. **Handle signals gracefully**: Implement SIGTERM handlers for clean shutdown
6. **Configure grace periods**: Allow enough time for cleanup operations

## Graceful Shutdown

SilentCast implements graceful process termination when scripts timeout, giving processes a chance to clean up before being forcefully terminated.

### How Graceful Shutdown Works

1. **Timeout Exceeded**: When a script runs longer than its configured timeout
2. **SIGTERM Sent**: The process receives a termination signal (SIGTERM on Unix, taskkill on Windows)
3. **Grace Period**: The process has a configurable grace period to clean up and exit
4. **Force Termination**: If the process doesn't exit within the grace period, SIGKILL is sent

### Configuring Grace Period

```yaml
grimoire:
  database_backup:
    type: script
    command: "pg_dump mydb > backup.sql"
    timeout: 30              # 30 seconds to complete
    grace_period: 10         # 10 seconds for graceful shutdown
    show_output: true
    description: "Database backup with graceful shutdown"
    
  quick_task:
    type: script
    command: "curl -X POST https://api.example.com/webhook"
    timeout: 5               # 5 seconds timeout
    grace_period: 2          # Only 2 seconds grace (quick cleanup)
    
  default_grace:
    type: script
    command: "./long-process.sh"
    timeout: 60              # 1 minute timeout
    grace_period: 0          # Use default 5 seconds
```

### Signal Handling Examples

#### Bash/Shell Script
```bash
#!/bin/bash
# Handle SIGTERM for graceful shutdown
cleanup() {
    echo "Received shutdown signal, cleaning up..."
    # Save state, close connections, etc.
    rm -f /tmp/process.lock
    echo "Cleanup complete"
    exit 0
}

trap cleanup SIGTERM SIGINT

# Main process
echo "Process started"
while true; do
    # Do work
    sleep 1
done
```

#### Python Script
```python
import signal
import sys
import time

def signal_handler(signum, frame):
    print("Received shutdown signal, cleaning up...")
    # Perform cleanup
    save_state()
    close_connections()
    print("Cleanup complete")
    sys.exit(0)

# Register signal handlers
signal.signal(signal.SIGTERM, signal_handler)
signal.signal(signal.SIGINT, signal_handler)

# Main process
print("Process started")
while True:
    # Do work
    time.sleep(1)
```

#### Node.js Script
```javascript
// Handle graceful shutdown
process.on('SIGTERM', () => {
    console.log('Received SIGTERM, shutting down gracefully');
    server.close(() => {
        console.log('HTTP server closed');
        database.close(() => {
            console.log('Database connection closed');
            process.exit(0);
        });
    });
});

// Main process
console.log('Server started');
```

### Platform Differences

#### Unix/Linux/macOS
- Uses SIGTERM for graceful shutdown
- Falls back to SIGKILL after grace period
- Process groups are used to terminate child processes

#### Windows
- Uses `taskkill` command for graceful termination
- Falls back to `taskkill /F` for force termination
- `/T` flag terminates process trees

### Best Practices

1. **Always implement signal handlers** in long-running scripts
2. **Set appropriate grace periods** based on cleanup requirements
3. **Log shutdown progress** to help debug timeout issues
4. **Test signal handling** independently from SilentCast
5. **Avoid blocking operations** in signal handlers
6. **Save state frequently** to minimize data loss

```yaml
grimoire:
  deployment:
    type: script
    command: |
      echo "Starting deployment..."
      
      # Each step has its own timeout
      timeout 30 npm test || exit 1
      timeout 60 npm run build || exit 1
      timeout 120 ./deploy.sh || exit 1
      
      echo "Deployment complete!"
    timeout: 240  # Overall 4-minute timeout
    show_output: true
    description: "Deploy with cascading timeouts"
```

## Interactive Scripts

### User Input

```yaml
grimoire:
  interactive_menu:
    type: script
    command: |
      echo "Choose an option:"
      echo "1) Development"
      echo "2) Staging"
      echo "3) Production"
      
      read -p "Enter choice (1-3): " choice
      
      case $choice in
        1) ENV="development" ;;
        2) ENV="staging" ;;
        3) ENV="production" ;;
        *) echo "Invalid choice"; exit 1 ;;
      esac
      
      echo "Deploying to $ENV..."
      ./deploy.sh --env=$ENV
    keep_open: true
```

### Using External Tools

```yaml
grimoire:
  # Using 'gum' for beautiful CLI prompts
  pretty_input:
    type: script
    command: |
      # Install: https://github.com/charmbracelet/gum
      NAME=$(gum input --placeholder "Project name")
      TYPE=$(gum choose "node" "python" "go" "rust")
      
      echo "Creating $TYPE project: $NAME"
      mkdir -p "$NAME"
      cd "$NAME"
      
      case $TYPE in
        node) npm init -y ;;
        python) python -m venv venv ;;
        go) go mod init "$NAME" ;;
        rust) cargo init ;;
      esac
    keep_open: true
    
  # Using 'fzf' for fuzzy finding
  fuzzy_find:
    type: script
    command: |
      # Find and open file in editor
      FILE=$(find . -type f -name "*.js" | fzf)
      [ -n "$FILE" ] && code "$FILE"
```

## Script Examples by Use Case

### Git Workflows

```yaml
grimoire:
  git_quick_commit:
    type: script
    command: |
      # Stage all changes
      git add -A
      
      # Generate commit message
      CHANGES=$(git diff --cached --stat | tail -1)
      DEFAULT_MSG="Update: $CHANGES"
      
      # Get commit message
      echo "Enter commit message (or press Enter for default):"
      echo "Default: $DEFAULT_MSG"
      read -r MESSAGE
      
      # Use default if empty
      MESSAGE="${MESSAGE:-$DEFAULT_MSG}"
      
      # Commit and push
      git commit -m "$MESSAGE"
      git push
    keep_open: true
    
  git_cleanup_branches:
    type: script
    command: |
      echo "Cleaning up merged branches..."
      
      # Fetch latest
      git fetch --prune
      
      # Delete merged local branches
      git branch --merged main | grep -v "main\|master" | xargs -n 1 git branch -d
      
      # Show remaining branches
      echo -e "\nRemaining branches:"
      git branch -a
    show_output: true
```

### Docker Operations

```yaml
grimoire:
  docker_cleanup:
    type: script
    command: |
      echo "🧹 Cleaning up Docker resources..."
      
      # Remove stopped containers
      echo "→ Removing stopped containers"
      docker container prune -f
      
      # Remove unused images
      echo "→ Removing unused images"
      docker image prune -a -f
      
      # Remove unused volumes
      echo "→ Removing unused volumes"
      docker volume prune -f
      
      # Remove unused networks
      echo "→ Removing unused networks"
      docker network prune -f
      
      # Show disk usage
      echo -e "\n📊 Disk usage after cleanup:"
      docker system df
    show_output: true
    
  docker_logs_live:
    type: script
    command: |
      # List running containers
      echo "Select container:"
      CONTAINER=$(docker ps --format "table {{.Names}}" | tail -n +2 | fzf)
      
      # Follow logs
      [ -n "$CONTAINER" ] && docker logs -f --tail=100 "$CONTAINER"
    keep_open: true
```

### Development Tasks

```yaml
grimoire:
  run_tests_with_coverage:
    type: script
    command: |
      # Clear previous coverage
      rm -rf coverage/
      
      # Run tests with coverage
      npm test -- --coverage --watchAll=false
      
      # Generate HTML report
      npx nyc report --reporter=html
      
      # Open coverage report
      if command -v xdg-open &> /dev/null; then
        xdg-open coverage/index.html
      elif command -v open &> /dev/null; then
        open coverage/index.html
      fi
    working_dir: "${PROJECT_ROOT}"
    
  database_backup:
    type: script
    command: |
      # Set variables
      DB_NAME="myapp"
      BACKUP_DIR="$HOME/backups"
      TIMESTAMP=$(date +%Y%m%d_%H%M%S)
      BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.sql"
      
      # Create backup directory
      mkdir -p "$BACKUP_DIR"
      
      # Perform backup
      echo "Backing up database: $DB_NAME"
      pg_dump $DB_NAME > "$BACKUP_FILE"
      
      # Compress
      gzip "$BACKUP_FILE"
      
      # Verify
      if [ -f "${BACKUP_FILE}.gz" ]; then
        SIZE=$(du -h "${BACKUP_FILE}.gz" | cut -f1)
        echo "✅ Backup complete: ${BACKUP_FILE}.gz ($SIZE)"
        
        # Clean old backups (keep last 5)
        ls -t "$BACKUP_DIR"/${DB_NAME}_*.sql.gz | tail -n +6 | xargs rm -f
      else
        echo "❌ Backup failed!"
        exit 1
      fi
    show_output: true
```

### System Administration

```yaml
grimoire:
  system_health_check:
    type: script
    command: |
      echo "=== System Health Check ==="
      echo "Date: $(date)"
      echo "Hostname: $(hostname)"
      echo ""
      
      # CPU Usage
      echo "🖥️  CPU Usage:"
      top -bn1 | grep "Cpu(s)" | awk '{print "   " $2 " user, " $4 " system"}'
      
      # Memory Usage
      echo -e "\n💾 Memory Usage:"
      free -h | grep Mem | awk '{print "   Used: " $3 " / " $2 " (" int($3/$2 * 100) "%)"}'
      
      # Disk Usage
      echo -e "\n💿 Disk Usage:"
      df -h / | tail -1 | awk '{print "   Used: " $3 " / " $2 " (" $5 ")"}'
      
      # Network
      echo -e "\n🌐 Network:"
      ip -4 addr show | grep inet | grep -v 127.0.0.1 | awk '{print "   " $2}'
      
      # Services
      echo -e "\n🔧 Key Services:"
      for service in nginx postgresql redis; do
        if systemctl is-active --quiet $service; then
          echo "   ✓ $service is running"
        else
          echo "   ✗ $service is not running"
        fi
      done
    show_output: true
    
  log_analysis:
    type: script
    command: |
      # Analyze nginx access logs
      echo "=== Top 10 IP addresses ==="
      awk '{print $1}' /var/log/nginx/access.log | sort | uniq -c | sort -rn | head -10
      
      echo -e "\n=== Top 10 requested URLs ==="
      awk '{print $7}' /var/log/nginx/access.log | sort | uniq -c | sort -rn | head -10
      
      echo -e "\n=== HTTP Status Codes ==="
      awk '{print $9}' /var/log/nginx/access.log | sort | uniq -c | sort -rn
      
      echo -e "\n=== Errors in last hour ==="
      grep "$(date -d '1 hour ago' '+%d/%b/%Y:%H')" /var/log/nginx/error.log || echo "No errors found"
    show_output: true
```

## Best Practices

### 1. Use Set Options for Safety

```yaml
grimoire:
  safe_script:
    type: script
    command: |
      #!/bin/bash
      set -euo pipefail  # Exit on error, undefined vars, pipe failures
      
      # Your script here
```

### 2. Add Progress Indicators

```yaml
grimoire:
  long_task:
    type: script
    command: |
      echo "Processing files..."
      
      total=$(find . -name "*.js" | wc -l)
      count=0
      
      find . -name "*.js" | while read file; do
        count=$((count + 1))
        echo "[$count/$total] Processing: $file"
        # Process file
      done
      
      echo "✅ Complete!"
```

### 3. Provide Clear Feedback

```yaml
grimoire:
  clear_feedback:
    type: script
    command: |
      echo "🚀 Starting backup process..."
      
      if perform_backup; then
        echo "✅ Backup completed successfully!"
        notify-send "Backup Complete" "All files backed up safely"
      else
        echo "❌ Backup failed! Check logs for details."
        notify-send "Backup Failed" "Please check the logs"
        exit 1
      fi
```

### 4. Make Scripts Idempotent

```yaml
grimoire:
  idempotent_setup:
    type: script
    command: |
      # Create directory only if it doesn't exist
      [ -d "$HOME/projects" ] || mkdir -p "$HOME/projects"
      
      # Install package only if not installed
      command -v node >/dev/null 2>&1 || {
        echo "Installing Node.js..."
        # Installation commands
      }
      
      echo "Environment ready!"
```

## Debugging Scripts

### Debug Mode

```yaml
grimoire:
  debug_script:
    type: script
    command: |
      #!/bin/bash -x  # Enable debug mode
      
      # Shows each command as it's executed
      echo "Starting debug script"
      VAR="test"
      echo "Variable is: $VAR"
```

### Logging

```yaml
grimoire:
  logged_script:
    type: script
    command: |
      LOG_FILE="$HOME/.local/share/silentcast/script.log"
      
      # Function to log with timestamp
      log() {
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
      }
      
      log "Script started"
      log "Running process..."
      
      # Your script logic
      
      log "Script completed"
```

## Next Steps

- Learn about [Environment Variables](/guide/env-vars) in detail
- Explore [Platform-specific](/guide/platforms) script features
- Check the [CLI Reference](/api/) for script-related options
- Share your scripts in our [Community Repository](https://github.com/SphereStacking/silentcast-scripts)