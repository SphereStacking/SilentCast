# Windows Guide

This guide covers Windows-specific features, configuration, and troubleshooting for SilentCast.

## Installation

### Download and Install
1. Download the Windows binary from [GitHub Releases](https://github.com/SphereStacking/silentcast/releases)
2. Extract to a folder (e.g., `C:\Program Files\SilentCast\`)
3. Add the folder to your system PATH (optional)
4. Run `silentcast.exe`

### Windows Terminal Integration
For the best experience, install [Windows Terminal](https://aka.ms/terminal) from the Microsoft Store.

## Windows-Specific Configuration

### Shell Options
SilentCast supports multiple shell environments on Windows:

```yaml
grimoire:
  # Command Prompt
  cmd_example:
    type: script
    command: "dir"
    shell: "cmd"
    show_output: true
    
  # PowerShell
  ps_example:
    type: script
    command: "Get-ChildItem"
    shell: "powershell"
    show_output: true
    
  # PowerShell Core (if installed)
  pwsh_example:
    type: script
    command: "Get-Process"
    shell: "pwsh"
    show_output: true
```

### Windows Terminal Integration
Open specific tabs or panes in Windows Terminal:

```yaml
grimoire:
  # Open new Windows Terminal window
  new_terminal:
    type: app
    command: "wt"
    description: "New Windows Terminal"
    
  # Open terminal in specific directory
  dev_terminal:
    type: app
    command: "wt"
    args: ["-d", "C:\\dev"]
    description: "Dev directory terminal"
    
  # Open multiple tabs
  multi_terminal:
    type: app
    command: "wt"
    args: ["new-tab", "PowerShell", ";", "new-tab", "cmd"]
    description: "Terminal with PowerShell and CMD tabs"
```

### File Explorer Integration
```yaml
grimoire:
  # Open current directory in Explorer
  explorer_here:
    type: script
    command: "explorer ."
    description: "Open Explorer here"
    
  # Open specific folder
  documents:
    type: script
    command: "explorer C:\\Users\\%USERNAME%\\Documents"
    description: "Open Documents folder"
```

## Windows-Specific Spells Examples

### Development Tools
```yaml
spells:
  # Visual Studio Code
  "e": "vscode"
  "e,c": "vscode_current"
  
  # Terminal operations
  "t": "terminal"
  "t,a": "admin_terminal"
  
  # Git operations
  "g,s": "git_status"
  "g,p": "git_push"

grimoire:
  vscode:
    type: app
    command: "code"
    description: "Open VS Code"
    
  vscode_current:
    type: app
    command: "code"
    args: ["."]
    description: "Open VS Code in current directory"
    
  terminal:
    type: app
    command: "wt"
    description: "Windows Terminal"
    
  admin_terminal:
    type: app
    command: "wt"
    admin: true  # Run as administrator
    description: "Admin Terminal"
    
  git_status:
    type: script
    command: "git status --short"
    show_output: true
    description: "Git status"
    
  git_push:
    type: script
    command: "git push"
    show_output: true
    terminal: true
    description: "Git push"
```

### System Operations
```yaml
grimoire:
  # Windows system tools
  task_manager:
    type: script
    command: "taskmgr"
    description: "Task Manager"
    
  system_info:
    type: script
    command: "systeminfo"
    show_output: true
    description: "System Information"
    
  # Network operations
  ip_config:
    type: script
    command: "ipconfig /all"
    show_output: true
    description: "IP Configuration"
    
  ping_google:
    type: script
    command: "ping google.com"
    show_output: true
    description: "Ping Google"
```

### File Operations
```yaml
grimoire:
  # Directory listing with details
  list_files:
    type: script
    command: "dir /a"
    show_output: true
    description: "List all files"
    
  # Find files
  find_file:
    type: script
    command: "where notepad.exe"
    show_output: true
    description: "Find notepad.exe"
    
  # Disk usage
  disk_usage:
    type: script
    command: "wmic logicaldisk get size,freespace,caption"
    show_output: true
    description: "Disk usage"
```

## Administrator Privileges

SilentCast supports running commands with elevated privileges:

```yaml
grimoire:
  # Admin command prompt
  admin_cmd:
    type: script
    command: "cmd"
    admin: true
    terminal: true
    keep_open: true
    description: "Admin Command Prompt"
    
  # Install software (requires admin)
  install_chocolatey:
    type: script
    command: "powershell -Command \"Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))\""
    admin: true
    show_output: true
    description: "Install Chocolatey"
```

## PowerShell Configuration

### Execution Policy
If PowerShell scripts aren't running, check the execution policy:

```powershell
# Check current policy
Get-ExecutionPolicy

# Set to allow local scripts (run as Administrator)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### PowerShell Profile Integration
Add SilentCast helpers to your PowerShell profile:

```powershell
# Add to $PROFILE
function sc-reload {
    Write-Host "Reloading SilentCast configuration..."
    # Signal SilentCast to reload (if feature exists)
}

function sc-test {
    param($spell)
    Write-Host "Testing spell: $spell"
    # Test spell execution
}
```

## Windows-Specific Troubleshooting

### Common Issues

#### Permission Errors
- Run SilentCast as Administrator for system operations
- Check User Account Control (UAC) settings
- Verify file/folder permissions

#### Antivirus Interference
- Add SilentCast to antivirus exclusions
- Some antivirus software may block global hotkey registration
- Check Windows Defender exclusions

#### Path Issues
- Use full paths for executables: `C:\Program Files\App\app.exe`
- Check PATH environment variable
- Use quotes for paths with spaces: `"C:\Program Files\App\app.exe"`

### Registry and System Integration

#### Windows Registry Access
```yaml
grimoire:
  registry_query:
    type: script
    command: "reg query HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion"
    show_output: true
    description: "Query Windows version from registry"
```

#### Windows Services
```yaml
grimoire:
  service_status:
    type: script
    command: "sc query Spooler"
    show_output: true
    description: "Check Print Spooler service"
    
  restart_service:
    type: script
    command: "net stop Spooler && net start Spooler"
    admin: true
    show_output: true
    description: "Restart Print Spooler"
```

## Performance Optimization

### Windows-Specific Optimizations
- Use `show_output: false` for background tasks
- Prefer native Windows commands over third-party tools
- Use Windows Terminal for better performance with console applications

### Resource Management
```yaml
grimoire:
  memory_usage:
    type: script
    command: "tasklist /fi \"imagename eq silentcast.exe\""
    show_output: true
    description: "Check SilentCast memory usage"
    
  cleanup_temp:
    type: script
    command: "del /q /s %TEMP%\\*"
    admin: true
    description: "Clean temporary files"
```

## Integration Examples

### Development Workflow
```yaml
spells:
  # Project management
  "p,o": "open_project"
  "p,b": "build_project"
  "p,t": "test_project"
  "p,d": "deploy_project"
  
  # Docker operations
  "d,u": "docker_up"
  "d,d": "docker_down"
  "d,l": "docker_logs"

grimoire:
  open_project:
    type: script
    command: "code C:\\dev\\myproject"
    description: "Open project in VS Code"
    
  build_project:
    type: script
    command: "dotnet build"
    working_dir: "C:\\dev\\myproject"
    show_output: true
    description: "Build .NET project"
    
  docker_up:
    type: script
    command: "docker-compose up -d"
    working_dir: "C:\\dev\\myproject"
    show_output: true
    description: "Start Docker containers"
```

### System Administration
```yaml
grimoire:
  # Event logs
  system_events:
    type: script
    command: "wevtutil qe System /c:10 /rd:true /f:text"
    show_output: true
    description: "Recent system events"
    
  # Process management
  kill_process:
    type: script
    command: "taskkill /f /im notepad.exe"
    show_output: true
    description: "Kill notepad processes"
    
  # Network diagnostics
  network_reset:
    type: script
    command: "ipconfig /release && ipconfig /renew"
    admin: true
    show_output: true
    description: "Reset network configuration"
```

## Best Practices for Windows

1. **Use specific shell types**: Specify `shell: "cmd"` or `shell: "powershell"` explicitly
2. **Handle spaces in paths**: Use quotes for file paths with spaces
3. **Test admin permissions**: Test admin commands in a safe environment first
4. **Use Windows Terminal**: For better Unicode and color support
5. **Environment variables**: Use Windows environment variables like `%USERPROFILE%`, `%TEMP%`
6. **Error handling**: Use `show_output: true` to see error messages

## Security Considerations

- Be cautious with `admin: true` - only use when necessary
- Validate user input in scripts
- Use specific paths instead of relying on PATH when possible
- Regular security updates for both SilentCast and system
- Consider using Windows Defender Application Guard for untrusted scripts