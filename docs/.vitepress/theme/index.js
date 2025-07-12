import { h } from 'vue'
import DefaultTheme from 'vitepress/theme'
import './style.css'

// カスタムコンポーネントをインポート
// import CustomHero from './components/CustomHero.vue'
// import VersionDisplay from './components/VersionDisplay.vue'

export default {
  extends: DefaultTheme,
  Layout: () => {
    return h(DefaultTheme.Layout, null, {
      // スロットをカスタマイズ
      // 'home-hero-before': () => h(VersionDisplay),
      // 'home-features-after': () => h(CustomSection),
    })
  },
  enhanceApp({ app }) {
    // グローバルコンポーネントを登録
    // app.component('CustomHero', CustomHero)
    // app.component('VersionDisplay', VersionDisplay)
  }
}