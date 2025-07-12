# Environment Variables Guide

SilentCast supports environment variables throughout its configuration, allowing dynamic and flexible setups. This guide covers how to use environment variables effectively in your spellbook.

## Overview

Environment variables can be used in:
- Configuration values
- Command execution
- Script environments
- Path specifications

## Using Environment Variables

### Basic Syntax

Use `${VARIABLE_NAME}` syntax in your configuration:

```yaml
grimoire:
  open_project:
    type: app
    command: "${EDITOR}"  # Uses $EDITOR environment variable
    working_dir: "${PROJECT_DIR}"
    
  build_project:
    type: script
    command: "npm run build"
    env:
      NODE_ENV: "${BUILD_ENV:-production}"  # With default value
```

### Default Values

Provide fallbacks with `${VAR:-default}` syntax:

```yaml
grimoire:
  editor:
    type: app
    command: "${EDITOR:-code}"  # Use VS Code if $EDITOR not set
    
  terminal:
    type: app  
    command: "${TERMINAL:-wt}"  # Default to Windows Terminal
```

## Built-in Variables

SilentCast provides these variables automatically:

| Variable | Description | Example |
|----------|-------------|---------|
| `${HOME}` | User home directory | `/home/username` or `C:\Users\username` |
| `${USER}` | Current username | `john` |
| `${PWD}` | Current working directory | `/home/user/projects` |
| `${SILENTCAST_CONFIG}` | Config directory path | `~/.config/silentcast` |
| `${SILENTCAST_VERSION}` | SilentCast version | `1.0.0` |
| `${SILENTCAST_OS}` | Operating system | `windows`, `darwin` |
| `${SILENTCAST_ARCH}` | System architecture | `amd64`, `arm64` |

### Platform-Specific Variables

::: code-group

```yaml [Windows]
# Windows-specific built-ins
${APPDATA}        # C:\Users\username\AppData\Roaming
${LOCALAPPDATA}   # C:\Users\username\AppData\Local
${PROGRAMFILES}   # C:\Program Files
${USERPROFILE}    # C:\Users\username
${COMPUTERNAME}   # Machine name
```

```yaml [macOS]
# Unix-like built-ins
${HOME}           # /Users/username
${USER}           # username
${HOSTNAME}       # machine.local
${SHELL}          # /bin/bash
${TMPDIR}         # /tmp or /var/tmp
```

:::

## Configuration Examples

### Dynamic Paths

```yaml
grimoire:
  # Use different editors based on environment
  smart_editor:
    type: app
    command: "${EDITOR:-${VISUAL:-code}}"
    
  # Project-specific paths
  open_work_project:
    type: app
    command: "code"
    args: ["${WORK_PROJECT_DIR}/src"]
    
  # User-specific config
  edit_config:
    type: app
    command: "${EDITOR}"
    args: ["${HOME}/.config/myapp/config.json"]
```

### Conditional Commands

```yaml
grimoire:
  # Development vs Production
  run_server:
    type: script
    command: "npm run ${NODE_ENV:-development}"
    env:
      PORT: "${PORT:-3000}"
      API_URL: "${API_URL:-http://localhost:8080}"
      
  # Platform-aware commands  
  open_file_manager:
    type: app
    command: "${FILE_MANAGER:-explorer}"  # Set per platform
```

### Secret Management

```yaml
grimoire:
  # Never hardcode secrets
  deploy_app:
    type: script
    command: |
      aws s3 sync ./dist s3://${BUCKET_NAME}/
    env:
      AWS_ACCESS_KEY_ID: "${AWS_ACCESS_KEY_ID}"
      AWS_SECRET_ACCESS_KEY: "${AWS_SECRET_ACCESS_KEY}"
      AWS_REGION: "${AWS_REGION:-us-east-1}"
```

## Setting Environment Variables

### For SilentCast Process

Set variables before starting SilentCast:

::: code-group

```bash [macOS]
# Export for current session
export EDITOR=vim
export PROJECT_DIR=~/projects/myapp
silentcast

# Or inline
EDITOR=vim PROJECT_DIR=~/projects silentcast

# In shell profile (~/.bashrc, ~/.zshrc)
export SILENTCAST_CONFIG=~/.config/silentcast
export EDITOR=code
```

```powershell [Windows]
# Set for current session
$env:EDITOR = "notepad++"
$env:PROJECT_DIR = "C:\Projects\MyApp"
silentcast

# Set permanently (user)
[Environment]::SetEnvironmentVariable("EDITOR", "code", "User")

# Set permanently (system)
[Environment]::SetEnvironmentVariable("EDITOR", "code", "Machine")
```

:::

### For Action Execution

Set variables for specific actions:

```yaml
grimoire:
  test_with_env:
    type: script
    command: "npm test"
    env:
      NODE_ENV: "test"
      CI: "true"
      TEST_TIMEOUT: "10000"
      # Mix with system variables
      CUSTOM_PATH: "${PATH}:/custom/bin"
```

### Using .env Files

Load variables from .env files:

```yaml
grimoire:
  load_env_and_run:
    type: script
    command: |
      # Load .env file
      if [ -f .env ]; then
        export $(cat .env | xargs)
      fi
      
      # Now use the variables
      echo "API_URL: $API_URL"
      npm run dev
```

## Advanced Usage

### Variable Expansion Order

Variables are expanded in this order:
1. Built-in SilentCast variables
2. System environment variables
3. Action-specific env settings
4. Default values (if specified)

```yaml
grimoire:
  expansion_example:
    type: script
    command: 'echo "User: ${USER}, Custom: ${CUSTOM_VAR:-default}"'
    env:
      CUSTOM_VAR: "override"  # This takes precedence
```

### Nested Variables

Variables can reference other variables:

```yaml
grimoire:
  nested_vars:
    type: script
    command: "echo ${MESSAGE}"
    env:
      NAME: "${USER}"
      TIME: "$(date +%H:%M)"
      MESSAGE: "Hello ${NAME}, it's ${TIME}"
```

### Escaping Variables

To use literal `${...}` without expansion:

```yaml
grimoire:
  # Use single quotes to prevent expansion
  literal_dollar:
    type: script
    command: 'echo "Literal: ${NOT_EXPANDED}"'
    
  # Or escape with backslash
  escaped_var:
    type: script
    command: "echo \"Use \\${VAR} syntax\""
```

## Platform-Specific Patterns

### Cross-Platform Paths

```yaml
grimoire:
  # Use variables for platform differences
  cross_platform_config:
    type: app
    command: "${EDITOR}"
    args: ["${CONFIG_DIR}/settings.json"]
    
# In platform overrides:
# spellbook.windows.yml
env:
  CONFIG_DIR: "${APPDATA}/MyApp"
  
# spellbook.mac.yml  
env:
  CONFIG_DIR: "${HOME}/Library/Application Support/MyApp"
```

### Platform Detection

```yaml
grimoire:
  platform_aware:
    type: script
    command: |
      case "${SILENTCAST_OS}" in
        windows)
          echo "Running on Windows"
          ;;
        darwin)
          echo "Running on macOS"
          ;;
      esac
```

## Common Patterns

### Development Workflows

```yaml
# Set these in your shell profile
export DEV_ROOT=~/development
export GITHUB_USER=myusername
export EDITOR=code

grimoire:
  clone_repo:
    type: script
    command: |
      read -p "Repo name: " REPO
      git clone "https://github.com/${GITHUB_USER}/${REPO}.git" \
        "${DEV_ROOT}/${REPO}"
      cd "${DEV_ROOT}/${REPO}"
      ${EDITOR} .
    
  new_project:
    type: script
    command: |
      PROJECT_NAME="$1"
      mkdir -p "${DEV_ROOT}/${PROJECT_NAME}"
      cd "${DEV_ROOT}/${PROJECT_NAME}"
      npm init -y
      ${EDITOR} .
```

### API Keys and Secrets

```yaml
# Store secrets in environment, not config
grimoire:
  api_call:
    type: script
    command: |
      curl -H "Authorization: Bearer ${API_TOKEN}" \
           "${API_BASE_URL}/endpoint"
    env:
      API_BASE_URL: "${API_BASE_URL:-https://api.example.com}"
      # API_TOKEN comes from system environment
```

### Dynamic Configuration

```yaml
grimoire:
  adaptive_build:
    type: script
    command: |
      # Detect environment
      if [ "${CI}" = "true" ]; then
        npm run build:ci
      elif [ "${NODE_ENV}" = "production" ]; then
        npm run build:prod
      else
        npm run build:dev
      fi
```

## Best Practices

### 1. Use Descriptive Names

```yaml
# Good
${PROJECT_ROOT}
${API_BASE_URL}
${BUILD_OUTPUT_DIR}

# Avoid
${DIR}
${URL}
${OUT}
```

### 2. Provide Defaults

```yaml
grimoire:
  safe_command:
    type: script
    command: "deploy --env ${DEPLOY_ENV:-staging}"
    env:
      TIMEOUT: "${TIMEOUT:-30}"
      RETRIES: "${RETRIES:-3}"
```

### 3. Document Required Variables

```yaml
grimoire:
  requires_vars:
    type: script
    command: |
      # Required: AWS_PROFILE, BUCKET_NAME
      # Optional: AWS_REGION (default: us-east-1)
      
      : ${AWS_PROFILE:?Error: AWS_PROFILE not set}
      : ${BUCKET_NAME:?Error: BUCKET_NAME not set}
      
      aws s3 ls s3://${BUCKET_NAME}/
```

### 4. Group Related Variables

```yaml
# Group by prefix
export MYAPP_API_URL=https://api.example.com
export MYAPP_API_KEY=secret123
export MYAPP_TIMEOUT=30

grimoire:
  myapp_sync:
    type: script
    command: "myapp sync"
    env:
      API_URL: "${MYAPP_API_URL}"
      API_KEY: "${MYAPP_API_KEY}"
      TIMEOUT: "${MYAPP_TIMEOUT}"
```

## Troubleshooting

### Variables Not Expanding

<details>
<summary>Check variable syntax</summary>

```yaml
# Correct
command: "${HOME}/scripts/run.sh"

# Incorrect
command: "$HOME/scripts/run.sh"     # Missing braces
command: "$(HOME)/scripts/run.sh"   # Wrong syntax
command: '%HOME%/scripts/run.sh'    # Windows syntax won't work
```

</details>

<details>
<summary>Variable not found</summary>

Debug with echo:
```yaml
grimoire:
  debug_vars:
    type: script
    command: |
      echo "HOME: ${HOME}"
      echo "CUSTOM: ${CUSTOM_VAR:-not set}"
      echo "All vars:"
      env | sort
```

</details>

<details>
<summary>Platform differences</summary>

Use platform-specific overrides:
```yaml
# spellbook.yml
grimoire:
  open_app:
    type: app
    command: "${APP_PATH}"

# spellbook.windows.yml
env:
  APP_PATH: "C:\\Program Files\\MyApp\\app.exe"

# spellbook.mac.yml
env:
  APP_PATH: "/Applications/MyApp.app"
```

</details>

### Security Notes

1. **Never commit secrets**: Keep sensitive variables in environment only
2. **Use read-only variables**: For paths and configs that shouldn't change
3. **Validate input**: Check variables before using in commands
4. **Restrict access**: Use proper file permissions for .env files

## Common Environment Variables

### Development

```bash
# Editors
EDITOR=code           # or vim, emacs, sublime
VISUAL=code          # GUI editor

# Paths
PROJECT_ROOT=~/projects
DOTFILES=~/.dotfiles
SCRIPTS=~/scripts

# Development
NODE_ENV=development
DEBUG=true
VERBOSE=1
```

### Tools

```bash
# Version managers
NVM_DIR=~/.nvm
PYENV_ROOT=~/.pyenv
RBENV_ROOT=~/.rbenv

# Package managers
NPM_TOKEN=xxx
CARGO_HOME=~/.cargo
GOPATH=~/go

# Cloud
AWS_PROFILE=default
GOOGLE_APPLICATION_CREDENTIALS=~/keys/gcp.json
AZURE_SUBSCRIPTION_ID=xxx
```

## Next Steps

- Review [Configuration Guide](/guide/configuration) for variable usage
- Check [Script Execution](/guide/scripts) for environment examples
- See [Platform Guide](/guide/platforms) for OS-specific variables
- Explore [Configuration Samples](/guide/samples) for real-world usage