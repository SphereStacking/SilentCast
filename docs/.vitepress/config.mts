import { defineConfig } from 'vitepress'
import llmstxt from 'vitepress-plugin-llms'
import { copyOrDownloadAsMarkdownButtons } from 'vitepress-plugin-llms'
import tailwindcss from '@tailwindcss/vite'

// Get version from environment variable or default
const version = process.env.VITEPRESS_VERSION || '0.1.0-dev'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: "/SilentCast/",
  title: "SilentCast",
  description: "Silent hotkey-driven task runner for developers",
  vite: {
    plugins: [llmstxt(),tailwindcss()],
  },
  markdown: {
    config(md) {
      md.use(copyOrDownloadAsMarkdownButtons)
    }
  },
  head: [
    ['link', { rel: 'icon', href: '/logo.svg' }],
    ['meta', { name: 'theme-color', content: '#3c8772' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'en' }],
    ['meta', { property: 'og:title', content: 'SilentCast | Silent Hotkey Task Runner' }],
    ['meta', { property: 'og:description', content: 'Execute tasks instantly with keyboard shortcuts. Works on Windows and macOS. Lightweight and developer-friendly.' }],
    ['meta', { property: 'og:site_name', content: 'SilentCast' }],
    ['meta', { property: 'og:image', content: 'https://SilentCast.SphereStacking.com/og-image.png' }],
    ['meta', { property: 'og:url', content: 'https://SilentCast.SphereStacking.com/' }],
  ],

  cleanUrls: true,
  lastUpdated: true,
  
  // Ignore dead links in README.md files
  ignoreDeadLinks: [
    // Ignore localhost URLs in development
    /^https?:\/\/localhost/
  ],

  themeConfig: {
    logo: '/logo.svg',
    siteTitle: 'SilentCast',
    
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Config', link: '/config/' },
      { text: 'API', link: '/api/' },
      { text: 'Development', link: '/development/setup' },
      { text: 'Troubleshooting', link: '/troubleshooting/' },
      {
        text: `v${version}`,
        items: [
          { text: 'Changelog', link: '/CHANGELOG' },
          { text: 'Releases', link: 'https://github.com/SphereStacking/SilentCast/releases' },
          { text: 'Contributing', link: '/contributing' },
          { text: 'TDD Guide', link: '/guide/tdd-development' }
        ]
      }
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Introduction',
          collapsed: false,
          items: [
            { text: 'What is SilentCast?', link: '/guide/what-is-silentcast' },
            { text: 'Getting Started', link: '/guide/getting-started' },
            { text: 'Installation', link: '/guide/installation' }
          ]
        },
        {
          text: 'Core Concepts',
          collapsed: false,
          items: [
            { text: 'Configuration', link: '/guide/configuration' },
            { text: 'Spells (Hotkeys)', link: '/guide/spells' },
            { text: 'Grimoire (Actions)', link: '/guide/grimoire' },
            { text: 'Configuration Samples', link: '/guide/samples' }
          ]
        },
        {
          text: 'Features',
          collapsed: false,
          items: [
            { text: 'Script Execution', link: '/guide/scripts' },
            { text: 'Browser Detection', link: '/guide/browser-detection' },
            { text: 'Browser Launcher', link: '/guide/browser-launcher' },
            { text: 'Custom Shells', link: '/guide/custom-shells' },
            { text: 'Terminal Customization', link: '/guide/terminal-customization' },
            { text: 'Force Terminal', link: '/guide/force-terminal' }
          ]
        },
        {
          text: 'Platform & Services',
          collapsed: true,
          items: [
            { text: 'Platform Support', link: '/guide/platforms' },
            { text: 'Windows Service', link: '/guide/windows-service' },
            { text: 'macOS Service', link: '/guide/macos-service' },
            { text: 'Linux Service', link: '/guide/linux-service' },
            { text: 'Windows Guide', link: '/guide/windows' }
          ]
        },
        {
          text: 'Management',
          collapsed: true,
          items: [
            { text: 'Auto-start', link: '/guide/auto-start' },
            { text: 'Updates', link: '/guide/updates' },
            { text: 'Self-update', link: '/guide/self-update' },
            { text: 'Update Notifications', link: '/guide/update-notifications' },
            { text: 'Timeout Notifications', link: '/guide/timeout-notifications' },
            { text: 'Backup & Restore', link: '/guide/backup-restore' }
          ]
        },
        {
          text: 'Advanced',
          collapsed: true,
          items: [
            { text: 'Environment Variables', link: '/guide/env-vars' },
            { text: 'Logging', link: '/guide/logging' },
            { text: 'Performance Optimization', link: '/guide/performance-optimization' },
            { text: 'Interpreter Mode', link: '/guide/interpreter-mode' }
          ]
        },
        {
          text: 'Reference',
          collapsed: true,
          items: [
            { text: 'CLI Reference', link: '/guide/cli-reference' },
            { text: 'FAQ', link: '/guide/faq' },
            { text: 'Troubleshooting', link: '/guide/troubleshooting' }
          ]
        },
        {
          text: 'Development',
          collapsed: true,
          items: [
            { text: 'TDD Development', link: '/guide/tdd-development' }
          ]
        }
      ],
      '/config/': [
        {
          text: 'Configuration',
          items: [
            { text: 'Overview', link: '/config/' },
            { text: 'Configuration Guide', link: '/config/configuration-guide' },
            { text: 'File Structure', link: '/config/file-structure' }
          ]
        }
      ],
      '/api/': [
        {
          text: 'Reference',
          items: [
            { text: 'CLI Options', link: '/api/' },
            { text: 'Key Names', link: '/api/key-names' },
            { text: 'Exit Codes', link: '/api/exit-codes' },
            { text: 'Environment Variables', link: '/api/env-vars' },
            { text: 'Architecture', link: '/api/architecture' },
            { text: 'Building', link: '/api/building' },
            { text: 'Testing', link: '/api/testing' },
            { text: 'Contributing', link: '/api/contributing' }
          ]
        }
      ],
      '/development/': [
        {
          text: 'Development',
          items: [
            { text: 'Setup Guide', link: '/development/setup' },
            { text: 'TDD Best Practices', link: '/development/tdd-best-practices' },
            { text: 'Error Handling', link: '/development/error-handling' },
            { text: 'Architecture Improvements', link: '/development/architecture-improvements' },
            { text: 'Code Quality Analysis', link: '/development/code_quality_analysis_report' },
            { text: 'Code Quality Priorities', link: '/development/code_quality_fix_priorities' }
          ]
        }
      ],
      '/troubleshooting/': [
        {
          text: 'Troubleshooting',
          items: [
            { text: 'Overview', link: '/troubleshooting/' },
            { text: 'Installation Issues', link: '/troubleshooting/installation' },
            { text: 'Configuration Issues', link: '/troubleshooting/configuration' },
            { text: 'Hotkey Issues', link: '/troubleshooting/hotkeys' },
            { text: 'Action Issues', link: '/troubleshooting/actions' },
            { text: 'Performance Issues', link: '/troubleshooting/performance' },
            { text: 'Permission Issues', link: '/troubleshooting/permissions' },
            { text: 'Platform-Specific', link: '/troubleshooting/platform-specific' },
            { text: 'Debugging', link: '/troubleshooting/debugging' },
            { text: 'FAQ', link: '/troubleshooting/faq' },
            { text: 'Support', link: '/troubleshooting/support' }
          ]
        }
      ]
    },

    editLink: {
      pattern: 'https://github.com/SphereStacking/SilentCast/edit/main/docs/:path',
      text: 'Edit this page on GitHub'
    },

    search: {
      provider: 'local',
      options: {
        translations: {
          button: {
            buttonText: 'Search',
            buttonAriaLabel: 'Search docs'
          },
          modal: {
            footer: {
              selectText: 'to select',
              navigateText: 'to navigate',
              closeText: 'to close',
            },
          }
        }
      }
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2025-present SphereStacking'
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/SphereStacking/SilentCast' },
      { icon: 'x', link: 'https://x.com/SphereStacking' }
    ]
  }
})
