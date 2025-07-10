import { defineConfig } from 'vitepress'
import llmstxt from 'vitepress-plugin-llms'
import { copyOrDownloadAsMarkdownButtons } from 'vitepress-plugin-llms'


// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: "/SilentCast/",
  title: "SilentCast",
  description: "Silent hotkey-driven task runner for developers",
  vite: {
    plugins: [llmstxt()],
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

  themeConfig: {
    logo: '/logo.svg',
    siteTitle: 'SilentCast',
    
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Config', link: '/config/' },
      { text: 'API', link: '/api/' },
      {
        text: 'v1.0.0',
        items: [
          { text: 'Changelog', link: 'https://github.com/SphereStacking/SilentCast/releases' },
          { text: 'Contributing', link: '/contributing' }
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
          text: 'Essentials',
          collapsed: false,
          items: [
            { text: 'Configuration', link: '/guide/configuration' },
            { text: 'Configuration Samples', link: '/guide/samples' },
            { text: 'Shortcuts', link: '/guide/shortcuts' },
            { text: 'Actions', link: '/guide/actions' },
            { text: 'Platform Support', link: '/guide/platforms' }
          ]
        },
        {
          text: 'Advanced',
          collapsed: false,
          items: [
            { text: 'Script Execution', link: '/guide/scripts' },
            { text: 'Environment Variables', link: '/guide/env-vars' },
            { text: 'Auto-start', link: '/guide/auto-start' },
            { text: 'Logging', link: '/guide/logging' }
          ]
        }
      ],
      '/config/': [
        {
          text: 'Configuration',
          items: [
            { text: 'Overview', link: '/config/' },
            { text: 'File Structure', link: '/config/file-structure' },
            { text: 'Daemon Settings', link: '/config/daemon' },
            { text: 'Hotkey Settings', link: '/config/hotkeys' },
            { text: 'Spells & Grimoire', link: '/config/spells' },
            { text: 'Platform Overrides', link: '/config/platform-overrides' },
            { text: 'Examples', link: '/config/examples' }
          ]
        }
      ],
      '/api/': [
        {
          text: 'Reference',
          items: [
            { text: 'CLI Options', link: '/api/' },
            { text: 'Key Names', link: '/api/key-names' },
            { text: 'Environment Variables', link: '/api/env-vars' },
            { text: 'Exit Codes', link: '/api/exit-codes' }
          ]
        },
        {
          text: 'Development',
          items: [
            { text: 'Architecture', link: '/api/architecture' },
            { text: 'Building', link: '/api/building' },
            { text: 'Testing', link: '/api/testing' },
            { text: 'Contributing', link: '/api/contributing' }
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
