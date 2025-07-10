# SilentCast ドキュメント

このディレクトリにはSilentCastのVitePressドキュメントが含まれています。

## 開発

```bash
# 依存関係のインストール
cd docs
npm install

# 開発サーバーを起動
npm run docs:dev
```

http://localhost:5173 でドキュメントを確認できます。

## ビルド

```bash
# 静的サイトをビルド
npm run docs:build

# ビルドしたサイトをプレビュー
npm run docs:preview
```

## 構造

```
docs/
├── .vitepress/       # VitePress設定
│   └── config.js     # サイト設定
├── guide/            # ガイド
│   ├── introduction.md
│   ├── getting-started.md
│   └── ...
├── config/           # 設定リファレンス
│   ├── index.md
│   └── ...
├── api/              # API/開発者向け
│   ├── build.md
│   └── ...
└── index.md          # ホームページ
```

## デプロイ

GitHub Pagesへのデプロイ：

```yaml
# .github/workflows/deploy-docs.yml
name: Deploy Docs

on:
  push:
    branches: [main]
    paths:
      - 'docs/**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-node@v4
        with:
          node-version: 18
          
      - run: |
          cd docs
          npm install
          npm run docs:build
          
      - uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs/.vitepress/dist
```