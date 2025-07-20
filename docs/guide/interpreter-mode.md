# Interpreter Mode

Interpreter mode allows you to run scripts directly in interpreters without shell wrappers, providing cleaner execution for non-shell scripts.

## Overview

Interpreter mode:
- Bypasses shell completely for direct interpreter execution
- Automatically detects suitable interpreters from script content
- Provides better performance for interpreter-based scripts
- Eliminates shell parsing and escaping issues
- Offers enhanced security by avoiding shell interpretation

## When to Use Interpreter Mode

Use interpreter mode when:
- Running Python, Node.js, Ruby, or other interpreter scripts
- You want to avoid shell escaping issues
- The script doesn't use shell features (pipes, redirects, etc.)
- You need precise control over arguments
- Performance is critical

**Don't use interpreter mode when:**
- Script uses shell features (pipes, environment variable expansion)
- Script needs to run shell commands
- Script requires shell-specific functionality

## Configuration

### Basic Interpreter Mode

Enable interpreter mode with the `interpreter` flag:

```yaml
grimoire:
  python_direct:
    type: script
    command: |
      import sys
      print(f"Python {sys.version}")
      print("Running in interpreter mode")
    interpreter: true
    description: "Python in interpreter mode"
    show_output: true
```

### Specifying Interpreter

Combine with `shell` field to specify interpreter:

```yaml
grimoire:
  node_direct:
    type: script
    command: |
      console.log(`Node.js ${process.version}`);
      console.log('Running in interpreter mode');
    shell: "node"
    interpreter: true
    description: "Node.js in interpreter mode"
    show_output: true
```

## Automatic Detection

SilentCast can automatically detect interpreters from:

### 1. Shebang Lines

```yaml
grimoire:
  auto_python:
    type: script
    command: |
      #!/usr/bin/env python3
      import datetime
      print(f"Current time: {datetime.datetime.now()}")
    interpreter: true
    # No shell specified - will detect Python from shebang
```

### 2. Script Content

```yaml
grimoire:
  auto_detect:
    type: script
    command: |
      import os
      import sys
      
      def main():
          print("Auto-detected Python script")
          print(f"Working directory: {os.getcwd()}")
      
      if __name__ == "__main__":
          main()
    interpreter: true
    # Will detect Python from imports and syntax
```

## Examples

### Python Scripts

```yaml
spells:
  p,1: python_hello
  p,2: python_json
  p,3: python_args

grimoire:
  python_hello:
    type: script
    command: |
      print("Hello from Python interpreter mode!")
      import sys
      print(f"Arguments: {sys.argv[1:]}")
    interpreter: true
    args: ["world", "test"]
    description: "Python with arguments"
    show_output: true

  python_json:
    type: script
    command: |
      import json
      data = {
          "name": "SilentCast",
          "mode": "interpreter",
          "languages": ["Python", "Node.js", "Ruby"]
      }
      print(json.dumps(data, indent=2))
    interpreter: true
    shell: "python"
    description: "Python JSON processing"
    show_output: true

  python_args:
    type: script
    command: |
      import sys
      import argparse
      
      parser = argparse.ArgumentParser()
      parser.add_argument("--name", default="World")
      parser.add_argument("--count", type=int, default=1)
      
      args = parser.parse_args()
      
      for i in range(args.count):
          print(f"Hello, {args.name}! ({i+1})")
    interpreter: true
    args: ["--name", "SilentCast", "--count", "3"]
    description: "Python with argument parsing"
    show_output: true
```

### Node.js Scripts

```yaml
grimoire:
  node_hello:
    type: script
    command: |
      const args = process.argv.slice(2);
      console.log('Hello from Node.js interpreter mode!');
      console.log('Arguments:', args);
    interpreter: true
    shell: "node"
    args: ["world", "test"]
    description: "Node.js with arguments"
    show_output: true

  node_fs:
    type: script
    command: |
      const fs = require('fs');
      const path = require('path');
      
      console.log('Current directory:', process.cwd());
      console.log('Files:');
      fs.readdirSync('.').forEach(file => {
          const stats = fs.statSync(file);
          console.log(`  ${file} (${stats.isDirectory() ? 'dir' : 'file'})`);
      });
    interpreter: true
    shell: "node"
    description: "Node.js file system"
    show_output: true
```

### Ruby Scripts

```yaml
grimoire:
  ruby_hello:
    type: script
    command: |
      puts "Hello from Ruby interpreter mode!"
      puts "Arguments: #{ARGV.join(', ')}"
      
      # Ruby-specific features
      fruits = %w[apple banana orange]
      fruits.each_with_index do |fruit, i|
        puts "#{i + 1}. #{fruit.capitalize}"
      end
    interpreter: true
    shell: "ruby"
    args: ["world", "test"]
    description: "Ruby with arguments"
    show_output: true
```

## Comparison: Shell vs Interpreter Mode

### Shell Mode (Default)

```yaml
python_shell:
  type: script
  command: 'print("Hello from shell mode")'
  shell: "python"
  # Executes: python -c 'print("Hello from shell mode")'
```

**Execution flow:**
1. SilentCast → Shell → Python interpreter
2. Shell handles escaping and parsing
3. Shell passes command to Python

### Interpreter Mode

```yaml
python_interpreter:
  type: script
  command: 'print("Hello from interpreter mode")'
  shell: "python"
  interpreter: true
  # Executes: python -c 'print("Hello from interpreter mode")'
```

**Execution flow:**
1. SilentCast → Python interpreter (direct)
2. No shell involvement
3. Direct argument passing

## Performance Benefits

Interpreter mode provides:

1. **Faster startup** - No shell process creation
2. **Lower memory usage** - One less process
3. **Better argument handling** - No shell escaping
4. **Cleaner process tree** - Direct parent-child relationship

## Security Benefits

Interpreter mode offers:

1. **No shell injection** - Arguments passed directly
2. **Precise quoting** - No shell interpretation
3. **Reduced attack surface** - Fewer processes involved
4. **Safer environment** - No shell variable expansion

## Best Practices

### 1. Use for Pure Interpreter Scripts

Good candidates:
```yaml
# Python data processing
command: |
  import pandas as pd
  data = pd.read_csv('input.csv')
  result = data.groupby('category').sum()
  print(result)
interpreter: true

# Node.js API call
command: |
  const https = require('https');
  https.get('https://api.example.com/data', (res) => {
    // Handle response
  });
interpreter: true
```

### 2. Avoid for Shell-Heavy Scripts

Avoid interpreter mode for:
```yaml
# This needs shell features
command: |
  import subprocess
  result = subprocess.run('ls -la | grep ".py"', shell=True)
# Better to use regular shell mode for this
```

### 3. Handle Arguments Properly

```yaml
# Good: Use args for command-line arguments
command: |
  import sys
  print(f"Args: {sys.argv[1:]}")
interpreter: true
args: ["arg1", "arg2"]

# Avoid: Embedding arguments in command string
command: 'print("arg1 arg2")'  # Less flexible
```

### 4. Environment Variables

```yaml
command: |
  import os
  print(f"Custom var: {os.environ.get('MY_VAR', 'default')}")
interpreter: true
env:
  MY_VAR: "custom value"
```

## Troubleshooting

### Script Not Recognized as Interpreter Script

If SilentCast doesn't detect the right interpreter:

1. **Add explicit shell specification:**
   ```yaml
   shell: "python"
   interpreter: true
   ```

2. **Use shebang line:**
   ```yaml
   command: |
     #!/usr/bin/env python3
     # Your script here
   ```

3. **Check interpreter availability:**
   ```bash
   which python
   which node
   which ruby
   ```

### Arguments Not Working

Ensure you're using the `args` field:

```yaml
# Correct
interpreter: true
args: ["--option", "value"]

# Incorrect - will be treated as code
command: "script.py --option value"
```

### Environment Variables Not Available

Make sure environment is set correctly:

```yaml
interpreter: true
env:
  PYTHONPATH: "/custom/path"
  NODE_ENV: "development"
```

## Migration from Shell Mode

To migrate existing scripts to interpreter mode:

1. **Identify pure interpreter scripts:**
   - No shell pipes (`|`)
   - No shell redirects (`>`, `>>`)
   - No shell variables (`$VAR`)
   - No shell commands

2. **Add interpreter flag:**
   ```yaml
   # Before
   shell: "python"
   
   # After
   shell: "python"
   interpreter: true
   ```

3. **Move command-line arguments:**
   ```yaml
   # Before
   command: "script.py --option value"
   
   # After
   command: "# script content here"
   args: ["--option", "value"]
   ```

4. **Test thoroughly** - Verify behavior is identical