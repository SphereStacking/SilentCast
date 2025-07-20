# Grimoire Guide

Learn how to define powerful grimoire entries in your spellbook. From launching applications to running complex scripts, SilentCast can execute virtually any task with a simple spell.

## Understanding Grimoire Entries

In SilentCast, **grimoire entries** are defined in the `grimoire` section of your configuration. Each grimoire entry has a type that determines how it's executed.

### Grimoire Entry Anatomy

```yaml
grimoire:
  grimoire_entry_name:      # Unique identifier
    type: app              # Entry type: app, script, or url
    command: "code"        # What to execute
    description: "..."     # Human-readable description
    # Optional parameters depending on type
```

## Grimoire Entry Types

### Type: `app` - Application Launcher

Launches applications with optional arguments and environment settings.

#### Basic Application Launch

```yaml
grimoire:
  editor:
    type: app
    command: "code"
    description: "Open Visual Studio Code"
```

#### With Arguments

```yaml
grimoire:
  editor_with_folder:
    type: app
    command: "code"
    args: ["~/projects", "--new-window", "--goto", "src/main.js:10:5"]
    description: "Open VS Code with project folder at specific line"
```

#### With Working Directory

```yaml
grimoire:
  terminal_in_project:
    type: app
    command: "wt"
    working_dir: "~/projects/my-app"
    description: "Open terminal in project directory"
```

#### With Environment Variables

```yaml
grimoire:
  dev_server:
    type: app
    command: "npm"
    args: ["run", "dev"]
    working_dir: "${PROJECT_ROOT}"
    env:
      NODE_ENV: "development"
      PORT: "3000"
      DEBUG: "app:*"
    description: "Start development server"
```

### Type: `script` - Command Execution

Executes shell commands and scripts with full control over the execution environment.

#### Simple Commands

```yaml
grimoire:
  git_status:
    type: script
    command: "git status"
    description: "Show git repository status"
```

#### Multi-line Scripts

```yaml
grimoire:
  deploy:
    type: script
    command: |
      echo "Starting deployment..."
      npm test || exit 1
      npm run build
      rsync -avz dist/ user@server:/var/www/
      echo "Deployment complete!"
    description: "Test, build, and deploy application"
```

#### Show Output

```yaml
grimoire:
  system_info:
    type: script
    command: "neofetch"
    show_output: true      # Display output in notification
    description: "Show system information"
```

#### Keep Terminal Open

```yaml
grimoire:
  log_monitor:
    type: script
    command: "tail -f /var/log/app.log | grep ERROR"
    keep_open: true        # Keep terminal window open
    description: "Monitor application errors"
```

#### Custom Shell

```yaml
grimoire:
  python_script:
    type: script
    command: "print('Hello from Python!')"
    shell: "python3"       # Use Python as interpreter
    description: "Run Python code directly"
    
  powershell_script:
    type: script
    command: "Get-Process | Sort-Object CPU -Descending | Select-Object -First 10"
    shell: "pwsh"          # Use PowerShell Core
    description: "Show top CPU-consuming processes"
```

### Type: `url` - URL Opener

Opens URLs in the default browser or specified application.

#### Basic URL

```yaml
grimoire:
  documentation:
    type: url
    command: "https://silentcast.dev/docs"
    description: "Open SilentCast documentation"
```

#### With Specific Browser

```yaml
grimoire:
  github_in_firefox:
    type: url
    command: "https://github.com"
    browser: "firefox"     # Open in specific browser
    description: "Open GitHub in Firefox"
```

#### Local Files as URLs

```yaml
grimoire:
  local_docs:
    type: url
    command: "file:///home/user/docs/index.html"
    description: "Open local documentation"
```

## Advanced Action Patterns

### Dynamic Commands

Use environment variables and shell expansion:

```yaml
grimoire:
  open_current_project:
    type: app
    command: "${EDITOR:-code}"
    args: ["${PWD}"]
    description: "Open current directory in preferred editor"
    
  backup_with_timestamp:
    type: script
    command: |
      TIMESTAMP=$(date +%Y%m%d_%H%M%S)
      tar -czf "backup_${TIMESTAMP}.tar.gz" src/
      echo "Backup created: backup_${TIMESTAMP}.tar.gz"
    show_output: true
    description: "Create timestamped backup"
```

### Conditional Execution

```yaml
grimoire:
  smart_git_pull:
    type: script
    command: |
      # Check if we're in a git repository
      if [ -d .git ]; then
        # Check for uncommitted changes
        if [ -z "$(git status --porcelain)" ]; then
          git pull
        else
          echo "⚠️  Uncommitted changes detected!"
          echo "Please commit or stash changes first."
        fi
      else
        echo "❌ Not a git repository"
      fi
    show_output: true
    description: "Pull changes only if working directory is clean"
```

### Interactive Scripts

```yaml
grimoire:
  interactive_deploy:
    type: script
    command: |
      echo "Deploy to which environment?"
      echo "1) Development"
      echo "2) Staging"
      echo "3) Production"
      read -p "Choice (1-3): " choice
      
      case $choice in
        1) TARGET="dev" ;;
        2) TARGET="staging" ;;
        3) TARGET="prod" ;;
        *) echo "Invalid choice"; exit 1 ;;
      esac
      
      ./deploy.sh --target=$TARGET
    keep_open: true
    description: "Interactive deployment script"
```

### Chained Actions

Execute multiple related commands:

```yaml
grimoire:
  full_test_suite:
    type: script
    command: |
      echo "=== Running Full Test Suite ==="
      
      # Unit tests
      echo "→ Unit tests..."
      npm run test:unit || exit 1
      
      # Integration tests
      echo "→ Integration tests..."
      npm run test:integration || exit 1
      
      # E2E tests
      echo "→ E2E tests..."
      npm run test:e2e || exit 1
      
      # Coverage report
      echo "→ Generating coverage report..."
      npm run coverage
      
      echo "✅ All tests passed!"
    keep_open: true
    description: "Run complete test suite with coverage"
```

## Platform-Specific Actions

### macOS Actions

```yaml
grimoire:
  # Use AppleScript
  toggle_dark_mode:
    type: script
    command: |
      osascript -e 'tell app "System Events" to tell appearance preferences to set dark mode to not dark mode'
    description: "Toggle macOS dark mode"
    
  # Open in specific app
  preview_image:
    type: script
    command: "open -a Preview"
    args: ["${FILE}"]
    description: "Open image in Preview"
    
  # System controls
  sleep_display:
    type: script
    command: "pmset displaysleepnow"
    description: "Put display to sleep"
```

### Windows Actions

```yaml
grimoire:
  # Windows-specific commands
  open_hosts_file:
    type: app
    command: "notepad"
    args: ["C:\\Windows\\System32\\drivers\\etc\\hosts"]
    admin: true  # Run as administrator
    description: "Edit hosts file"
    
  # PowerShell automation
  cleanup_temp:
    type: script
    shell: "powershell"
    command: |
      Remove-Item -Path "$env:TEMP\*" -Recurse -Force -ErrorAction SilentlyContinue
      Write-Host "Temp files cleaned!"
    show_output: true
    description: "Clean temporary files"
    
  # System utilities
  flush_dns:
    type: script
    command: "ipconfig /flushdns"
    admin: true
    description: "Flush DNS cache"
```

## Action Parameters Reference

### Common Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `type` | string | Action type: `app`, `script`, or `url` | Required |
| `command` | string | Command to execute | Required |
| `description` | string | Human-readable description | Optional |
| `working_dir` | string | Working directory | Current directory |
| `env` | object | Environment variables | Inherited |

### App-Specific Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `args` | array | Command line arguments | `[]` |
| `admin` | boolean | Run with elevated privileges | `false` |
| `terminal` | boolean | Run in terminal window | `false` |

### Script-Specific Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `shell` | string | Shell interpreter | System default |
| `show_output` | boolean | Show output in notification | `false` |
| `keep_open` | boolean | Keep terminal open after execution | `false` |
| `timeout` | integer | Maximum execution time (seconds) | `0` (no timeout) |

### URL-Specific Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `browser` | string | Specific browser to use | System default |

## Best Practices

### 1. Use Descriptive Names

```yaml
# ❌ Poor naming
grimoire:
  action1:
    type: app
    command: "code"
    
# ✅ Good naming
grimoire:
  open_vscode_editor:
    type: app
    command: "code"
    description: "Open Visual Studio Code editor"
```

### 2. Add Meaningful Descriptions

```yaml
grimoire:
  deploy_production:
    type: script
    command: "./deploy.sh --prod"
    description: |
      Deploy to production environment
      - Runs tests
      - Builds application
      - Deploys to AWS
      - Sends notification to Slack
```

### 3. Handle Errors Gracefully

```yaml
grimoire:
  safe_deploy:
    type: script
    command: |
      set -e  # Exit on error
      
      # Check prerequisites
      command -v node || { echo "Node.js not installed"; exit 1; }
      [ -f package.json ] || { echo "Not in a Node project"; exit 1; }
      
      # Run deployment
      npm test && npm run build && npm run deploy
    show_output: true
```

### 4. Use Environment Variables

```yaml
grimoire:
  flexible_editor:
    type: app
    command: "${EDITOR:-code}"  # Use $EDITOR or default to VS Code
    args: ["${PROJECT_DIR:-$HOME/projects}"]
```

### 5. Provide Feedback

```yaml
grimoire:
  long_running_task:
    type: script
    command: |
      echo "🚀 Starting process..."
      # Long running command
      ./process_data.sh
      echo "✅ Process complete!"
      
      # Send notification
      notify-send "SilentCast" "Data processing complete!"
    keep_open: true
    show_output: true
```

## Examples by Use Case

### Development Workflow

```yaml
grimoire:
  # Start development environment
  dev_start:
    type: script
    command: |
      docker-compose up -d
      npm run dev
    keep_open: true
    description: "Start full development stack"
    
  # Run tests with coverage
  test_coverage:
    type: script
    command: "npm run test -- --coverage --watch"
    keep_open: true
    description: "Run tests in watch mode with coverage"
    
  # Quick commit
  quick_commit:
    type: script
    command: |
      git add -A
      git commit -m "$(gum input --placeholder 'Commit message')"
      git push
    description: "Stage, commit, and push changes"
```

### System Administration

```yaml
grimoire:
  # SSH to servers
  ssh_production:
    type: app
    command: "ssh"
    args: ["user@prod.example.com", "-p", "22"]
    terminal: true
    description: "SSH to production server"
    
  # Monitor logs
  monitor_nginx:
    type: script
    command: "sudo tail -f /var/log/nginx/error.log | grep -v 'info'"
    keep_open: true
    description: "Monitor nginx error logs"
    
  # System health check
  health_check:
    type: script
    command: |
      echo "=== System Health Check ==="
      echo "CPU: $(top -bn1 | grep 'Cpu(s)' | awk '{print $2}')%"
      echo "Memory: $(free -h | grep Mem | awk '{print $3 "/" $2}')"
      echo "Disk: $(df -h / | tail -1 | awk '{print $5}')"
      echo "Uptime: $(uptime -p)"
    show_output: true
    description: "Quick system health check"
```

## Troubleshooting Actions

### Action Not Executing

1. **Check the command path**:
   ```yaml
   grimoire:
     my_app:
       type: app
       command: "/full/path/to/app"  # Use full path if not in PATH
   ```

2. **Verify working directory**:
   ```yaml
   grimoire:
     my_script:
       type: script
       command: "./script.sh"
       working_dir: "/path/to/scripts"  # Ensure script exists here
   ```

3. **Check permissions**:
   ```yaml
   grimoire:
     admin_task:
       type: script
       command: "sudo systemctl restart nginx"
       admin: true  # May need admin flag
   ```

### Debugging Actions

```yaml
grimoire:
  debug_action:
    type: script
    command: |
      echo "Current directory: $(pwd)"
      echo "User: $(whoami)"
      echo "Path: $PATH"
      echo "Environment: $(env | sort)"
    show_output: true
    keep_open: true
    description: "Debug environment information"
```

## Next Steps

- Explore [Script Execution](/guide/scripts) for advanced scripting
- Learn about [Environment Variables](/guide/env-vars)
- Check [Platform-specific features](/guide/platforms)
- See [real-world examples](https://github.com/SphereStacking/silentcast-examples)