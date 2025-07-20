# Custom Shell Support

SilentCast supports running scripts with custom shells, allowing you to use your preferred shell or interpreter for script execution.

## Overview

The custom shell feature:
- Automatically detects available shells on your system
- Supports common shells (bash, zsh, fish, PowerShell, cmd)
- Supports interpreters (Python, Node.js, Ruby, Perl)
- Detects shell from script shebang
- Provides fallback to system default shell
- Platform-specific shell detection

## Configuration

### Specifying a Shell

Use the `shell` field in your grimoire configuration:

```yaml
grimoire:
  python_script:
    type: script
    command: "print('Hello from Python')"
    shell: "python"
    description: "Run Python code"

  powershell_script:
    type: script
    command: "Write-Host 'Hello from PowerShell'"
    shell: "powershell"
    description: "Run PowerShell command"
```

### Available Shells

Common shells you can specify:

| Shell | Names | Platform |
|-------|-------|----------|
| Bash | `bash`, `sh` | Unix/Linux/macOS |
| Zsh | `zsh` | Unix/Linux/macOS |
| Fish | `fish` | Unix/Linux/macOS |
| Dash | `dash` | Unix/Linux |
| PowerShell | `powershell`, `pwsh` | Windows/Cross-platform |
| Command Prompt | `cmd` | Windows |
| Python | `python`, `python3` | Cross-platform |
| Node.js | `node`, `nodejs` | Cross-platform |
| Ruby | `ruby` | Cross-platform |
| Perl | `perl` | Cross-platform |

## Shell Detection

### Automatic Detection

SilentCast automatically detects shells by:
1. Checking system paths
2. Reading `/etc/shells` (Unix/Linux)
3. Checking Windows registry
4. Scanning common installation directories

### List Available Shells

To see what shells are available:

```yaml
grimoire:
  list_shells:
    type: script
    command: |
      echo "Available shells:"
      # This would show detected shells
    description: "List available shells"
    terminal: true
```

### Shebang Detection

Scripts with shebang are automatically executed with the correct interpreter:

```yaml
grimoire:
  shebang_script:
    type: script
    command: |
      #!/usr/bin/env python3
      import sys
      print(f"Python {sys.version}")
    description: "Script with shebang"
    # No shell specified - will use Python from shebang
```

## Examples

### Multi-Language Scripts

```yaml
spells:
  p,y: run_python
  n,j: run_node
  r,b: run_ruby

grimoire:
  run_python:
    type: script
    command: |
      import datetime
      print(f"Current time: {datetime.datetime.now()}")
    shell: "python"
    description: "Python timestamp"
    show_output: true

  run_node:
    type: script
    command: |
      console.log(`Node.js ${process.version}`);
      console.log(`Current time: ${new Date()}`);
    shell: "node"
    description: "Node.js info"
    show_output: true

  run_ruby:
    type: script
    command: |
      puts "Ruby #{RUBY_VERSION}"
      puts "Current time: #{Time.now}"
    shell: "ruby"
    description: "Ruby info"
    show_output: true
```

### Platform-Specific Shells

```yaml
grimoire:
  # Windows-specific
  windows_info:
    type: script
    command: |
      Get-ComputerInfo | Select-Object CsName, WindowsVersion, WindowsBuildLabEx
    shell: "powershell"
    description: "Windows system info"
    show_output: true

  # Unix/Linux-specific
  unix_info:
    type: script
    command: |
      uname -a
      echo "Shell: $SHELL"
      echo "User: $USER"
    shell: "bash"
    description: "Unix system info"
    show_output: true
```

### Shell Features

```yaml
grimoire:
  # Bash arrays and loops
  bash_features:
    type: script
    command: |
      fruits=("apple" "banana" "orange")
      for fruit in "${fruits[@]}"; do
        echo "I like $fruit"
      done
    shell: "bash"
    description: "Bash array example"
    show_output: true

  # PowerShell objects
  powershell_features:
    type: script
    command: |
      Get-Process | 
        Where-Object {$_.CPU -gt 10} |
        Select-Object Name, CPU |
        Sort-Object CPU -Descending |
        Select-Object -First 5
    shell: "powershell"
    description: "Top 5 CPU processes"
    show_output: true

  # Python with modules
  python_features:
    type: script
    command: |
      import json
      import requests
      
      # Note: This is an example, adjust URL as needed
      data = {"status": "running", "timestamp": "now"}
      print(json.dumps(data, indent=2))
    shell: "python"
    description: "Python JSON example"
    show_output: true
```

### Interactive Shells

```yaml
grimoire:
  # Python REPL
  python_repl:
    type: script
    command: "python"
    shell: "sh"  # Launch through shell
    terminal: true
    description: "Python interactive shell"

  # Node.js REPL
  node_repl:
    type: script
    command: "node"
    terminal: true
    description: "Node.js interactive shell"

  # Custom shell
  custom_shell:
    type: script
    command: "/usr/bin/zsh"
    terminal: true
    keep_open: true
    description: "Zsh terminal"
```

## Shell-Specific Arguments

Different shells use different command flags:

| Shell | Command Flag | Example |
|-------|--------------|---------|
| bash/zsh/sh | `-c` | `bash -c "echo hello"` |
| PowerShell | `-Command` | `powershell -Command "Write-Host 'hello'"` |
| PowerShell Core | `-c` | `pwsh -c "Write-Host 'hello'"` |
| cmd | `/c` | `cmd /c "echo hello"` |
| Python | `-c` | `python -c "print('hello')"` |
| Node.js | `-e` | `node -e "console.log('hello')"` |
| Ruby | `-e` | `ruby -e "puts 'hello'"` |

## Environment Variables

Shell-specific environment variables are automatically set:

### Python
- `PYTHONUNBUFFERED=1` - Ensures output isn't buffered

### Node.js
- `NODE_NO_WARNINGS=1` - Suppresses warnings

### PowerShell
- `PSModulePath` - Cleared to avoid conflicts

## Troubleshooting

### Shell Not Found

If a shell isn't found:

1. **Check Installation**: Ensure the shell is installed
   ```bash
   which python
   where powershell
   ```

2. **Use Full Path**: Specify the full path
   ```yaml
   shell: "/usr/local/bin/python3"
   ```

3. **Check PATH**: Ensure shell is in system PATH

### Wrong Shell Version

To use a specific version:

```yaml
grimoire:
  python3_specific:
    type: script
    command: "print('Hello')"
    shell: "python3"  # Not just "python"
    
  node_version:
    type: script
    command: "console.log(process.version)"
    shell: "/usr/local/bin/node"  # Full path
```

### Permission Errors

If you get permission errors:

1. **Check Executable**: Ensure shell is executable
   ```bash
   chmod +x /path/to/shell
   ```

2. **Use System Shell**: Fall back to system default
   ```yaml
   shell: ""  # Uses system default
   ```

## Best Practices

1. **Use Shebangs**: For portability, use shebang lines
   ```yaml
   command: |
     #!/usr/bin/env python3
     # Script content
   ```

2. **Test Shell Availability**: Check if shell exists before using
   ```yaml
   command: |
     if command -v python3 &> /dev/null; then
       python3 -c "print('Python available')"
     else
       echo "Python not found"
     fi
   ```

3. **Handle Cross-Platform**: Use different shells per platform
   ```yaml
   # In spellbook.windows.yml
   shell: "powershell"
   
   # In spellbook.darwin.yml
   shell: "zsh"
   
   # In spellbook.linux.yml
   shell: "bash"
   ```

4. **Error Handling**: Always handle potential shell errors
   ```yaml
   command: |
     set -e  # Exit on error
     # Your script here
   ```