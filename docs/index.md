---
layout: home

hero:
  name: "SilentCast"
  text: "Silent Hotkey Task Runner"
  tagline: Execute tasks instantly with keyboard shortcuts. Works on Windows and macOS. Lightweight and developer-friendly.
  image:
    src: /logo.svg
    alt: SilentCast
  background:
    image: linear-gradient(135deg, #667eea 0%, #764ba2 100%)
    filter: blur(72px)
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/SphereStacking/silentcast

features:
  - icon: 🎯
    title: Global Hotkeys
    details: System-wide keyboard shortcuts that work in any application. No need to focus a specific window.
    link: /guide/shortcuts
    linkText: Learn about shortcuts
    
  - icon: ⚡
    title: Fast and Lightweight
    details: Minimal CPU usage (~1%), low memory footprint (~15MB). Built with Go for efficiency.
    link: /guide/what-is-silentcast#performance
    linkText: Performance details
    
  - icon: 🌐
    title: Cross-Platform
    details: Works on Windows and macOS with platform-specific optimizations.
    link: /guide/platforms
    linkText: Platform guides
    
  - icon: 🎨
    title: Multi-Key Sequences
    details: Support for VS Code-style key sequences like 'g,s' for git status. Configurable timeouts.
    link: /guide/shortcuts#sequences
    linkText: Sequence shortcuts
    
  - icon: 🛠️
    title: Developer Focused
    details: YAML configuration, auto-reload on changes, structured logging, and detailed documentation.
    link: /guide/configuration
    linkText: Configuration guide
    
  - icon: 🤖
    title: Automation Ready
    details: Launch applications with arguments, execute scripts, set working directories, and use environment variables.
    link: /guide/scripts
    linkText: Automation guide
---

## Quick Example

::: code-group

```yaml [spellbook.yml]
# Configuration example
daemon:
  auto_start: true
  log_level: info

hotkeys:
  prefix: "alt+space"  # Activation key
  timeout: 1000

spells:
  # Single key shortcuts
  e: "editor"          # Alt+Space, E -> Open editor
  t: "terminal"        # Alt+Space, T -> Open terminal
  
  # Multi-key sequences (VS Code style)
  "g,s": "git_status"  # Alt+Space, G, S -> Git status
  "d,b": "docker_build" # Alt+Space, D, B -> Docker build

grimoire:
  editor:
    type: app
    command: "code"
    description: "Open VS Code"
  
  git_status:
    type: script
    command: "git status"
    description: "Show git status"
```

```bash [Installation]
# macOS
curl -sSL https://silentcast.dev/install.sh | bash

# Windows (PowerShell)
iwr -useb https://silentcast.dev/install.ps1 | iex

# Or download binaries
# https://github.com/SphereStacking/silentcast/releases
```

```bash [Usage]
# Start SilentCast
silentcast

# Start without system tray
silentcast --no-tray

# Use custom config
silentcast --config ~/my-spellbook.yml

# Then press Alt+Space and use your shortcuts
```

:::

## Why SilentCast?

<div class="features-comparison">

### 🎯 **Built for Developers**

Designed specifically for developer workflows. Launch your IDE, run build scripts, check git status with simple keystrokes.

### 🧩 **Simple Configuration**

Uses standard YAML format for configuration. No scripting languages or complex GUIs required.

### 🚀 **Lightweight**

Built with Go for minimal resource usage. No heavy runtimes or unnecessary features.

### 🔒 **Local Only**

Runs entirely on your machine. No cloud services, telemetry, or account required.

</div>

## Ready to Get Started?

<div class="cta-section">

Quick installation to start using keyboard shortcuts for your tasks.

[Get Started →](/guide/getting-started)

</div>

<style>

/* 6. Aurora - オーロラ */
.VPHero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  margin-bottom: 20px
}

/* Hero image (logo) customization */
.VPHero .VPImage {
  width: 180px !important;
  height: 180px !important;
  max-width: 180px !important;
  max-height: 180px !important;
}

/* Larger on desktop */
@media (min-width: 768px) {
  .VPHero .VPImage {
    width: 240px !important;
    height: 240px !important;
    max-width: 240px !important;
    max-height: 240px !important;
  }
}

/* Even larger on wide screens */
@media (min-width: 1200px) {
  .VPHero .VPImage {
    width: 300px !important;
    height: 300px !important;
    max-width: 300px !important;
    max-height: 300px !important;
  }
}

.features-comparison {
  margin: 3rem 0;
}

.features-comparison h3 {
  margin: 2rem 0 1rem;
}

.community-links {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin: 2rem 0;
}

.community-links a {
  flex: 1;
  min-width: 200px;
  padding: 1rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  text-align: center;
  transition: all 0.3s;
}

.community-links a:hover {
  border-color: var(--vp-c-brand);
  transform: translateY(-2px);
}

.cta-section {
  text-align: center;
  margin: 4rem 0;
  padding: 3rem;
  background: var(--vp-c-bg-soft);
  border-radius: 12px;
}

.cta-button {
  display: inline-block;
  padding: 12px 24px;
  background: var(--vp-c-brand);
  color: white;
  border-radius: 6px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.3s;
}

.cta-button:hover {
  background: var(--vp-c-brand-dark);
  transform: translateY(-2px);
}
</style>
