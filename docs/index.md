---
layout: home

hero:
  name: "SilentCast"
  text: "Silent Hotkey Task Runner"
  tagline: Execute tasks instantly with keyboard shortcuts. Works on Windows and macOS. Lightweight and developer-friendly.
  pretext:
    text: "⚠️ Currently in Development"
    color: "warning"
  image:
    src: /logo.svg
    alt: SilentCast
  actions:
    - theme: brand
      text: Get Started 🚀
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

:::warning Development Status
SilentCast is currently under active development. Features may change and bugs may exist. Please use at your own risk and report any issues on [GitHub](https://github.com/SphereStacking/silentcast/issues).
:::

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

<div class="grid md:grid-cols-2 gap-8 my-12">

<div class="space-y-6">

### 🎯 **Built for Developers**

Designed specifically for developer workflows. Launch your IDE, run build scripts, check git status with simple keystrokes.

### 🧩 **Simple Configuration**

Uses standard YAML format for configuration. No scripting languages or complex GUIs required.

</div>

<div class="space-y-6">

### 🚀 **Lightweight**

Built with Go for minimal resource usage. No heavy runtimes or unnecessary features.

### 🔒 **Local Only**

Runs entirely on your machine. No cloud services, telemetry, or account required.

</div>

</div>


## 🚀 Ready to Get Started?

<div class="my-16 text-center">
  <div class="inline-block relative px-16 py-12 rounded-3xl bg-gradient-to-br from-white/5 to-white/[0.02] border border-indigo-500/20 overflow-hidden transition-all duration-300 hover:scale-110 hover:border-indigo-500/40 hover:shadow-[0_20px_25px_-5px_rgba(99,102,241,0.1),0_10px_10px_-5px_rgba(99,102,241,0.04),inset_0_1px_0_0_rgba(255,255,255,0.1)] group">
    <div class="absolute inset-[-50%] bg-[radial-gradient(circle_at_center,rgba(99,102,241,0.3)_0%,transparent_70%)] opacity-0 transition-opacity duration-300 pointer-events-none group-hover:opacity-50 group-hover:animate-pulse"></div>
    <div class="relative z-10">
      <h3 class="text-3xl font-bold mb-4 bg-gradient-to-r from-indigo-500 to-purple-500 bg-clip-text text-transparent">
        Start Using Keyboard Shortcuts Today
      </h3>
      <p class="text-lg text-gray-600 dark:text-gray-400 mb-8 max-w-2xl mx-auto">
        Install SilentCast in seconds and boost your productivity with powerful hotkeys
      </p>
      <a href="/guide/getting-started" class="inline-flex items-center gap-2 px-8 py-3.5 border-4 border-gradient-to-r from-indigo-500 to-purple-600 font-semibold text-lg rounded-xl no-underline transition-all duration-300 relative overflow-hidden shadow-[0_10px_15px_-3px_rgba(99,102,241,0.3),0_4px_6px_-2px_rgba(99,102,241,0.2)] group/btn ">
        <span class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent opacity-0 transition-opacity duration-300 group-hover/btn:opacity-80"></span>
        <span class="relative z-10 ">Get Started</span>
        <span class="relative z-10 text-xl transition-transform duration-300 group-hover/btn:translate-x-1">→</span>
      </a>
    </div>
  </div>
</div>

<style>
/* Hero animations */
.VPHero .VPImage {
  /* animation: float 6s ease-in-out infinite; */
  filter: drop-shadow(0 25px 25px rgb(0 0 0 / 0.15));
  width: 11rem;
  height: 11rem;
  transition: transform 0.3s ease;
}

@media (min-width: 768px) {
  .VPHero .VPImage {
    width: 15rem;
    height: 15rem;
  }
}

@media (min-width: 1024px) {
  .VPHero .VPImage {
    width: 18rem;
    height: 18rem;
  }
}
</style>

