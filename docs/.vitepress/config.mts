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
      {
        text: `v${version}`,
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
            { text: 'Configuration Guide', link: '/config/configuration-guide' }
          ]
        }
      ],
      '/api/': [
        {
          text: 'Reference',
          items: [
            { text: 'CLI Options', link: '/api/' },
            { text: 'Key Names', link: '/api/key-names' },
            { text: 'Testing', link: '/api/testing' }
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
